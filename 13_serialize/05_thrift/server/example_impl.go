package main

import (
	"context"
	"fmt"
)

type ExampleServiceImpl struct{}

func (e *ExampleServiceImpl) SayHello(ctx context.Context, name string) (err error) {
	fmt.Printf("Hello, %s!\n", name)
	return nil
}
