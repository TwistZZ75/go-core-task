package main

import (
	"context"
	"testing"
	"time"
)

// Вспомогательные функции

// readN читает до n значений из канала с общим таймаутом.
// Возвращает прочитанный срез и флаг: успели ли прочитать все n.
func readN(ch <-chan int, n int, timeout time.Duration) ([]int, bool) {
	got := make([]int, 0, n)
	deadline := time.After(timeout)
	for i := 0; i < n; i++ {
		select {
		case v, ok := <-ch:
			if !ok {
				return got, false
			}
			got = append(got, v)
		case <-deadline:
			return got, false
		}
	}
	return got, true
}

// Базовые проверки

func TestRndFunc_ProducesValues(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := rndFunc(ctx)

	got, ok := readN(ch, 5, time.Second)
	if !ok {
		t.Fatalf("не удалось прочитать 5 значений, получено: %v", got)
	}
	if len(got) != 5 {
		t.Errorf("прочитано %d значений, ожидалось 5", len(got))
	}
}

func TestRndFunc_WithinRange(t *testing.T) {
	const min, max = 0, 99

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch := rndFunc(ctx)

	got, ok := readN(ch, 200, 2*time.Second)
	if !ok {
		t.Fatalf("не удалось прочитать 200 значений, получено %d", len(got))
	}
	for i, v := range got {
		if v < min || v > max {
			t.Errorf("значение[%d] = %d вне диапазона [%d, %d]", i, v, min, max)
		}
	}
}

// Закрытие канала

func TestRndFunc_ClosesOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ch := rndFunc(ctx)

	// Читаем пару значений, затем отменяем.
	for i := 0; i < 3; i++ {
		select {
		case <-ch:
		case <-time.After(time.Second):
			t.Fatalf("не дождались значения #%d", i)
		}
	}

	cancel()

	// Канал должен закрыться, читаем до закрытия с таймаутом.
	deadline := time.After(time.Second)
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return // успех
			}
		case <-deadline:
			t.Fatalf("канал не закрылся после отмены контекста")
		}
	}
}

func TestRndFunc_ClosesOnTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	ch := rndFunc(ctx)

	// Читаем всё, что приходит, до закрытия канала.
	deadline := time.After(3 * time.Second)
	count := 0
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				if count == 0 {
					t.Errorf("канал закрылся, не выдав ни одного значения")
				}
				return
			}
			count++
		case <-deadline:
			t.Fatalf("канал не закрылся после истечения контекста")
		}
	}
}

func TestRndFunc_CancelBeforeRead(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ch := rndFunc(ctx)

	cancel() // отменяем сразу, не читая

	// Канал должен закрыться, даже если никто не читал.
	deadline := time.After(time.Second)
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return
			}
		case <-deadline:
			t.Fatalf("канал не закрылся после отмены контекста")
		}
	}
}
