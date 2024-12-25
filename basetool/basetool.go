package basetool

/*
#cgo CFLAGS: -I/usr/include/python3.10
#cgo LDFLAGS: -lpython3.10
#include <Python.h>
*/
import "C"

import (
	"context"
	_ "embed"
	"unsafe"

	"github.com/pkg/errors"
)

const (
	name        = "basetool"
	description = "basetool tools"
)

//go:embed basetool.py
var source string

type BaseTool struct{}

func (b BaseTool) Init(_ context.Context) error {
	C.Py_Initialize()

	if C.Py_IsInitialized() == 0 {
		return errors.New("failed to init python\n")
	}

	return nil
}

func (b BaseTool) Deinit(_ context.Context) error {
	C.Py_Finalize()

	return nil
}

func (b BaseTool) Name(_ context.Context) string {
	return name
}

func (b BaseTool) Description(_ context.Context) string {
	return description
}

func (b BaseTool) Call(_ context.Context, _ ...interface{}) (string, error) {
	cstr := C.CString(source)
	defer C.free(unsafe.Pointer(cstr))

	globals := C.PyDict_New()
	defer C.Py_DecRef(globals)

	locals := C.PyDict_New()
	defer C.Py_DecRef(locals)

	result := C.PyRun_String(cstr, C.Py_file_input, globals, locals)
	if result == nil {
		C.PyErr_Print()
		return "", errors.New("failed to run python\n")
	}

	defer C.Py_DecRef(result)

	resultStr := C.PyObject_Str(result)
	defer C.Py_DecRef(resultStr)

	resultCStr := C.PyUnicode_AsUTF8(resultStr)

	return C.GoString(resultCStr), nil
}
