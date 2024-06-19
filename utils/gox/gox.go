package gox

import (
	"context"
	"fmt"
	"sync"
)

func RunSafe(ctx context.Context, wg *sync.WaitGroup, fn func(ctx context.Context)) {
	wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovered from panic: %v\n", r)
			}
			wg.Done()
		}()
		fn(ctx)
	}()
}
