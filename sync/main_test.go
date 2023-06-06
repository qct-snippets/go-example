package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestWaitGroup(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		i := i
		go func() {
			defer wg.Done()
			println(i)
		}()
	}
	wg.Wait()
}

func TestOnce(t *testing.T) {
	var once sync.Once
	once.Do(func() {
		println("1")
	})
	once.Do(func() {
		println("2")
	})
}

func TestErrorGroup(t *testing.T) {
	errGroup, ctx := errgroup.WithContext(context.Background())
	for i := 0; i < 6; i++ {
		i := i
		errGroup.Go(func() error {
			select {
			case <-ctx.Done():
				fmt.Printf("Canceled:%d\n", i)
				return ctx.Err()
			default:
				if i > 3 {
					return fmt.Errorf("error: %d", i)
				}
				fmt.Printf("Success: %d\n", i)
				return nil
			}
		})
	}

	if err := errGroup.Wait(); err != nil {
		log.Fatal(err)
	}
}
