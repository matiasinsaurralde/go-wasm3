package main

import (
	"bytes"
	"io/ioutil"
	golog "log"

	wasm3 "github.com/matiasinsaurralde/go-wasm3"
)

var (
	allocateFn wasm3.FunctionWrapper
	execFn     wasm3.FunctionWrapper
	runtime    *wasm3.Runtime
	print      = golog.Print
	printf     = golog.Printf
)

const (
	wasmFilename = "boa.wasm"
)

func initRuntimeAndModule() error {
	runtime = wasm3.NewRuntime(&wasm3.Config{
		Environment: wasm3.NewEnvironment(),
		StackSize:   1024 * 1024,
	})

	wasmBytes, err := ioutil.ReadFile(wasmFilename)
	if err != nil {
		return err
	}

	module, err := runtime.ParseModule(wasmBytes)
	if err != nil {
		return err
	}
	_, err = runtime.LoadModule(module)
	if err != nil {
		return err
	}
	return nil
}

func mapCalls() error {
	var err error
	allocateFn, err = runtime.FindFunction("boa_alloc")
	if err != nil {
		return err
	}

	execFn, err = runtime.FindFunction("boa_exec3")
	if err != nil {
		return err
	}
	return nil
}

func allocate(input string) (int, error) {
	ptr, err := allocateFn(len(input))
	if err != nil {
		return 0, nil
	}
	pos := ptr
	for _, ch := range input {
		runtime.Memory()[pos] = byte(ch)
		pos++
	}
	return ptr, nil
}

func exec(ptr, length int) (string, error) {
	outPtr, err := execFn(ptr, length)
	if err != nil {
		return "", err
	}
	printf("\"boa_exec3\" returned, output pointer is %d\n", outPtr)
	buf := new(bytes.Buffer)
	for {
		ch := runtime.Memory()[outPtr]
		if ch == 0 {
			break
		}
		buf.WriteByte(ch)
		outPtr++
	}
	printf("Read %d bytes from WASM memory, starting in %d\n", buf.Len(), outPtr)
	return buf.String(), nil
}

func main() {
	err := initRuntimeAndModule()
	if err != nil {
		panic(err)
	}
	defer runtime.Destroy()
	err = mapCalls()
	if err != nil {
		panic(err)
	}
	jsInput := "var s=\"test\"; typeof(s)"
	inputLength := len(jsInput)
	printf(("JS Input is: %s\n"), jsInput)

	ptr, err := allocate(jsInput)
	printf("Allocated %d bytes in WASM memory (\"boa_alloc\"), pointer is %d\n", inputLength, ptr)

	printf("Calling \"boa_exec3\" with arguments: (ptr=%d, length=%d)\n", ptr, inputLength)
	out, err := exec(ptr, inputLength)
	if err != nil {
		panic(err)
	}
	printf("JS output is: %s\n", out)
}
