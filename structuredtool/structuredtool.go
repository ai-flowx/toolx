package structuredtool

/*
#cgo CFLAGS: -I/usr/include/python3.10
#cgo LDFLAGS: -lpython3.10
#include <Python.h>
*/
import "C"

import (
	"context"
	_ "embed"
)

const (
	name        = "structuredtool"
	description = "structuredtool tools"
)

//go:embed structuredtool.py
var source string

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
