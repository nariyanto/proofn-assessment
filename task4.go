package main

import (
	"fmt"
	"sort"
)

func main() {
	// input := []int{4, 3, 2, 2, 1, 1}
	// amount := 4

	input := []int{50, 50, 50, 50, 25, 25, 25, 25, 25, 25}
	amount := 300
	
	result := getNumberOfWays(amount, input)

	fmt.Println(result)
}

func getNumberOfWays(n int, coins []int) int {
	sort.Ints(coins)
	
	result := 0
	duplicates := make(map[int][]int)

	for i := 0; i < len(coins); i++ {
		
		comb := []int{}
		counter := 0
		x := i
		amount := n

		for amount > 0 && counter < len(coins) {
			if x == (len(coins)) {
				x = 0
			}

			checkOverlap := amount - coins[x]
			
			if checkOverlap >= 0 && (checkOverlap >= coins[x] || checkOverlap == 0) || i == x {
				amount = amount - coins[x]
				comb = append(comb, coins[x])
			}
			
			if amount == 0 {
				sort.Ints(comb)
				if !isExists(duplicates, comb) {
					duplicates[result] = comb
					result++
				}
				break;
			}
			
			counter++
			x++
		}
	}

	// fmt.Println(duplicates);

	return result
}

func isExists(arr map[int][]int, b []int) bool  {
	if len(arr) == 0 {
		return false
	}

	// exist
	for _,a := range arr {
		if len(a) == len(b) {
			for x := range a {
				if a[x] == b[x] && x == len(a)-1 {
					return true
				}
			}	
		}
	}

	return false
}