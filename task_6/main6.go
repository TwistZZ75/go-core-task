package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ch := rndFunc(ctx)

	for i := 0; i < 10; i++ {
		v, ok := <-ch
		if !ok {
			break
		}
		fmt.Println(v)
	}

	fmt.Println("--- завершение по таймауту ---")
	for v := range ch {
		fmt.Println(v)
	}
	fmt.Println("канал закрыт")
}

func rndFunc(ctx context.Context) <-chan int {
	ch := make(chan int)

	go func() {
		defer close(ch)
		for {
			value, err := rand.Int(rand.Reader, big.NewInt(100))
			if err != nil {
				fmt.Println("rand error: ", err)
				continue
			}

			select {
			case <-ctx.Done():
				return
			case ch <- int(value.Int64()):
				// значение отправлено получателю, продолжаем
			}
		}
	}()

	return ch
}
