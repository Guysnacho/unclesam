package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var file_filters = []string{".git", ".go", ".mod", ".md", ".png"}
var transaction_filters = []string{"supabase", "tailwind", "tamagui", "digitalocean", "git", "sponsor", "notion"}

// Handle error checks
func check(e error, message string) {
	if e != nil {
		fmt.Printf("\n\n%s\n", message)
		panic(e)
	}
}

func main() {
	fmt.Println("I'm doing taxes I guess")
	data, err := os.ReadDir("./")

	check(err, "Failed to read directory.")

	fmt.Println("Searching for CSVs")

	var expenses = 0.00

	for _, file := range data {
		// Skip project files
		if check_filter(file.Name(), file_filters) {
			continue
		}

		// Check CSVs
		if check_filter(file.Name(), []string{".CSV", ".csv"}) {
			fmt.Println("\n\nProcessing - " + strings.Split(file.Name(), "_")[0])
			expenses += process(file)
		} else {
			fmt.Println(file.Name() + " - Is not processable")
		}
	}

	fmt.Println("\n\n###############################")
	fmt.Printf("Compiled Expenses - %f\n", expenses)
	fmt.Println("###############################")
}

func process(entry os.DirEntry) float64 {
	info, err := entry.Info()
	check(err, "Failed to process file")

	// Open file
	file, err := os.Open("./" + info.Name())
	check(err, "Failed to open file")

	reader := csv.NewReader(file)

	header, err := reader.Read()
	check(err, "Failed to read csv")

	desc_col, amt_col := parse_header(header)

	rows, err := reader.ReadAll()

	check(err, "Ran into an issue reading CSV")

	var expenses = 0.00

	for i, row := range rows {
		// Process all rows. I'd like to do it row by row but let's do it the dumb way first
		if !check_filter(row[desc_col], transaction_filters) {
			continue
		}
		amt, err := strconv.ParseFloat(row[amt_col], 64)

		check(err, "Failed to parse float")

		fmt.Printf("#%d | %s | $%f\n", i, row[desc_col], math.Abs(amt))
		expenses += math.Abs(amt)
	}

	fmt.Println("\n\n###############################")
	fmt.Printf("Total Expenses - %f\n", expenses)
	fmt.Println("###############################")
	return expenses
}

func check_filter(desc string, filters []string) bool {
	for _, filter := range filters {
		if strings.Contains(strings.ToLower(desc), filter) {
			return true
		}
	}
	return false
}

func parse_header(header_row []string) (int, int) {
	var NOT_FOUND = -123
	var description int = NOT_FOUND
	var amount int = NOT_FOUND

	fmt.Println("\nHeader Row")

	for idx, header := range header_row {
		fmt.Printf("%d | %s\n", idx, header)
		if strings.EqualFold(header, "description") {
			description = idx
		} else if strings.EqualFold(header, "amount") {
			amount = idx
		}
	}

	if description == NOT_FOUND || amount == NOT_FOUND {
		panic("Couldn't find a column for descriptions or amounts.")
	}

	return description, amount
}
