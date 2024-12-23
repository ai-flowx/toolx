package hello

import (
	"context"
	"fmt"
)

const (
	name        = "hello"
	description = "tools hello"
)

type Hello struct{}

func New() (*Hello, error) {
	return &Hello{}, nil
}

func (h Hello) Name() string {
	return name
}

func (h Hello) Description() string {
	return description
}

func (h Hello) Call(ctx context.Context, args ...interface{}) (string, error) {
	return fmt.Sprintf("%v\n", args), nil
}
