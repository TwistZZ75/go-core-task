package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func Merge[T any](ctx context.Context, channels ...<-chan T) <-chan T {
	out := make(chan T)

	if len(channels) == 0 {
		close(out)
		return out
	}

	var wg sync.WaitGroup
	wg.Add(len(channels))

	for _, ch := range channels {
		go func(in <-chan T) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case v, ok := <-in:
					if !ok {
						return
					}
					select {
					case <-ctx.Done():
						return
					case out <- v:
					}
				}
			}
		}(ch)
	}

	// Отдельная горутина ждёт завершения и закрывает out.
	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Три входных канала с разными наборами значений.
	ch1 := producer(ctx, []int{1, 2, 3})
	ch2 := producer(ctx, []int{10, 20, 30})
	ch3 := producer(ctx, []int{100, 200})

	merged := Merge(ctx, ch1, ch2, ch3)

	sum := 0
	count := 0
	for v := range merged {
		fmt.Println(v)
		sum += v
		count++
	}

	fmt.Printf("Получено %d значений, сумма = %d\n", count, sum)
}

// producer — вспомогательная функция: отдаёт значения из слайса в канал
// и закрывает его по завершении
func producer(ctx context.Context, values []int) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for _, v := range values {
			select {
			case <-ctx.Done():
				return
			case ch <- v:
			}
		}
	}()
	return ch
}
