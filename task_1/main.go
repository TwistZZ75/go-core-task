package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
)

func main() {
	decInt, octInt, hexInt, fl, str, b, cmlx := createVariables()

	fmt.Println("Переменные и их типы")
	defineType(decInt, octInt, hexInt, fl, str, b, cmlx)

	fmt.Println("\nОбъединённая строка")
	commonString := createString(decInt, octInt, hexInt, fl, str, b, cmlx)

	rune_slice := convertToRune(commonString)

	saltStr := "go-2026"

	hash := hashRunesWithSalt(rune_slice, saltStr)

	fmt.Printf("\nSHA256 (с солью %q в середине)\n", saltStr)
	fmt.Println(hash)
}

func createVariables() (decInt int, octInt int, hexInt int,
	fl float64, str string, b bool, cmlx complex64) {
	decInt = 54
	octInt = 0o27
	hexInt = 0x22
	fl = 5.17
	str = "first test task"
	b = true
	cmlx = complex64(complex(1.5, 2.5))

	return
}

func defineType(decInt int, octInt int, hexInt int,
	fl float64, str string, b bool, cmlx complex64) {

	fmt.Printf("decInt = %v (тип: %T)\n", decInt, decInt)
	fmt.Printf("octInt = %v (тип: %T)\n", octInt, octInt)
	fmt.Printf("hexInt = %v (тип: %T)\n", hexInt, hexInt)
	fmt.Printf("fl = %v (тип: %T)\n", fl, fl)
	fmt.Printf("str = %v (тип: %T)\n", str, str)
	fmt.Printf("b = %v (тип: %T)\n", b, b)
	fmt.Printf("cmlx = %v (тип: %T)\n", cmlx, cmlx)
}

func createString(decInt int, octInt int, hexInt int,
	fl float64, str string, b bool, cmlx complex64) string {

	commonStr := strconv.Itoa(decInt) + strconv.Itoa(octInt) +
		strconv.Itoa(hexInt) + strconv.FormatFloat(fl, 'f', -1, 64) +
		str + strconv.FormatBool(b) +
		strconv.FormatComplex(complex128(cmlx), 'f', -1, 64)

	fmt.Println(commonStr)
	return commonStr
}

func convertToRune(str string) []rune {
	return []rune(str)
}

func hashRunesWithSalt(r []rune, salt string) string {
	mid := len(r) / 2
	saltRunes := []rune(salt)

	combined := make([]rune, 0, len(r)+len(saltRunes))
	combined = append(combined, r[:mid]...)
	combined = append(combined, saltRunes...)
	combined = append(combined, r[mid:]...)

	sum := sha256.Sum256([]byte(string(combined)))
	return hex.EncodeToString(sum[:])
}
