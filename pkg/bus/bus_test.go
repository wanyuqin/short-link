package bus

import (
	"context"
	"fmt"
	"testing"
)

type User struct {
	UserName string
	Age      int
}

func sub1(ctx context.Context, user User) error {
	fmt.Println("sub1", user.UserName, user.Age)
	return nil
}

func sub2(ctx context.Context, user User) error {
	fmt.Println("sub2", user.UserName, user.Age)
	return nil
}

func TestAsyncEventBus_Publish(t *testing.T) {
	bus := NewAsyncEventBus()
	bus.AddEventListener("topic-1", sub1)
	bus.AddEventListener("topic-1", sub2)

	err := bus.Publish(context.Background(), "topic-1", User{
		UserName: "foo",
		Age:      1,
	})
	if err != nil {
		t.Error(err)
	}
}
