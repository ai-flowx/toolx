package runnables

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
	name        = "runnables"
	description = "runnables tools"
)

//go:embed runnables.py
var source string

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
