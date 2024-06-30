package gox

import (
	"context"
	"fmt"
	"sync"
)

// TODO 优化
func RunSafe(ctx context.Context, wg *sync.WaitGroup, fn func(ctx context.Context)) {
	if wg != nil {
		wg.Add(1)
	}
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

func Run(ctx context.Context, fn func(ctx context.Context)) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovered from panic: %v\n", r)
			}
		}()
		fn(ctx)
	}()
}
