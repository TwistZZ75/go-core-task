package main

import "fmt"

func main() {
	slice1 := []string{"apple", "banana", "cherry", "date", "43", "lead", "gno1"}
	slice2 := []string{"banana", "date", "fig"}

	fmt.Println("slice 1: ", slice1)
	fmt.Println("slice 2: ", slice2)
	fmt.Println("difference: ", diff(slice1, slice2))
}

func diff(slice1, slice2 []string) []string {
	result := make([]string, 0, len(slice1))

	excludeMap := make(map[string]struct{})
	// пустая структура, потому что у нас мапа без значений,
	// а пуста структура не занимает места
	for _, val := range slice2 {
		excludeMap[val] = struct{}{}
	}

	for _, val := range slice1 {
		if _, excluded := excludeMap[val]; excluded {
			continue
		}
		result = append(result, val)
	}

	return result

}
