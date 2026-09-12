package main

import "fmt"

type StringIntMap struct {
	data map[string]int
}

func NewStringIntMap() *StringIntMap {
	return &StringIntMap{
		data: make(map[string]int)}
}

func (m *StringIntMap) Add(key string, value int) {
	if m.data == nil {
		m.data = make(map[string]int)
	}
	m.data[key] = value
}

func (m *StringIntMap) Remove(key string) {
	delete(m.data, key)
}

func (m *StringIntMap) Copy() map[string]int {
	copy := make(map[string]int, len(m.data))
	for k, v := range m.data {
		copy[k] = v
	}
	return copy
}

func (m *StringIntMap) Exists(key string) bool {
	_, ok := m.data[key]
	return ok
}

func (m *StringIntMap) Get(key string) (int, bool) {
	v, ok := m.data[key]
	return v, ok
}

// для тестов
func (m *StringIntMap) Len() int {
	return len(m.data)
}

func main() {
	m := NewStringIntMap()

	m.Add("one", 1)
	m.Add("fourty two", 42)
	m.Add("sixty nine", 69)
	m.Add("one hundred", 100)

	fmt.Println(m.data)

	if v, ok := m.Get("two"); ok {
		fmt.Println("Get(two) =", v)
	}

	fmt.Println("Exists(four) =", m.Exists("four"))

	cp := m.Copy()
	cp["four"] = 4
	fmt.Println("После изменения копии:")
	fmt.Println("оригинал: Exists(four) =", m.Exists("four"))
	_, ok := cp["four"]
	fmt.Println("копия: Exists(four) =", ok)

	m.Remove("one")
	fmt.Println("После Remove(one): Exists(one) =", m.Exists("one"))
}
