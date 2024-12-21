package hello

import (
	"context"
	"fmt"
)

type Hello struct{}

func New() (*Hello, error) {
	return &Hello{}, nil
}

func (h Hello) Name() string {
	return "hello"
}

func (h Hello) Description() string {
	return "tools hello"
}

func (h Hello) Call(ctx context.Context, args ...interface{}) (string, error) {
	return fmt.Sprintf("%v\n", args), nil
}
