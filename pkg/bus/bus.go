package bus

import "context"

type HandlerFunc any

type Msg any

type Bus interface {
	Publish(ctx context.Context, msg Msg) error
	AddEventListener(handler HandlerFunc)
}

type AsyncEventBus struct {
	handlers map[string][]HandlerFunc
}

func (a AsyncEventBus) Publish(ctx context.Context, msg Msg) error {
	//TODO implement me
	panic("implement me")
}

func (a AsyncEventBus) AddEventListener(handler HandlerFunc) {
	//TODO implement me
	panic("implement me")
}

func NewAsyncEventBus() Bus {
	return &AsyncEventBus{
		handlers: make(map[string][]HandlerFunc),
	}
}
