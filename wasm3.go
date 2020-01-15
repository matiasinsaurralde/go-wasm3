package wasm3

/*
#cgo CFLAGS: -Iinclude
#cgo darwin LDFLAGS: -L${SRCDIR}/lib/darwin -lm3
#cgo linux LDFLAGS: -L${SRCDIR}/lib/linux -lm3 -lm
#include "m3.h"
#include "m3_api_libc.h"
*/
import "C"

import(
	"unsafe"
	"fmt"
	"errors"
)

// RuntimeT is an alias for IM3Runtime
type RuntimeT C.IM3Runtime
// EnvironmentT is an alias for IM3Environment
type EnvironmentT C.IM3Environment
// ModuleT is an alias for IM3Module
type ModuleT C.IM3Module
// FunctionT is an alias for IM3Function
type FunctionT C.IM3Function
// ResultT is an alias for M3Result
type ResultT C.M3Result

var(
	errParseModule = errors.New("Parse error")
	errLoadModule = errors.New("Load error")
	errFuncLookupFailed = errors.New("Function lookup failed")
)

// Runtime wraps a WASM3 runtime
type Runtime struct {
	ptr RuntimeT
	Environment *Environment
}

// Ptr returns a IM3Runtime pointer
func(r *Runtime) Ptr() C.IM3Runtime {
	return (C.IM3Runtime)(r.ptr)
}
// Load wraps the parse and load module calls.
func(r *Runtime) Load(wasmBytes []byte) (ModuleT, error) {
	result := C.m3Err_none
	bytes := C.CBytes(wasmBytes)
	length := len(wasmBytes)
	var module C.IM3Module
	fmt.Printf("module=%p\n",module)
	result = C.m3_ParseModule(
		r.Environment.Ptr(),
		&module,
		(*C.uchar)(bytes),
		C.uint(length),
	)
	if result != nil {
		return nil, errParseModule
	}
	result = C.m3_LoadModule(
		r.Ptr(),
		module,
	)
	if result != nil {
		return nil, errLoadModule
	}
	// result = C.m3_LinkSpecTest((C.IM3Runtime)(r.Ptr()).modules)
	return (ModuleT)(module), nil
}

// FindFunction calls m3_FindFunction and returns a call function
func(r *Runtime) FindFunction(funcName string) (Function, error) {
	result := C.m3Err_none
	var f C.IM3Function
	cFuncName := C.CString(funcName)
	defer C.free(unsafe.Pointer(cFuncName))
	result = C.m3_FindFunction(
		&f,
		r.Ptr(),
		cFuncName,
	)
	if result != nil {
		return nil, errFuncLookupFailed
	}
	fnWrapper := func(args... string) {	
		length := len(args)
		cArgs := make([]*C.char, length)
		for i, v := range args {
			cVal := C.CString(v)
			cArgs[i] = cVal
		}
		C.m3_CallWithArgs(f, C.uint(length), &cArgs[0])
	}
	return fnWrapper, nil
}

// NewRuntime initializes a new runtime
// TODO: nativeStackInfo is passed as NULL
func NewRuntime(env *Environment, stackSize uint) *Runtime {
	ptr := C.m3_NewRuntime(
		env.Ptr(),
		C.uint(stackSize),
		nil,
	)
	return &Runtime{
		ptr: (RuntimeT)(ptr),
		Environment: env,
	}
}

// Function is a function wrapper
type Function func(args ...string)

// Environment wraps a WASM3 environment
type Environment struct {
	ptr EnvironmentT
}

// Ptr returns a pointer to IM3Environment
func(e *Environment) Ptr() C.IM3Environment {
	return (C.IM3Environment)(e.ptr)
}

// NewEnvironment initializes a new environment
func NewEnvironment() *Environment {
	ptr := C.m3_NewEnvironment()
	return &Environment{
		ptr: (EnvironmentT)(ptr),
	}
}
