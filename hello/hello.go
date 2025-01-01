package hello

import (
	"context"
	"fmt"
)

const (
	name        = "hello"
	description = "hello tools"
)

type Hello struct{}

func (h *Hello) Init(_ context.Context) error {
	return nil
}

func (h *Hello) Deinit(_ context.Context) error {
	return nil
}

func (h *Hello) Name(_ context.Context) string {
	return name
}

func (h *Hello) Description(_ context.Context) string {
	return description
}

func (h *Hello) Call(_ context.Context, _ func(context.Context, interface{}) (interface{}, error), args ...interface{}) (string, error) {
	return fmt.Sprintf("%v\n", args), nil
}
