package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// Handle error checks
func check(e error, message string) {
	if e != nil {
		panic(e)
	}
}

func main() {
	fmt.Println("I'm doing taxes")
	data, err := os.ReadDir("./")

	check(err, "Failed to read directory.")

	fmt.Println("Searching for CSV\n")

	for _, file := range data {
		if strings.Contains(file.Name(), ".CSV") || strings.Contains(file.Name(), ".csv") {
			fmt.Println("\n\nProcessing - " + file.Name())
			process(file)
		} else {
			fmt.Println(file.Name() + " - Is not processable")
		}
	}
}

func process(entry os.DirEntry) {
	info, err := entry.Info()
	check(err, "Failed to process file")

	file, err := os.Open("./" + info.Name())
	check(err, "Failed to open file")

	reader := csv.NewReader(file)

	header, err := reader.Read()
	check(err, "Failed to read csv")

	description, amount := parse_header(header)

	rows, err := reader.ReadAll()

	check(err, "Ran into an issue reading CSV")

	for i, row := range rows {
		fmt.Printf("#%d | %s | $%s", i, row[description], row[amount])
	}
}

func parse_header(header_row []string) (int, int) {
	var NOT_FOUND = -123
	var description int = NOT_FOUND
	var amount int = NOT_FOUND

	for idx, header := range header_row {
		fmt.Println(header)
		if strings.EqualFold(header, "description") {
			description = idx
		} else if strings.EqualFold(header, "amount") {
			amount = idx
		}
	}

	if description == NOT_FOUND || amount == NOT_FOUND || true {
		panic("Couldn't find a column for descriptions or amounts.")
	}
	return description, amount
}
