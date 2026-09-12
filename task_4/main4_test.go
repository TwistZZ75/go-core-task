package main

import (
	"reflect"
	"testing"
)

func TestDiff(t *testing.T) {
	cases := []struct {
		name   string
		slice1 []string
		slice2 []string
		want   []string
	}{
		{
			name:   "пример из задания",
			slice1: []string{"apple", "banana", "cherry", "date", "43", "lead", "gno1"},
			slice2: []string{"banana", "date", "fig"},
			want:   []string{"apple", "cherry", "43", "lead", "gno1"},
		},
		{
			name:   "полное пересечение",
			slice1: []string{"a", "b", "c"},
			slice2: []string{"a", "b", "c"},
			want:   []string{},
		},
		{
			name:   "нет пересечения",
			slice1: []string{"a", "b"},
			slice2: []string{"x", "y"},
			want:   []string{"a", "b"},
		},
		{
			name:   "второй слайс пустой",
			slice1: []string{"a", "b"},
			slice2: []string{},
			want:   []string{"a", "b"},
		},
		{
			name:   "первый слайс пустой",
			slice1: []string{},
			slice2: []string{"a", "b"},
			want:   []string{},
		},
		{
			name:   "оба пустые",
			slice1: []string{},
			slice2: []string{},
			want:   []string{},
		},
		{
			name:   "дубликаты в первом слайсе сохраняются",
			slice1: []string{"a", "b", "a", "c", "b"},
			slice2: []string{"c"},
			want:   []string{"a", "b", "a", "b"},
		},
		{
			name:   "дубликаты во втором слайсе",
			slice1: []string{"a", "b", "c"},
			slice2: []string{"b", "b", "b"},
			want:   []string{"a", "c"},
		},
		{
			name:   "регистрозависимость",
			slice1: []string{"Apple", "apple"},
			slice2: []string{"apple"},
			want:   []string{"Apple"},
		},
		{
			name:   "пустые строки",
			slice1: []string{"", "a", ""},
			slice2: []string{"a"},
			want:   []string{"", ""},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := diff(tc.slice1, tc.slice2)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("diff(%v, %v) = %v, ожидалось %v",
					tc.slice1, tc.slice2, got, tc.want)
			}
		})
	}
}

func TestDiff_ReturnsNonNil(t *testing.T) {
	got := diff([]string{"a"}, []string{"a"})
	if got == nil {
		t.Errorf("diff вернул nil, ожидался пустой (не nil) слайс")
	}
}

func TestDiff_PreservesOrder(t *testing.T) {
	slice1 := []string{"z", "y", "x", "w", "v"}
	slice2 := []string{"x"}

	want := []string{"z", "y", "w", "v"}
	got := diff(slice1, slice2)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("порядок не сохранён: got = %v, ожидалось %v", got, want)
	}
}
