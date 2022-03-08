package wasm3

/*
#cgo                                 CFLAGS:  -I${SRCDIR}/include
#cgo darwin                          LDFLAGS: -framework Security -lm -lm3
#cgo darwin,amd64,!ios,!iossimulator LDFLAGS: -L${SRCDIR}/lib/macosx-x86_64
#cgo darwin,arm64,ios,!iossimulator  LDFLAGS: -L${SRCDIR}/lib/iphoneos-arm64
#cgo darwin,arm64,ios,iossimulator   LDFLAGS: -L${SRCDIR}/lib/iphonesimulator-arm64
#cgo darwin,amd64,ios,iossimulator   LDFLAGS: -L${SRCDIR}/lib/iphonesimulator-x86_64
#cgo linux                           LDFLAGS: -lm3 -lm
#cgo linux,arm64,android             LDFLAGS: -L${SRCDIR}/lib/android-aarch64
#cgo linux,amd64,android             LDFLAGS: -L${SRCDIR}/lib/android-x86_64
#cgo linux,amd64,!android            LDFLAGS: -L${SRCDIR}/lib/linux-x86_64

#include "m3_api_libc.h"
#include "m3_api_wasi.h"
#include "m3_env.h"
#include "go-wasm3.h"

// module_get_function is a helper function for the module Go struct
IM3Function module_get_function(IM3Module i_module, int index) {
	IM3Function f = & i_module->functions [index];
	return f;
}

int call(IM3Function i_function, uint32_t i_argc, int i_argv[]) {
	int result = 0;
	IM3Module module = i_function->module;
	IM3Runtime runtime = module->runtime;
	u64* stack = (u64*)(runtime->stack);
	IM3FuncType ftype = i_function->funcType;
	for (int i = 0; i < ftype->numArgs; i++) {
		int v = i_argv[i];
		u64* s = &stack[i];
		*(u64*)(s) = v;
	}
	m3StackCheckInit();
	M3Result call_result = Call(i_function->compiled, (m3stack_t)(stack), runtime->memory.mallocated, d_m3OpDefaultArgs);
	if(call_result != NULL) {
		set_error(call_result);
		return -1;
	}
	switch (ftype->returnType) {
		case c_m3Type_i32:
			result = *(u32*)(stack);
			break;
		case c_m3Type_i64:
		default:
			result =  *(u32*)(stack);
	};
	return result;
}

int get_allocated_memory_length(IM3Runtime i_runtime) {
	return i_runtime->memory.mallocated->length;
}

u8* get_allocated_memory(IM3Runtime i_runtime) {
	return m3MemData(i_runtime->memory.mallocated);
}
*/
import "C"

import(
	"unsafe"
	"errors"
	"reflect"
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

// Config holds the runtime and environment configuration
type Config struct {
	Environment *Environment
	StackSize uint
	EnableWASI bool
}

// Runtime wraps a WASM3 runtime
type Runtime struct {
	ptr RuntimeT
	cfg *Config
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
		r.cfg.Environment.Ptr(),
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
	if r.cfg.EnableWASI {
		C.m3_LinkWASI(r.Ptr().modules)
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
	if r.cfg.EnableWASI {
		C.m3_LinkWASI(r.Ptr().modules)
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
	fn := &Function{
		ptr: (FunctionT)(f),
	}
	// var fnWrapper FunctionWrapper
	// fnWrapper = fn.Call
	return FunctionWrapper(fn.Call), nil
}

// Destroy free calls m3_FreeRuntime
func(r *Runtime) Destroy() {
	C.m3_FreeRuntime(r.Ptr());
	r.cfg.Environment.Destroy()
}

// Memory allows access to runtime Memory.
// Taken from Wasmer extension: https://github.com/wasmerio/go-ext-wasm
func(r *Runtime) Memory() []byte {
	mem := C.get_allocated_memory(
		r.Ptr(),
	)
	var data = (*uint8)(mem)
	length := r.GetAllocatedMemoryLength()
	var header reflect.SliceHeader
	header = *(*reflect.SliceHeader)(unsafe.Pointer(&header))
	header.Data = uintptr(unsafe.Pointer(data))
	header.Len = int(length)
	header.Cap = int(length)
	return *(*[]byte)(unsafe.Pointer(&header))
}

// GetAllocatedMemoryLength returns the amount of allocated runtime memory
func(r *Runtime) GetAllocatedMemoryLength() int {
	length := C.get_allocated_memory_length(r.Ptr())
	return int(length)
}

// ParseModule is a helper that calls the same function in env.
func(r *Runtime) ParseModule(wasmBytes []byte) (*Module, error) {
	return r.cfg.Environment.ParseModule(wasmBytes)
}

// NewRuntime initializes a new runtime
// TODO: nativeStackInfo is passed as NULL
func NewRuntime(cfg *Config) *Runtime {
	// env *Environment, stackSize uint
	ptr := C.m3_NewRuntime(
		cfg.Environment.Ptr(),
		C.uint(cfg.StackSize),
		nil,
	)
	return &Runtime{
		ptr: (RuntimeT)(ptr),
		cfg: cfg,
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
type FunctionWrapper func(args ...interface{}) (int, error)

// Ptr returns a pointer to IM3Function
func(f *Function) Ptr() C.IM3Function {
	return (C.IM3Function)(f.ptr)
}

// CallWithArgs wraps m3_CallWithArgs
func(f *Function) CallWithArgs(args... string) {
	length := len(args)
	cArgs := make([]*C.char, length)
	for i, v := range args {
		cVal := C.CString(v)
		cArgs[i] = cVal
	}
	C.m3_CallWithArgs(f.Ptr(), C.uint(length), &cArgs[0])
}

// Call implements a better call function
// TODO: support diferent types
func(f *Function) Call(args... interface{}) (int, error) {
	length := len(args)
	if length == 0 {
		result := C.call(f.Ptr(), 0, nil)
		if result == -1 {
			return int(result), errors.New(LastErrorString())
		}
		return int(result), nil
	}
	cArgs := make([]C.int, length)
	for i, v := range args {
		val := v.(int)
		n := C.int(val)
		cArgs[i] = n
	}
	result := C.call(f.Ptr(), C.uint(length), &cArgs[0])
	if result == -1 {
		return int(result), errors.New(LastErrorString())
	}
	return int(result), nil
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
