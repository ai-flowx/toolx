package structuredtool

/*
#cgo CFLAGS: -I/usr/include/python3.13
#cgo LDFLAGS: -lpython3.13

#define Py_LIMITED_API

#include <Python.h>

extern void Py_Initialize(void);
extern void Py_Finalize(void);
extern int PyRun_SimpleString(const char*);

static PyObject* PyInit_gomodule(void) {
	static PyMethodDef methods[] = {
		{NULL, NULL, 0, NULL}
	};

	static PyModuleDef module = {
		PyModuleDef_HEAD_INIT,
		"gomodule",
		"Go module for Python",
		-1,
		methods
	};

	return PyModule_Create(&module);
}

static void init_python() {
	Py_Initialize();
}

static void call_python_function(const char* script) {
	PyRun_SimpleString(script);
}

static void finalize_python() {
	Py_Finalize();
}
*/
import "C"

import (
	"context"
	"unsafe"
)

const (
	name        = "structuredtool"
	description = "structuredtool tools"
)

type StructuredTool struct{}

func (s StructuredTool) Init(_ context.Context) error {
	C.init_python()
	return nil
}

func (s StructuredTool) Deinit(_ context.Context) error {
	C.finalize_python()
	return nil
}

func (s StructuredTool) Name(_ context.Context) string {
	return name
}

func (s StructuredTool) Description(_ context.Context) string {
	return description
}

func (s StructuredTool) Call(ctx context.Context, args ...interface{}) (string, error) {
	script := `print("hello structuredtool")`

	buf := C.CString(script)
	defer C.free(unsafe.Pointer(buf))

	C.call_python_function(buf)

	return "", nil
}