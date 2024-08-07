package bus

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"short-link/utils/gox"
	"sync"
)

type HandlerFunc any

type Msg any

type Bus interface {
	Publish(ctx context.Context, topic string, msg Msg) error
	AddEventListener(topic string, handler HandlerFunc)
}

type AsyncEventBus struct {
	handlers map[string][]HandlerFunc
	lock     sync.RWMutex
}

var (
	asyncEventBus *AsyncEventBus
	once          sync.Once
)

func GetEventBus() *AsyncEventBus {
	return asyncEventBus
}

func NewAsyncEventBus() *AsyncEventBus {
	once.Do(func() {
		asyncEventBus = &AsyncEventBus{
			handlers: make(map[string][]HandlerFunc),
		}
	})
	return asyncEventBus
}

func (bus *AsyncEventBus) Publish(ctx context.Context, topic string, msg Msg) error {
	bus.lock.RLock()
	handlers, exists := bus.handlers[topic]
	bus.lock.RUnlock()
	if !exists {
		return nil
	}
	params := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(msg)}
	return callListeners(handlers, params)
}

func (bus *AsyncEventBus) AddEventListener(topic string, handler HandlerFunc) {
	bus.lock.Lock()
	defer bus.lock.Unlock()
	v := reflect.ValueOf(handler)
	if v.Type().Kind() != reflect.Func {
		panic("handler must be of type HandlerFunc")
	}
	bus.handlers[topic] = append(bus.handlers[topic], handler)
}

func callListeners(handlers []HandlerFunc, params []reflect.Value) error {
	var errs error
	wg := gox.NewWaitGroup()
	errChan := make(chan error, len(handlers))
	for _, handler := range handlers {
		fn := handler
		wg.RunSafe(context.Background(), func(ctx context.Context) {
			ret := reflect.ValueOf(fn).Call(params)
			if len(ret) > 0 && !ret[0].IsNil() {
				e := ret[0].Interface()
				if e != nil {
					err, ok := e.(error)
					if ok {
						errChan <- err
					}
					errChan <- fmt.Errorf("expected listener to return an error, got '%T'", e)
				}
			}
		})
	}

	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			errs = errors.Join(errs, err)
		}
		return nil
	}
	return errs
}
