package main

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestWaitGroup_AddZero_WaitReturnsImmediately(t *testing.T) {
	wg := NewSemaWaitGroup()
	wg.Add(0)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// успех
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Wait() не вернулся при нулевом счётчике")
	}
}

func TestWaitGroup_SingleDone(t *testing.T) {
	wg := NewSemaWaitGroup()
	wg.Add(1)

	go func() {
		time.Sleep(50 * time.Millisecond)
		wg.Done()
	}()

	start := time.Now()
	wg.Wait()
	elapsed := time.Since(start)

	if elapsed < 40*time.Millisecond {
		t.Errorf("Wait() вернулся слишком рано: %v", elapsed)
	}
}

func TestWaitGroup_MultipleWorkers(t *testing.T) {
	wg := NewSemaWaitGroup()

	const workers = 100
	var completed atomic.Int64

	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			completed.Add(1)
		}()
	}

	wg.Wait()

	if got := completed.Load(); got != workers {
		t.Errorf("выполнено %d, ожидалось %d", got, workers)
	}
}

func TestWaitGroup_DoneBlocksUntilWait(t *testing.T) {
	wg := NewSemaWaitGroup()
	wg.Add(1)

	doneReturned := make(chan struct{})
	go func() {
		wg.Done() // должен заблокироваться: в семафор никто не читает
		close(doneReturned)
	}()

	// Убеждаемся, что Done() действительно заблокирован
	select {
	case <-doneReturned:
		t.Fatal("Done() вернулся до вызова Wait() — семантика семафора нарушена")
	case <-time.After(100 * time.Millisecond):
		// ожидаемо
	}

	wg.Wait()

	// После Wait() горутина должна разблокироваться
	select {
	case <-doneReturned:
		// успех
	case <-time.After(time.Second):
		t.Fatal("Done() не разблокировался после Wait()")
	}
}

func TestWaitGroup_WaitBlocksUntilAllDone(t *testing.T) {
	wg := NewSemaWaitGroup()

	const workers = 5
	wg.Add(workers)

	// Запускаем воркеров с разной задержкой
	for i := 0; i < workers; i++ {
		go func(id int) {
			time.Sleep(time.Duration(50*(id+1)) * time.Millisecond)
			wg.Done()
		}(i)
	}

	start := time.Now()
	wg.Wait()
	elapsed := time.Since(start)

	// Самый медленный воркер спит 250 мс. Wait должен ждать его
	if elapsed < 240*time.Millisecond {
		t.Errorf("Wait() вернулся слишком рано: %v", elapsed)
	}
	if elapsed > 500*time.Millisecond {
		t.Errorf("Wait() вернулся слишком поздно: %v", elapsed)
	}
}

func TestWaitGroup_AddMultipleTimes(t *testing.T) {
	wg := NewSemaWaitGroup()

	// Add можно вызывать несколько раз — счётчик аккумулируется
	wg.Add(2)
	wg.Add(3)

	const total = 5
	for i := 0; i < total; i++ {
		go wg.Done()
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Wait() не вернулся после всех Done")
	}
}

func TestWaitGroup_NegativeAddPanics(t *testing.T) {
	wg := NewSemaWaitGroup()

	defer func() {
		if r := recover(); r == nil {
			t.Error("ожидалась паника при отрицательном Add")
		}
	}()

	wg.Add(-1)
}

// Wait() должен корректно дождаться всех воркеров,
// даже если некоторые из них завершаются ДО вызова Wait()
func TestWaitGroup_WorkersFinishBeforeWait(t *testing.T) {
	wg := NewSemaWaitGroup()

	const workers = 10
	wg.Add(workers)

	// Быстрые воркеры: успевают завершиться, пока главная горутина спит
	for i := 0; i < workers; i++ {
		go func() {
			wg.Done() // заблокируется на отправке в семафор
		}()
	}

	// Даём воркерам шанс дойти до Done() и «зависнуть» на семафоре
	time.Sleep(50 * time.Millisecond)

	// Теперь Wait() должен принять все токены и вернуться
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Wait() не разблокировал воркеров")
	}
}
