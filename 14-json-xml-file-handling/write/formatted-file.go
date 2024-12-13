package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	userName, orderId := "alka", "od123456"
	filepath := "files/" + userName + "_" + orderId + ".pdf"
	file, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	item1, item2, item3 := "Shoes", "Bag", "Kettle"
	price1, price2, price3 := 234, 229, 344
	_, err = fmt.Fprintf(file, "Username: %s\nOrder Number: %s\nItem 1: %s\nPrice 1: %d\nItem 2: %s\nPrice 2: %d\nItem 3: %s\nPrice 3: %d\n", userName, orderId, item1, price1, item2, price2, item3, price3)
	if err != nil {
		log.Fatal("Error in saving file:", file, "Error:", err)
	}
	log.Println("File saved successfully")

}
