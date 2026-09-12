package main

import (
	"reflect"
	"testing"
)

func TestIntersection(t *testing.T) {
	cases := []struct {
		name     string
		a        []int
		b        []int
		wantBool bool
		wantNums []int
	}{
		{
			name:     "пример из задания",
			a:        []int{65, 3, 58, 678, 64},
			b:        []int{64, 2, 3, 43},
			wantBool: true,
			wantNums: []int{64, 3},
		},
		{
			name:     "полное совпадение",
			a:        []int{1, 2, 3},
			b:        []int{1, 2, 3},
			wantBool: true,
			wantNums: []int{1, 2, 3},
		},
		{
			name:     "нет пересечений",
			a:        []int{1, 2, 3},
			b:        []int{4, 5, 6},
			wantBool: false,
			wantNums: []int{},
		},
		{
			name:     "первый пустой",
			a:        []int{},
			b:        []int{1, 2, 3},
			wantBool: false,
			wantNums: []int{},
		},
		{
			name:     "второй пустой",
			a:        []int{1, 2, 3},
			b:        []int{},
			wantBool: false,
			wantNums: []int{},
		},
		{
			name:     "оба пустые",
			a:        []int{},
			b:        []int{},
			wantBool: false,
			wantNums: []int{},
		},
		{
			name:     "один общий элемент",
			a:        []int{1, 2, 3},
			b:        []int{5, 3, 7},
			wantBool: true,
			wantNums: []int{3},
		},
		{
			name:     "отрицательные значения",
			a:        []int{-1, -2, -3},
			b:        []int{0, -2, 5, -3},
			wantBool: true,
			wantNums: []int{-2, -3},
		},
		{
			name:     "нули",
			a:        []int{0, 1, 2},
			b:        []int{3, 0, 5},
			wantBool: true,
			wantNums: []int{0},
		},
		{
			name:     "порядок соответствует b",
			a:        []int{10, 20, 30, 40},
			b:        []int{40, 20, 10},
			wantBool: true,
			wantNums: []int{40, 20, 10},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotBool, gotNums := intersection(tc.a, tc.b)

			if gotBool != tc.wantBool {
				t.Errorf("bool = %v, ожидалось %v", gotBool, tc.wantBool)
			}
			if !reflect.DeepEqual(gotNums, tc.wantNums) {
				t.Errorf("nums = %v, ожидалось %v", gotNums, tc.wantNums)
			}
		})
	}
}

func TestIntersection_EmptyResultIsNonNil(t *testing.T) {
	// Даже при отсутствии пересечений результат должен быть не nil
	_, nums := intersection([]int{1, 2}, []int{3, 4})
	if nums == nil {
		t.Errorf("при отсутствии пересечений ожидался пустой (не nil) срез")
	}
}

func TestIntersection_BoolMatchesLen(t *testing.T) {
	// bool и длина среза должны быть согласованы между собой.
	cases := []struct {
		a, b []int
	}{
		{[]int{1, 2, 3}, []int{2, 3, 4}},
		{[]int{1, 2}, []int{3, 4}},
		{[]int{}, []int{}},
		{[]int{5}, []int{5}},
		{[]int{5}, []int{6}},
	}
	for _, tc := range cases {
		ok, nums := intersection(tc.a, tc.b)
		if ok && len(nums) == 0 {
			t.Errorf("bool = true, но срез пуст: a=%v b=%v", tc.a, tc.b)
		}
		if !ok && len(nums) != 0 {
			t.Errorf("bool = false, но срез непуст: %v", nums)
		}
	}
}
