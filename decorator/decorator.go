package decorator

/*
#cgo CFLAGS: -I/usr/include/python3.13
#cgo LDFLAGS: -lpython3.13
#define Py_LIMITED_API
#include <Python.h>

extern void Py_Initialize(void);
extern void Py_Finalize(void);
extern PyObject *PyRun_String(const char *str, int start, PyObject *globals, PyObject *locals);
extern const char *PyUnicode_AsUTF8(PyObject *unicode);

static void init_python() {
	Py_Initialize();
}

static void finalize_python() {
	Py_Finalize();
}

static PyObject* run_python(const char* source) {
	PyObject* result = PyRun_String(source, Py_eval_input, PyEval_GetGlobals(), PyEval_GetLocals());
	return result;
}

static char* get_result(PyObject* obj) {
	if (!obj) {
		return NULL;
	}

	PyObject* ret = PyObject_Str(obj);
	if (!ret) {
		return NULL;
	}

	return (char*)PyUnicode_AsUTF8(ret);
}
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

func (d Decorator) Init(_ context.Context) error {
	C.init_python()
	return nil
}

func (d Decorator) Deinit(_ context.Context) error {
	C.finalize_python()
	return nil
}

func (d Decorator) Name(_ context.Context) string {
	return name
}

func (d Decorator) Description(_ context.Context) string {
	return description
}

func (d Decorator) Call(ctx context.Context, args ...interface{}) (string, error) {
	cstr := C.CString(source)
	defer C.free(unsafe.Pointer(cstr))

	buf := C.run_python(cstr)
	if buf == nil {
		return "", errors.New("failed to run python\n")
	}

	defer C.Py_DecRef(buf)

	ret := C.get_result(buf)
	if ret == nil {
		return "", errors.New("failed to get result\n")
	}

	return C.GoString(ret), nil
}
