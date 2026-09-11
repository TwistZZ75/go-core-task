package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

func main() {
	originalSlice := randomSlice()
	fmt.Println("Generated slice:")
	fmt.Println(originalSlice)

	newSlice := sliceExample(originalSlice)
	fmt.Println("After sliceExample:")
	fmt.Println(newSlice)

	n := 15
	fmt.Println("Adding element: ", n)
	withAdded := addElements(originalSlice, n)
	fmt.Println("After addElements:")
	fmt.Println(withAdded)
	fmt.Println("originalSlice: ", originalSlice)

	cloned := copySlice(originalSlice)
	fmt.Println("After copySlice:")
	fmt.Println(cloned)
	if len(cloned) > 0 {
		cloned[0] = -12345
	}
	fmt.Println("copy after change: ", cloned)
	fmt.Println("originalSlice: ", originalSlice)

	if len(originalSlice) > 0 {
		idx := 3
		removed := removeElement(originalSlice, idx)
		fmt.Printf("after deleting element with index %d: %v\n", idx, removed)
		fmt.Println("originalSlice: ", originalSlice)
	}
}

func randomSlice() []int {
	len := 10
	rndSlice := make([]int, len)

	for i := 0; i < len; i++ {
		value, err := rand.Int(rand.Reader, big.NewInt(100))
		if err != nil {
			panic(err)
		}
		rndSlice[i] = int(value.Int64())
	}

	return rndSlice
}

func sliceExample(origSlice []int) []int {
	newSlice := make([]int, 0, len(origSlice))

	for _, v := range origSlice {
		if v%2 == 0 {
			newSlice = append(newSlice, v)
		}
	}

	return newSlice
}

func addElements(s []int, v int) []int {
	result := make([]int, len(s), len(s)+1)
	copy(result, s)
	result = append(result, v)
	return result
}

func copySlice(s []int) []int {
	result := make([]int, len(s))
	copy(result, s)
	return result
}

func removeElement(s []int, index int) []int {
	if index < 0 || index >= len(s) {
		return copySlice(s)
	}
	result := make([]int, 0, len(s)-1)
	result = append(result, s[:index]...)
	result = append(result, s[index+1:]...)
	return result
}
