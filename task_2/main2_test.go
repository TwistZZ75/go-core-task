package main

import (
	"reflect"
	"testing"
)

func TestRandomSlice_Length(t *testing.T) {
	s := randomSlice()
	if len(s) != 10 {
		t.Errorf("длина = %d, ожидалось 10", len(s))
	}
}

func TestSliceExample(t *testing.T) {
	cases := []struct {
		name string
		in   []int
		want []int
	}{
		{"смешанные", []int{1, 2, 3, 4, 5, 6}, []int{2, 4, 6}},
		{"только чётные", []int{2, 4, 6}, []int{2, 4, 6}},
		{"только нечётные", []int{1, 3, 5}, []int{}},
		{"пустой", []int{}, []int{}},
		{"нули и отрицательные", []int{0, -2, -3, 4}, []int{0, -2, 4}},
		{"только нули", []int{0, 0, 0}, []int{0, 0, 0}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := sliceExample(tc.in)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("sliceExample(%v) = %v, ожидалось %v", tc.in, got, tc.want)
			}
		})
	}
}

func TestSliceExample_EmptyResultNotNil(t *testing.T) {
	// результат должен быть всегда не nil,
	// чтобы его можно было безопасно передавать дальше.
	got := sliceExample([]int{1, 3, 5})
	if got == nil {
		t.Errorf("sliceExample вернул nil для непустого входного слайса")
	}
	if len(got) != 0 {
		t.Errorf("длина = %d, ожидалось 0", len(got))
	}
}

func TestSliceExample_DoesNotMutateInput(t *testing.T) {
	in := []int{1, 2, 3, 4, 5}
	original := append([]int(nil), in...)
	_ = sliceExample(in)
	if !reflect.DeepEqual(in, original) {
		t.Errorf("входной слайс изменён: %v, ожидалось %v", in, original)
	}
}

func TestAddElements(t *testing.T) {
	in := []int{1, 2, 3}
	got := addElements(in, 4)
	want := []int{1, 2, 3, 4}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("addElements(%v, 4) = %v, ожидалось %v", in, got, want)
	}
}

func TestAddElements_Empty(t *testing.T) {
	got := addElements([]int{}, 7)
	if !reflect.DeepEqual(got, []int{7}) {
		t.Errorf("addElements([], 7) = %v, ожидалось [7]", got)
	}
}

func TestAddElements_NegativeValue(t *testing.T) {
	got := addElements([]int{1, 2}, -99)
	if !reflect.DeepEqual(got, []int{1, 2, -99}) {
		t.Errorf("addElements([1 2], -99) = %v, ожидалось [1 2 -99]", got)
	}
}

func TestAddElements_DoesNotMutateInput(t *testing.T) {
	in := []int{1, 2, 3}
	original := append([]int(nil), in...)

	_ = addElements(in, 999)

	if !reflect.DeepEqual(in, original) {
		t.Errorf("входной слайс изменён: %v, ожидалось %v", in, original)
	}
}

func TestAddElements_DoesNotMutateInputWithCapacity(t *testing.T) {
	// Даже если у исходного слайса есть свободная ёмкость,
	// append не должен писать в его массив.
	base := make([]int, 3, 10)
	base[0], base[1], base[2] = 1, 2, 3
	original := append([]int(nil), base...)

	_ = addElements(base, 42)

	if !reflect.DeepEqual(base, original) {
		t.Errorf("входной слайс изменён: %v, ожидалось %v", base, original)
	}
}

func TestCopySlice_Equal(t *testing.T) {
	in := []int{5, 6, 7, 8}
	got := copySlice(in)
	if !reflect.DeepEqual(got, in) {
		t.Errorf("copySlice(%v) = %v, ожидалось равенство", in, got)
	}
}

func TestCopySlice_Independent(t *testing.T) {
	in := []int{1, 2, 3}
	cp := copySlice(in)

	cp[0] = 100
	cp[1] = 200

	if in[0] != 1 || in[1] != 2 {
		t.Errorf("изменение копии повлияло на исходный: %v", in)
	}
}

func TestCopySlice_Empty(t *testing.T) {
	got := copySlice([]int{})
	if len(got) != 0 {
		t.Errorf("ожидался пустой слайс, получено %v", got)
	}
}

func TestRemoveElement(t *testing.T) {
	cases := []struct {
		name  string
		in    []int
		index int
		want  []int
	}{
		{"середина", []int{1, 2, 3, 4, 5}, 2, []int{1, 2, 4, 5}},
		{"первый", []int{1, 2, 3}, 0, []int{2, 3}},
		{"последний", []int{1, 2, 3}, 2, []int{1, 2}},
		{"единственный", []int{9}, 0, []int{}},
		{"индекс вне диапазона (большой)", []int{1, 2, 3}, 10, []int{1, 2, 3}},
		{"отрицательный индекс", []int{1, 2, 3}, -1, []int{1, 2, 3}},
		{"пустой", []int{}, 0, []int{}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := removeElement(tc.in, tc.index)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("removeElement(%v, %d) = %v, ожидалось %v",
					tc.in, tc.index, got, tc.want)
			}
		})
	}
}
