package structuredtool

import (
	"context"
)

const (
	name        = "structuredtool"
	description = "structuredtool tools"
)

type StructuredTool struct{}

func (s StructuredTool) Init(_ context.Context) error {
	return nil
}

func (s StructuredTool) Deinit(_ context.Context) error {
	return nil
}

func (s StructuredTool) Name(_ context.Context) string {
	return name
}

func (s StructuredTool) Description(_ context.Context) string {
	return description
}

func (s StructuredTool) Call(ctx context.Context, args ...interface{}) (string, error) {
	return "", nil
}
