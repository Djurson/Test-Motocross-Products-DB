package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func foo() {
	file, err := os.Open("test.csv")
	if err != nil {
		fmt.Println("Fel vid öppning av fil: ", err)
	}
	reader := csv.NewReader(file)
	reader.Comma = ';'

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Fel vid läsning av CSV:", err)
		return
	}

	if len(records) < 1 {
		fmt.Println("CSV-filen är tom")
		return
	}

	readCSVRecords(records)
}

func readCSVRecords(records [][]string) {
	for _, row := range records {
		subCategoryName := row[0]
		brandName := row[4]
		modelName := row[5]
		modYears := row[6]
		productCode := "KT" + row[9]
		productName := row[10]

		startYear, endYear := parseModelYears(modYears)

		fmt.Println("Category:", subCategoryName, "Bike brand name:", brandName,
			"Bike model:", modelName, "Model year:", startYear, "-", endYear,
			"Product code:", productCode, "Product name:", productName, "\n")

	}
}
