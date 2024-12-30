package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	fmt.Println(RupeesToPaise("99.00 "))
	fmt.Println(RupeesToPaise("99.999"))
	fmt.Println(RupeesToPaise("099.9"))
	fmt.Println(RupeesToPaise("099.99"))
	fmt.Println(RupeesToPaise("99.9v"))
	fmt.Println(RupeesToPaise("99."))
	fmt.Println(RupeesToPaise("95"))
	fmt.Println(RupeesToPaise("A9.99"))
	fmt.Println(RupeesToPaise("425.5"))
	fmt.Println(RupeesToPaise("-425.5"))
	fmt.Println(RupeesToPaise("425.-5"))
	fmt.Println(RupeesToPaise("425.555"))
}

func RupeesToPaise(priceStr string) (uint64, error) {
	fmt.Println("Input price:", priceStr)
	//trim extra space from price
	priceStr = strings.Trim(priceStr, " ")

	//split the price based by dot(.)
	prices := strings.Split(priceStr, ".")
	var rupee, paisa uint64
	if len(prices) == 0 || len(prices) > 2 {
		return 0, fmt.Errorf("invalid price, empty price field or more than one dot(.)")
	}

	rupee, err := strconv.ParseUint(prices[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid price, not a valid number")
	}

	if len(prices) == 2 {
		paisa, err = strconv.ParseUint(prices[1], 10, 64)
		if err != nil || paisa > 99 {
			return 0, fmt.Errorf("invalid price, please provide price in valid format")
		}

		// append 0 if paisa part has only one digit
		// e.g INR 99.2 => Convert it to 9900 + 20 = 9920
		if paisa < 10 {
			paisa *= 10
		}
	}
	return rupee*100 + paisa, nil
}
