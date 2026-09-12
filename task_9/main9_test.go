package main

import (
	"context"
	"math"
	"testing"
	"time"
)

// collect читает всё из канала до закрытия (с таймаутом)
func collectFloats(ch <-chan float64, timeout time.Duration) ([]float64, bool) {
	var got []float64
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

func TestProducer(t *testing.T) {
	ctx := context.Background()
	ch := Producer(ctx, []uint8{1, 2, 3})

	var got []uint8
	for v := range ch {
		got = append(got, v)
	}

	want := []uint8{1, 2, 3}
	if len(got) != len(want) {
		t.Fatalf("got = %v, want = %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("got[%d] = %d, want %d", i, got[i], want[i])
		}
	}
}

func TestProducer_ContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	values := make([]uint8, 1000)
	ch := Producer(ctx, values)

	<-ch
	<-ch
	cancel()

	// канал должен закрыться, а не зависнуть
	deadline := time.After(time.Second)
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return
			}
		case <-deadline:
			t.Fatal("канал не закрылся после отмены")
		}
	}
}

func TestPipeline(t *testing.T) {
	ctx := context.Background()

	in := make(chan uint8, 5)
	for _, v := range []uint8{1, 2, 3, 4, 5} {
		in <- v
	}
	close(in)

	out := Pipeline(ctx, in)
	got, ok := collectFloats(out, time.Second)
	if !ok {
		t.Fatal("канал не закрылся")
	}

	want := []float64{1, 8, 27, 64, 125}
	if len(got) != len(want) {
		t.Fatalf("got = %v, want = %v", got, want)
	}
	for i := range want {
		if math.Abs(got[i]-want[i]) > 1e-9 {
			t.Errorf("got[%d] = %v, want %v", i, got[i], want[i])
		}
	}
}

func TestPipeline_ContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	// бесконечный источник
	in := make(chan uint8)
	go func() {
		defer close(in)
		i := uint8(0)
		for {
			select {
			case <-ctx.Done():
				return
			case in <- i:
				i++
			}
		}
	}()

	out := Pipeline(ctx, in)
	<-out
	<-out
	cancel()

	deadline := time.After(time.Second)
	for {
		select {
		case _, ok := <-out:
			if !ok {
				return
			}
		case <-deadline:
			t.Fatal("канал не закрылся после отмены")
		}
	}
}

func TestPipelineWithProducer(t *testing.T) {
	ctx := context.Background()

	input := []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	ch1 := Producer(ctx, input)
	ch2 := Pipeline(ctx, ch1)

	got, ok := collectFloats(ch2, 2*time.Second)
	if !ok {
		t.Fatal("канал не закрылся")
	}
	if len(got) != len(input) {
		t.Fatalf("got %d значений, want %d", len(got), len(input))
	}
	for i, v := range input {
		want := float64(v) * float64(v) * float64(v)
		if math.Abs(got[i]-want) > 1e-9 {
			t.Errorf("got[%d] = %v, want %v", i, got[i], want)
		}
	}
}
