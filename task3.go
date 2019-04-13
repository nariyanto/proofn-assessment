package main

import (
	"fmt"
	"math"
	"strconv"
)

func main() {
	a := 13
	b := -2

	r, diff := divide(a, b)

	fmt.Println("result :" + strconv.Itoa(r))
	fmt.Println("diff :" + strconv.Itoa(diff))
}

func divide(dividend int, divisor int) (int, int) {
	summary := 0
	sign := 1
	if (dividend < 0) || (divisor < 0) {
		sign = -1
	}

	// Handling 0
	if dividend == 0 {
		return 0, 0
	}
	if divisor == 0 {
		return math.MaxInt32, 0
	}

	// Handling negative
	if dividend < 0 {
		dividend = dividend * -1
	}
	if divisor < 0 {
		divisor = divisor * -1
	}

	// Handling quotient
	for dividend >= divisor {
		dividend -= divisor
		summary++
	}
	summary = summary * sign

	return summary, dividend
}
