package wasm3

/*
#cgo CFLAGS: -Iinclude
#cgo darwin LDFLAGS: -L${SRCDIR}/lib/darwin -lm3
#cgo linux LDFLAGS: -L${SRCDIR}/lib/linux -lm3 -lm
#include "m3.h"
#include "m3_api_libc.h"
#include "m3_env.h"

// module_get_function is a helper function for the module Go struct
IM3Function module_get_function(IM3Module i_module, int index) {
	IM3Function f = & i_module->functions [index];
	return f;
}
*/
import "C"

import(
	"unsafe"
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
// This will be replaced by env.ParseModule and Runtime.LoadModule.
func(r *Runtime) Load(wasmBytes []byte) (*Module, error) {
	result := C.m3Err_none
	bytes := C.CBytes(wasmBytes)
	length := len(wasmBytes)
	var module C.IM3Module
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
	result = C.m3_LinkSpecTest(r.Ptr().modules)
	if result != nil {
		return nil, errors.New("LinkSpecTest failed")
	}
	m := NewModule((ModuleT)(module))
	return m, nil
}

// LoadModule wraps m3_LoadModule and returns a module object
func(r *Runtime) LoadModule(module *Module) (*Module, error) {
	result := C.m3Err_none
	result = C.m3_LoadModule(
		r.Ptr(),
		module.Ptr(),
	)
	if result != nil {
		return nil, errLoadModule
	}
	result = C.m3_LinkSpecTest(r.Ptr().modules)
	if result != nil {
		return nil, errors.New("LinkSpecTest failed")
	}
	return module, nil
}

// FindFunction calls m3_FindFunction and returns a call function
func(r *Runtime) FindFunction(funcName string) (FunctionWrapper, error) {
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
	var fnWrapper FunctionWrapper
	fnWrapper = func(args... string) {	
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

// Destroy free calls m3_FreeRuntime
func(r *Runtime) Destroy() {
    C.m3_FreeRuntime(r.Ptr());
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

// Module wraps a WASM3 module.
type Module struct {
	ptr ModuleT
	numFunctions int
}

// Ptr returns a pointer to IM3Module
func(m *Module) Ptr() C.IM3Module {
	return (C.IM3Module)(m.ptr)
}

// GetFunction provides access to IM3Function->functions
func(m *Module) GetFunction(index uint) (*Function, error) {
	if uint(m.NumFunctions()) <= index {
		return nil, errFuncLookupFailed
	}
	ptr := C.module_get_function(m.Ptr(), C.int(index))
	name := C.GoString(ptr.name)
	return &Function{
		ptr: (FunctionT)(ptr),
		Name: name,
	}, nil
}

// GetFunctionByName is a helper to lookup functions by name
// TODO: could be optimized by caching function names and pointer on the Go side, right after the load call.
func(m *Module) GetFunctionByName(lookupName string) (*Function, error) {
	var fn *Function
	for i :=0 ; i < m.NumFunctions(); i++ {
		ptr := C.module_get_function(m.Ptr(), C.int(i))
		name := C.GoString(ptr.name)
		if name != lookupName {
			continue	
		}
		fn = &Function{
			ptr: (FunctionT)(ptr),
			Name: name,
		}
		return fn, nil
	}
	return nil, errFuncLookupFailed
}

// NumFunctions provides access to numFunctions.
func(m *Module) NumFunctions() int {
	// In case the number of functions hasn't been resolved yet, retrieve the int and keep it in the structure
	if m.numFunctions == -1 {
		m.numFunctions = int(m.Ptr().numFunctions)
	}
	return m.numFunctions
}

// NewModule wraps a WASM3 moduke
func NewModule(ptr ModuleT) *Module {
	return &Module{
		ptr: ptr,
		numFunctions: -1,
	}
}

// Function is a function wrapper
type Function struct {
	ptr FunctionT
	// fnWrapper FunctionWrapper
	Name string
}

// FunctionWrapper is used to wrap WASM3 call methods and make the calls more idiomatic
// TODO: this is very limited, we need to handle input and output types appropriately
type FunctionWrapper func(args ...string)

// Ptr returns a pointer to IM3Function
func(f *Function) Ptr() C.IM3Function {
	return (C.IM3Function)(f.ptr)
}

// Call wraps m3_CallWithArgs
func(f *Function) Call(args... string) {
	length := len(args)
	cArgs := make([]*C.char, length)
	for i, v := range args {
		cVal := C.CString(v)
		cArgs[i] = cVal
	}
	C.m3_CallWithArgs(f.Ptr(), C.uint(length), &cArgs[0])
}

// Environment wraps a WASM3 environment
type Environment struct {
	ptr EnvironmentT
}

// ParseModule wraps m3_ParseModule
func(e *Environment) ParseModule(wasmBytes []byte) (*Module, error) {
	result := C.m3Err_none
	bytes := C.CBytes(wasmBytes)
	length := len(wasmBytes)
	var module C.IM3Module
	result = C.m3_ParseModule(
		e.Ptr(),
		&module,
		(*C.uchar)(bytes),
		C.uint(length),
	)
	if result != nil {
		return nil, errParseModule
	}
	return NewModule((ModuleT)(module)), nil
}
// Ptr returns a pointer to IM3Environment
func(e *Environment) Ptr() C.IM3Environment {
	return (C.IM3Environment)(e.ptr)
}

// Destroy calls m3_FreeEnvironment
func(e *Environment) Destroy() {
	C.m3_FreeEnvironment(e.Ptr())
}

// NewEnvironment initializes a new environment
func NewEnvironment() *Environment {
	ptr := C.m3_NewEnvironment()
	return &Environment{
		ptr: (EnvironmentT)(ptr),
	}
}
