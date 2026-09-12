package main

import "fmt"

func main() {
	slice1 := []int{65, 3, 58, 678, 64}
	slice2 := []int{64, 2, 3, 43}

	fmt.Println("slice 1: ", slice1)
	fmt.Println("slice 2: ", slice2)
	ok, result := intersection(slice1, slice2)
	fmt.Println("difference: ", ok, result)
}

func intersection(slice1, slice2 []int) (bool, []int) {
	if len(slice1) == 0 || len(slice2) == 0 {
		return false, []int{}
	}
	result := make([]int, 0, len(slice1))

	excludeMap := make(map[int]struct{})
	// пустая структура, потому что у нас мапа без значений,
	// а пуста структура не занимает места
	for _, val := range slice1 {
		excludeMap[val] = struct{}{}
	}

	for _, val := range slice2 {
		if _, excluded := excludeMap[val]; !excluded {
			continue
		}
		result = append(result, val)
	}

	if len(result) == 0 {
		return false, []int{}
	}
	return true, result

}
