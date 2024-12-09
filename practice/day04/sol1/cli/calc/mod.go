package calc

import "fmt"

func Mod(a, b int) int {
	if b == 0 {
		fmt.Errorf("second operator cannot be zero")
		return 0
	}
	return a % b
}
