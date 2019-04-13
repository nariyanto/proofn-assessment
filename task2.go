package main

import (
	"fmt"
)

func main() {
	input := []int{1, -3, 5, -2, 9, -8, -6, 4}
	length := len(input)

	result := maxContiguousArraySum(input, length)

	fmt.Println(result)
}

func maxContiguousArraySum(arr []int, size int) ([]int) {
	currMax := -1
	max := 0
	start := 0
	end := 0
	s := 0

	var result []int

	for i := 0; i < size; i++ {
		max = max + arr[i]
		if currMax < max {
			currMax = max
			start = s
			end = i
		}

		if max < 0 {
			max = 0
			s = i + 1
		}
	}

	for j := start; j <= end; j++ {
		result = append(result, arr[j])
	}

	return result
}
