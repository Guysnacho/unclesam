package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

var file_filters = []string{".git", ".go", ".mod", ".md", ".png"}
var transaction_filters = []string{"supabase", "tailwind", "tamagui", "digitalocean", "git", "sponsor"}

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

	for _, file := range data {
		// Skip project files
		if check_filter(file.Name(), file_filters) {
			continue
		}

		// Check CSVs
		if check_filter(file.Name(), []string{".CSV", ".csv"}) {
			fmt.Println("\n\nProcessing - " + strings.Split(file.Name(), "_")[0])
			process(file)
		} else {
			fmt.Println(file.Name() + " - Is not processable")
		}
	}
}

func process(entry os.DirEntry) {
	info, err := entry.Info()
	check(err, "Failed to process file")

	// Open file
	file, err := os.Open("./" + info.Name())
	check(err, "Failed to open file")

	reader := csv.NewReader(file)

	header, err := reader.Read()
	check(err, "Failed to read csv")

	description, amount := parse_header(header)

	rows, err := reader.ReadAll()

	check(err, "Ran into an issue reading CSV")

	for i, row := range rows {
		// Process all rows. I'd like to do it row by row but let's do it the dumb way first
		if !check_filter(row[description], transaction_filters) {
			continue
		}
		fmt.Printf("#%d | %s | $%s\n", i, row[description], row[amount])
	}
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
