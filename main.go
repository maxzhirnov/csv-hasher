package main

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"encoding/csv"
	"fmt"
	"hash"
	"io"
	"os"

	"github.com/360EntSecGroup-Skylar/excelize"
)

func hashData(data []string, h hash.Hash) []string {
	hashedData := make([]string, len(data))
	for i, val := range data {
		h.Reset()
		io.WriteString(h, val)
		if val == "" {
			hashedData[i] = ""
			continue
		}
		hashedData[i] = fmt.Sprintf("%x", h.Sum(nil))
	}
	return hashedData
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("csv-hasher <input_file.xlsx> <hash_algorithm> (sha256 or md5)")
		os.Exit(1)
	}

	hashAlgorithmArg := os.Args[2]
	if hashAlgorithmArg != "sha256" && hashAlgorithmArg != "md5" {
		fmt.Println("Invalid hash algorithm. Use sha256 or md5")
		os.Exit(1)
	}

	var hashAlgorithm hash.Hash

	switch hashAlgorithmArg {
	case "sha256":
		hashAlgorithm = sha256.New()
	case "md5":
		hashAlgorithm = md5.New()
	}

	outputName := fmt.Sprintf("hashed-output-%s.csv", hashAlgorithmArg)

	xlFile, err := excelize.OpenFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Create a new file for writing hashed data
	csvFile, err := os.Create(outputName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(bufio.NewWriter(csvFile))
	defer writer.Flush()

	// Iterate over each sheet in the XLSX file
	for _, sheetName := range xlFile.GetSheetMap() {
		// Iterate over each row in the sheet
		for rowIndex, rowValues := range xlFile.GetRows(sheetName) {
			var data []string
			// Iterate over each cell in the row
			for _, cellValue := range rowValues {
				data = append(data, cellValue)
			}
			// Hash each cell value (excluding the first cell) and write to the CSV file
			if rowIndex == 0 {
				err := writer.Write(data)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			} else {
				hashedData := hashData(data, hashAlgorithm)
				err := writer.Write(hashedData)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
		}
	}
	fmt.Println("Done")
}
