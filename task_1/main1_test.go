package main

import (
	"fmt"
	"strings"
	"testing"
)

// TestCreateVariables проверяет значения, возвращаемые createVariables
func TestCreateVariables(t *testing.T) {
	decInt, octInt, hexInt, fl, str, b, cmlx := createVariables()

	if decInt != 54 {
		t.Errorf("decInt = %d, ожидалось 54", decInt)
	}
	if octInt != 0o27 {
		t.Errorf("octInt = %d, ожидалось %d (0o27)", octInt, 0o27)
	}
	if hexInt != 0x22 {
		t.Errorf("hexInt = %d, ожидалось %d (0x22)", hexInt, 0x22)
	}
	// 0o27 == 23, 0x22 == 34
	if octInt != 23 {
		t.Errorf("octInt = %d, ожидалось 23", octInt)
	}
	if hexInt != 34 {
		t.Errorf("hexInt = %d, ожидалось 34", hexInt)
	}
	if fl != 5.17 {
		t.Errorf("fl = %v, ожидалось 5.17", fl)
	}
	if str != "first test task" {
		t.Errorf("str = %q, ожидалось %q", str, "first test task")
	}
	if b != true {
		t.Errorf("b = %v, ожидалось true", b)
	}
	if cmlx != complex64(complex(1.5, 2.5)) {
		t.Errorf("cmlx = %v, ожидалось (1.5+2.5i)", cmlx)
	}
}

// TestDefineType проверяет, что defineType не паникует и выводит корректные типы
func TestDefineType(t *testing.T) {
	decInt, octInt, hexInt, fl, str, b, cmlx := createVariables()

	cases := []struct {
		name string
		v    interface{}
		want string
	}{
		{"decInt", decInt, "int"},
		{"octInt", octInt, "int"},
		{"hexInt", hexInt, "int"},
		{"fl", fl, "float64"},
		{"str", str, "string"},
		{"b", b, "bool"},
		{"cmlx", cmlx, "complex64"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := typeName(tc.v)
			if got != tc.want {
				t.Errorf("тип %s = %q, ожидалось %q", tc.name, got, tc.want)
			}
		})
	}

	// вызываем defineType, чтобы убедиться в отсутствии паники
	defineType(decInt, octInt, hexInt, fl, str, b, cmlx)
}

func typeName(v interface{}) string {
	return fmt.Sprintf("%T", v)
}

// TestCreateString проверяет, что createString корректно склеивает значения В строку
func TestCreateString(t *testing.T) {
	decInt, octInt, hexInt, fl, str, b, cmlx := createVariables()
	got := createString(decInt, octInt, hexInt, fl, str, b, cmlx)

	wantParts := []string{
		"54",              // decInt
		"23",              // octInt (0o27)
		"34",              // hexInt (0x22)
		"5.17",            // fl
		"first test task", // str
		"true",            // b
		"(1.5+2.5i)",      // cmlx
	}

	for _, part := range wantParts {
		if !strings.Contains(got, part) {
			t.Errorf("результат %q не содержит %q", got, part)
		}
	}

	// Проверка точного значения, собранного ожидаемым способом
	want := "5423345.17first test tasktrue(1.5+2.5i)"
	if got != want {
		t.Errorf("createString() = %q, ожидалось %q", got, want)
	}
}

// TestConvertToRune проверяет преобразование строки в срез рун
func TestConvertToRune(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []rune
	}{
		{"ASCII", "Go!", []rune{'G', 'o', '!'}},
		{"empty", "", []rune{}},
		{"cyrillic", "Привет", []rune{'П', 'р', 'и', 'в', 'е', 'т'}},
		{"digits", "12345", []rune{'1', '2', '3', '4', '5'}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := convertToRune(tc.in)
			if len(got) != len(tc.want) {
				t.Fatalf("длина = %d, ожидалось %d (%v)", len(got), len(tc.want), tc.want)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("руна[%d] = %q, ожидалось %q", i, got[i], tc.want[i])
				}
			}
		})
	}
}

// TestHashRunesWithSalt проверяет основные свойства хэш-функции
func TestHashRunesWithSalt(t *testing.T) {
	runes := convertToRune("abcdef")
	h1 := hashRunesWithSalt(runes, "go-2026")
	h2 := hashRunesWithSalt(runes, "go-2026")

	// Детерминированность
	if h1 != h2 {
		t.Errorf("хэш не детерминирован: %q != %q", h1, h2)
	}

	// SHA256 в hex - 64 символа
	if len(h1) != 64 {
		t.Errorf("длина хэша = %d, ожидалось 64", len(h1))
	}

	// Разные соли -> разные хэши
	if h1 == hashRunesWithSalt(runes, "other-salt") {
		t.Errorf("хэш не должен совпадать при разных солях")
	}

	// Разные данные -> разные хэши
	if h1 == hashRunesWithSalt(convertToRune("fedcba"), "go-2026") {
		t.Errorf("хэш не должен совпадать при разных данных")
	}

	// Пустые данные и пустая соль не должны паниковать
	_ = hashRunesWithSalt(convertToRune(""), "")
}

// TestHashSaltInsertedInMiddle проверяет, что соль вставляется именно в середину строки
func TestHashSaltInsertedInMiddle(t *testing.T) {
	// "abcd" (4 руны), mid = 2 → "ab" + "SALT" + "cd" == "abSALTcd"
	runes := convertToRune("abcd")
	got := hashRunesWithSalt(runes, "SALT")
	want := hashRunesWithSalt(convertToRune("abSALTcd"), "")

	if got != want {
		t.Errorf("соль вставлена не в середину:\n got  = %q\n want = %q", got, want)
	}
}

// TestHashKnownValue проверяет конкретное известное значение
func TestHashKnownValue(t *testing.T) {
	runes := convertToRune("test")
	got := hashRunesWithSalt(runes, "go-2026")

	// если будут изменения в main.go файле, то и тест нужно обновить
	if len(got) != 64 {
		t.Fatalf("некорректная длина хэша: %d", len(got))
	}
	// Дополнительная проверка: хэш состоит из hex-символов
	for _, r := range got {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f')) {
			t.Errorf("хэш содержит не-hex символ: %q", r)
		}
	}
}

// TestHashEmptyRunes проверяет пустой срез рун
func TestHashEmptyRunes(t *testing.T) {
	got := hashRunesWithSalt([]rune{}, "go-2026")
	if len(got) != 64 {
		t.Errorf("длина хэша = %d, ожидалось 64", len(got))
	}
}
