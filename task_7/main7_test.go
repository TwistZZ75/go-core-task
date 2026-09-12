package main

import (
	"context"
	"sort"
	"sync"
	"testing"
	"time"
)

// makeChan создаёт канал, который отдаёт переданные значения и закрывается
func makeChan(values ...int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for _, v := range values {
			ch <- v
		}
	}()
	return ch
}

// collect читает все значения из канала с таймаутом
func collect(t *testing.T, ch <-chan int, timeout time.Duration) ([]int, bool) {
	t.Helper()
	var got []int
	deadline := time.After(timeout)
	for {
		select {
		case v, ok := <-ch:
			if !ok {
				return got, true
			}
			got = append(got, v)
		case <-deadline:
			return got, false
		}
	}
}

func sortInts(s []int) []int {
	out := append([]int(nil), s...)
	sort.Ints(out)
	return out
}

// Базовые тесты

func TestMerge_SingleChannel(t *testing.T) {
	ctx := context.Background()
	in := makeChan(1, 2, 3)

	out := Merge(ctx, in)
	got, ok := collect(t, out, time.Second)
	if !ok {
		t.Fatalf("канал не закрылся, получено: %v", got)
	}

	want := []int{1, 2, 3}
	if !equalInts(sortInts(got), want) {
		t.Errorf("got = %v, want = %v", got, want)
	}
}

func TestMerge_MultipleChannels(t *testing.T) {
	ctx := context.Background()
	ch1 := makeChan(1, 2, 3)
	ch2 := makeChan(10, 20, 30)
	ch3 := makeChan(100, 200)

	out := Merge(ctx, ch1, ch2, ch3)
	got, ok := collect(t, out, time.Second)
	if !ok {
		t.Fatalf("канал не закрылся, получено: %v", got)
	}

	want := []int{1, 2, 3, 10, 20, 30, 100, 200}
	if !equalInts(sortInts(got), want) {
		t.Errorf("got = %v, want = %v", sortInts(got), want)
	}
}

func TestMerge_EmptyChannels(t *testing.T) {
	ctx := context.Background()
	ch1 := makeChan()
	ch2 := makeChan()
	ch3 := makeChan()

	out := Merge(ctx, ch1, ch2, ch3)
	got, ok := collect(t, out, time.Second)
	if !ok {
		t.Fatalf("канал не закрылся")
	}
	if len(got) != 0 {
		t.Errorf("ожидался пустой результат, получено: %v", got)
	}
}

func TestMerge_NoChannels(t *testing.T) {
	ctx := context.Background()

	out := Merge[int](ctx)

	// Канал должен быть уже закрыт
	select {
	case _, ok := <-out:
		if ok {
			t.Errorf("ожидался закрытый канал")
		}
	case <-time.After(100 * time.Millisecond):
		t.Errorf("канал не закрыт при отсутствии входов")
	}
}

func TestMerge_SomeEmptyChannels(t *testing.T) {
	ctx := context.Background()
	ch1 := makeChan()
	ch2 := makeChan(1, 2, 3)
	ch3 := makeChan()

	out := Merge(ctx, ch1, ch2, ch3)
	got, ok := collect(t, out, time.Second)
	if !ok {
		t.Fatalf("канал не закрылся")
	}
	if !equalInts(got, []int{1, 2, 3}) {
		t.Errorf("got = %v, want = [1 2 3]", got)
	}
}

// Отмена контекста

func TestMerge_ContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	infinite := func() <-chan int {
		ch := make(chan int)
		go func() {
			defer close(ch)
			i := 0
			for {
				select {
				case <-ctx.Done():
					return
				case ch <- i:
					i++
				}
			}
		}()
		return ch
	}

	out := Merge(ctx, infinite(), infinite(), infinite())

	for i := 0; i < 3; i++ {
		select {
		case <-out:
		case <-time.After(time.Second):
			t.Fatalf("не дождались значения #%d", i)
		}
	}

	cancel()

	// После отмены выходной канал должен закрыться
	deadline := time.After(time.Second)
	for {
		select {
		case _, ok := <-out:
			if !ok {
				return // успех
			}
		case <-deadline:
			t.Fatalf("канал не закрылся после отмены контекста")
		}
	}
}

// Утечки горутин

func TestMerge_NoGoroutineLeak(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	infinite := func() <-chan int {
		ch := make(chan int)
		go func() {
			defer close(ch)
			i := 0
			for {
				select {
				case <-ctx.Done():
					return
				case ch <- i:
					i++
				}
			}
		}()
		return ch
	}

	out := Merge(ctx, infinite(), infinite())

	// Читаем пару значений и отменяем.
	<-out
	<-out
	cancel()

	select {
	case <-out:
		// дренаж
	case <-time.After(time.Second):
		t.Fatalf("выходной канал не закрылся — возможна утечка горутин")
	}
}

// Конкурентное чтение

func TestMerge_ConcurrentReaders(t *testing.T) {
	ctx := context.Background()

	ch1 := makeChan(1, 2, 3, 4, 5)
	ch2 := makeChan(6, 7, 8, 9, 10)
	out := Merge(ctx, ch1, ch2)

	var mu sync.Mutex
	seen := make(map[int]bool)
	var wg sync.WaitGroup

	// Несколько горутин читают из одного выходного канала.
	for r := 0; r < 3; r++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range out {
				mu.Lock()
				if seen[v] {
					t.Errorf("значение %d прочитано дважды", v)
				}
				seen[v] = true
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	if len(seen) != 10 {
		t.Errorf("прочитано %d уникальных значений, ожидалось 10", len(seen))
	}
}

func equalInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
