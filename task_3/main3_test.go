package main

import (
	"reflect"
	"testing"
)

func TestNewStringIntMap_Empty(t *testing.T) {
	m := NewStringIntMap()
	if m.Len() != 0 {
		t.Errorf("новая карта не пуста: len = %d", m.Len())
	}
	if m.Exists("any") {
		t.Errorf("в новой карте не должно быть ключей")
	}
}

func TestAdd_NewKey(t *testing.T) {
	m := NewStringIntMap()
	m.Add("a", 1)

	if !m.Exists("a") {
		t.Fatalf("ключ %q не найден после Add", "a")
	}
	v, ok := m.Get("a")
	if !ok || v != 1 {
		t.Errorf("Get(a) = (%d, %v), ожидалось (1, true)", v, ok)
	}
	if m.Len() != 1 {
		t.Errorf("len = %d, ожидалось 1", m.Len())
	}
}

func TestAdd_Overwrite(t *testing.T) {
	m := NewStringIntMap()
	m.Add("a", 1)
	m.Add("a", 42)

	v, ok := m.Get("a")
	if !ok || v != 42 {
		t.Errorf("Get(a) = (%d, %v), ожидалось (42, true)", v, ok)
	}
	if m.Len() != 1 {
		t.Errorf("после перезаписи len = %d, ожидалось 1", m.Len())
	}
}
func TestAdd_EmptyKey(t *testing.T) {
	m := NewStringIntMap()
	m.Add("", 5)

	if !m.Exists("") {
		t.Errorf("пустой ключ должен быть допустим")
	}
	if v, ok := m.Get(""); !ok || v != 5 {
		t.Errorf("Get(\"\") = (%d, %v), ожидалось (5, true)", v, ok)
	}
}
func TestRemove_Existing(t *testing.T) {
	m := NewStringIntMap()
	m.Add("a", 1)
	m.Add("b", 2)

	m.Remove("a")

	if m.Exists("a") {
		t.Errorf("ключ %q должен быть удалён", "a")
	}
	if !m.Exists("b") {
		t.Errorf("ключ %q не должен был удалиться", "b")
	}
	if m.Len() != 1 {
		t.Errorf("len = %d, ожидалось 1", m.Len())
	}
}

func TestRemove_NonExisting(t *testing.T) {
	m := NewStringIntMap()
	m.Add("a", 1)

	// Удаление несуществующего ключа не должно паниковать и не менять карту
	m.Remove("missing")

	if m.Len() != 1 {
		t.Errorf("len = %d, ожидалось 1", m.Len())
	}
	if !m.Exists("a") {
		t.Errorf("существующий ключ %q был затронут", "a")
	}
}
func TestRemove_FromEmpty(t *testing.T) {
	m := NewStringIntMap()
	m.Remove("nope") // не должно паниковать
	if m.Len() != 0 {
		t.Errorf("len = %d, ожидалось 0", m.Len())
	}
}
func TestCopy_Equal(t *testing.T) {
	m := NewStringIntMap()
	m.Add("a", 1)
	m.Add("b", 2)
	m.Add("c", 3)

	cp := m.Copy()
	if !reflect.DeepEqual(cp, map[string]int{"a": 1, "b": 2, "c": 3}) {
		t.Errorf("Copy() = %v, ожидалось {a:1 b:2 c:3}", cp)
	}
}

func TestCopy_Independent(t *testing.T) {
	m := NewStringIntMap()
	m.Add("a", 1)

	cp := m.Copy()
	cp["a"] = 999
	cp["new"] = 42

	// Изменения в копии не должны затрагивать оригинал
	if v, _ := m.Get("a"); v != 1 {
		t.Errorf("изменение копии повлияло на оригинал: a = %d", v)
	}
	if m.Exists("new") {
		t.Errorf("добавление в копию попало в оригинал")
	}
}

func TestCopy_Empty(t *testing.T) {
	m := NewStringIntMap()
	cp := m.Copy()

	if cp == nil {
		t.Fatalf("Copy() не должен возвращать nil")
	}
	if len(cp) != 0 {
		t.Errorf("копия пустой карты должна быть пустой, получено %v", cp)
	}
}

func TestExists(t *testing.T) {
	m := NewStringIntMap()
	m.Add("yes", 1)

	if !m.Exists("yes") {
		t.Errorf("Exists(yes) = false, ожидалось true")
	}
	if m.Exists("no") {
		t.Errorf("Exists(no) = true, ожидалось false")
	}
}
func TestGet_Existing(t *testing.T) {
	m := NewStringIntMap()
	m.Add("k", 42)

	v, ok := m.Get("k")
	if !ok {
		t.Errorf("Get(k): ok = false, ожидалось true")
	}
	if v != 42 {
		t.Errorf("Get(k) = %d, ожидалось 42", v)
	}
}

func TestGet_Missing(t *testing.T) {
	m := NewStringIntMap()

	v, ok := m.Get("missing")
	if ok {
		t.Errorf("Get(missing): ok = true, ожидалось false")
	}
	if v != 0 {
		t.Errorf("Get(missing) = %d, ожидалось нулевое значение 0", v)
	}
}

func TestGet_AfterRemove(t *testing.T) {
	m := NewStringIntMap()
	m.Add("k", 42)
	m.Remove("k")

	if _, ok := m.Get("k"); ok {
		t.Errorf("после Remove ключ всё ещё доступен")
	}
}
