package decorator

/*
#cgo pkg-config: python
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
	name        = "decorator"
	description = "decorator tools"
)

//go:embed decorator.py
var source string

type Decorator struct{}

func (d *Decorator) Init(_ context.Context) error {
	C.Py_Initialize()

	if C.Py_IsInitialized() == 0 {
		return errors.New("failed to init python\n")
	}

	return nil
}

func (d *Decorator) Deinit(_ context.Context) error {
	C.Py_Finalize()

	return nil
}

func (d *Decorator) Name(_ context.Context) string {
	return name
}

func (d *Decorator) Description(_ context.Context) string {
	return description
}

func (d *Decorator) Call(_ context.Context, _ func(context.Context, interface{}) (interface{}, error), _ ...interface{}) (string, error) {
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
