package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func timeStamp(date, time string, year int) string {
	result := strings.Split(date, "/") // d/m/yyyy
	var dm [2]int
	for i, value := range result[:2] {
		x, err := strconv.Atoi(value)
		if err != nil {
			log.Fatal(err) // TODO: not Fatal
		}
		dm[i] = x
	}
	return fmt.Sprintf("%d-%02d-%02d %s", year, dm[1], dm[0], time) // 2025-05-01 00:24:00
}

func isLeapYear(year int) bool {
	if year%400 == 0 {
		return true
	} else if year%100 == 0 {
		return false
	} else if year%4 == 0 {
		return true
	}

	return false
}

func fileName(year int) string {
	var fakeYear int
	if isLeapYear(year) {
		fakeYear = 2008
	} else {
		fakeYear = 2007
	}
	csvFile := fmt.Sprintf("hhpc%d.csv", fakeYear)
	return csvFile
}

func readCSV(csvfile string) ([][]string, error) {
	file, err := os.Open(csvfile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';' // change Comma to the delimiter in the file
	reader.Read()      // use Read to remove the first line

	rows, err := reader.ReadAll() // ReadAll to read all the data in the file and return a 2D array of strings [][]string.
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func createData(now time.Time, data Data) {
	year := now.Year()
	csvFile := fileName(year)
	fmt.Printf("CSV file: '%s'\n", csvFile)

	rows, err := readCSV(csvFile)
	if err != nil {
		log.Fatal(err)
	}

	for i := range rows {
		if strings.Join(rows[i][2:], "-") == "?-?-?-?-?-?-" {
			copy(rows[i][2:], []string{"0.214", "0.1", "240.0", "1.4", "0.0", "0.0", "0.0"})
			if i > 0 {
				rows[i][8] = rows[i-1][8] // Sub_metering_3
			}
		}
		key := timeStamp(rows[i][0], rows[i][1], year)
		data[key] = rows[i][2:]
	}
	fmt.Println("The length of the map is", len(data))
}
