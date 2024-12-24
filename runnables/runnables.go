package runnables

import (
	"context"
)

const (
	name        = "runnables"
	description = "runnables tools"
)

type Runnables struct{}

func (r Runnables) Init(_ context.Context) error {
	return nil
}

func (r Runnables) Deinit(_ context.Context) error {
	return nil
}

func (r Runnables) Name(_ context.Context) string {
	return name
}

func (r Runnables) Description(_ context.Context) string {
	return description
}

func (r Runnables) Call(ctx context.Context, args ...interface{}) (string, error) {
	return "", nil
}
