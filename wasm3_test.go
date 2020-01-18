package wasm3

import (
	"io/ioutil"
	"testing"
)

const (
	sumModulePath = "examples/sum/sum.wasm"
)

var (
	sumModuleBytes []byte
)

func init() {
	var err error
	sumModuleBytes, err = ioutil.ReadFile(sumModulePath)
	if err != nil {
		panic(err)
	}
}
func TestEnvRuntimeCycle(t *testing.T) {
	runtime := NewRuntime(&Config{
		Environment: NewEnvironment(),
		StackSize:   64 * 1024,
	})
	defer runtime.Destroy()
}

func TestParseModule(t *testing.T) {
	env := NewEnvironment()
	_, err := env.ParseModule([]byte(""))
	if err == nil {
		t.Fatal("Invalid input should error")
	}
	module, err := env.ParseModule(sumModuleBytes)
	if err != nil {
		t.Fatal("Couldn't parse valid WASM module")
	}
	if module.ptr == nil {
		t.Fatal("Internal module pointer is nil")
	}
}

func TestLoadModule(t *testing.T) {
	runtime := NewRuntime(&Config{
		Environment: NewEnvironment(),
		StackSize:   64 * 1024,
	})
	defer runtime.Destroy()
	module, _ := runtime.ParseModule(sumModuleBytes)
	_, err := runtime.LoadModule(module)
	if err != nil {
		t.Fatal("Couldn't load sample module")
	}
	_, err = runtime.FindFunction("nonexistent")
	if err == nil {
		t.Fatal("No error when referencing a nonexistent function")
	}
	_, err = runtime.FindFunction("sum")
	if err != nil {
		t.Fatal("Couldn't find 'sum' test function")
	}
}

func TestModuleHelpers(t *testing.T) {
	env := NewEnvironment()
	module, _ := env.ParseModule(sumModuleBytes)
	if module.numFunctions != -1 {
		t.Fatal("Initial value for numFunctions should be -1")
	}
	fn, err := module.GetFunctionByName("sum")
	if err != nil {
		t.Fatal("Couldn't find 'sum' test function using name lookup")
	}
	if fn.Name == "" {
		t.Fatal("Function name is empty")
	}
	fn2, err := module.GetFunction(0)
	if err != nil {
		t.Fatal("Couldn't find 'sum' test function using index lookup")
	}
	if fn2.Name == "" {
		t.Fatal("Function name is empty")
	}
	if module.NumFunctions() != 1 {
		t.Fatal("Module NumFunctions should be 1")
	}
}
