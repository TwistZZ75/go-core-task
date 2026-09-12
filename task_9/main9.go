package main

import (
	"context"
	"fmt"
)

func Pipeline(ctx context.Context, in <-chan uint8) <-chan float64 {
	out := make(chan float64)

	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				cube := float64(v) * float64(v) * float64(v)

				select {
				case <-ctx.Done():
					return
				case out <- cube:
				}
			}
		}
	}()

	return out
}

// Producer пишет числа из values в канал
// и закрывает его по завершении или при отмене ctx
func Producer(ctx context.Context, values []uint8) <-chan uint8 {
	ch := make(chan uint8)
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

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	input := []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Println("Входные числа:", input)

	ch1 := Producer(ctx, input)
	ch2 := Pipeline(ctx, ch1)

	fmt.Println("Результаты (v³ как float64):")
	for v := range ch2 {
		fmt.Printf("%v\n", v)
	}
	fmt.Println("Конвейер завершён")
}
