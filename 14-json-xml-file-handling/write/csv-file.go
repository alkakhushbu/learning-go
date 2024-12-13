package main

import (
	"encoding/csv"
	"log"
	"os"
)

func main() {
	filepath := "files/data.csv"
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Println("Error opening file:", filepath, "Error:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	data := []string{"File", "Operation", "Best", "Practices"}
	err = writer.Write(data)
	if err != nil {
		log.Println("Error in writing data to file:", filepath, "Error:", err)
		return
	}
	log.Println("Data saved successfully!")
}
