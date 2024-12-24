package basetool

import (
	"context"
)

const (
	name        = "basetool"
	description = "basetool tools"
)

type BaseTool struct{}

func (b BaseTool) Init(_ context.Context) error {
	return nil
}

func (b BaseTool) Deinit(_ context.Context) error {
	return nil
}

func (b BaseTool) Name(_ context.Context) string {
	return name
}

func (b BaseTool) Description(_ context.Context) string {
	return description
}

func (b BaseTool) Call(ctx context.Context, args ...interface{}) (string, error) {
	return "", nil
}
