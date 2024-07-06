package gox

import (
	"context"
	"errors"
	"fmt"
	"log"
	"short-link/logs"
	"sync"
)

type WaitGroup struct {
	wg *sync.WaitGroup
}

func NewWaitGroup() *WaitGroup {
	return &WaitGroup{
		wg: &sync.WaitGroup{},
	}
}

func (w *WaitGroup) RunSafe(ctx context.Context, fn func(ctx context.Context)) {
	w.wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic: %v\n", r)
			}
		}()
		fn(ctx)
		w.wg.Done()

	}()
}

func (w *WaitGroup) Wait() {
	w.wg.Wait()
}

func Run(ctx context.Context, fn func(ctx context.Context)) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Error(fmt.Errorf("Recovered from panic: %v\n", r), "gox run panic")
			}
		}()
		fn(ctx)
	}()
}

type ErrorWaitGroup struct {
	wg   sync.WaitGroup
	lock sync.Mutex
	err  error
}

func NewErrorWaitGroup() *ErrorWaitGroup {
	return &ErrorWaitGroup{}
}

func (e *ErrorWaitGroup) RunSafe(ctx context.Context, fn func(ctx context.Context) error) {
	e.wg.Add(1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logs.Error(fmt.Errorf("Recovered from panic: %v\n", r), "run safe panic")
			}
			e.wg.Done()
		}()
		if err := fn(ctx); err != nil {
			e.lock.Lock()
			defer e.lock.Unlock()
			e.err = errors.Join(e.err, err)
		}
	}()
}

func (e *ErrorWaitGroup) Wait() error {
	e.wg.Wait()
	return e.err
}
