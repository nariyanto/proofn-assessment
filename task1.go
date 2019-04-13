package main

import "fmt"

func main() {
	arrInput := []int{9, 10, 8, 4, 5, 3}
	x := 15
	result := hasSumValue(arrInput, x)

	fmt.Println(result)
}

func hasSumValue(params []int, x int) bool {
	for a, i := range params {
		for b, j := range params {
			if a != b && i+j == x {
				return true
			}
		}
	}

	return false
}
