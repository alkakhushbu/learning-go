package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func main() {
	filepath := "read/data.csv"
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error in opening file:", file, "Error:", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error in reading csv file:", file, "Error:", err)
		return
	}
	// for _, record := range records {
	// 	fmt.Println(record)
	// }
	fmt.Println(records)

}
