package langchain

/*
#cgo CFLAGS: -I/usr/include/python3.12
#cgo LDFLAGS: -lpython3.12

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
	name        = "langchain"
	description = "langchain tools"
)

type LangChain struct{}

func New() (*LangChain, error) {
	return &LangChain{}, nil
}

func (l LangChain) Init() error {
	C.init_python()
	return nil
}

func (l LangChain) Deinit() error {
	C.finalize_python()
	return nil
}

func (l LangChain) Name() string {
	return name
}

func (l LangChain) Description() string {
	return description
}

func (l LangChain) Call(ctx context.Context, args ...interface{}) (string, error) {
	script := `print("hello langchain")`

	s := C.CString(script)
	defer C.free(unsafe.Pointer(s))

	C.call_python_function(s)

	return "", nil
}
