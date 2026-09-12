package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

type SemaWaitGroup struct {
	counter atomic.Int64  // сколько Done() ждёт Wait()
	sem     chan struct{} // сам семафор — небуферизированный канал
}

// NewSemaWaitGroup создаёт новый WaitGroup.
func NewSemaWaitGroup() *SemaWaitGroup {
	return &SemaWaitGroup{
		sem: make(chan struct{}),
	}
}

// Add увеличивает счётчик ожидаемых завершений на delta
func (wg *SemaWaitGroup) Add(delta int) {
	if delta < 0 {
		panic("WaitGroup: negative Add")
	}
	wg.counter.Add(int64(delta))
}

// Done сигнализирует о завершении одной задачи
// отправляет токен в семафор
// блокируется до тех пор, пока Wait() не примет этот токен
func (wg *SemaWaitGroup) Done() {
	wg.sem <- struct{}{}
}

// Wait блокируется, пока не получит из семафора ровно столько токенов,
// сколько было заявлено в Add
func (wg *SemaWaitGroup) Wait() {
	n := wg.counter.Load()
	for i := int64(0); i < n; i++ {
		<-wg.sem
	}
}

func main() {
	wg := NewSemaWaitGroup()

	const workers = 5
	wg.Add(workers)

	for i := 1; i <= workers; i++ {
		go func(id int) {
			defer wg.Done()
			time.Sleep(time.Duration(id) * 100 * time.Millisecond)
			fmt.Printf("worker %d finished\n", id)
		}(i)
	}

	start := time.Now()
	wg.Wait()
	fmt.Printf("all workers done in %v\n", time.Since(start))
}
