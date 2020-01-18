package main

import (
	"bytes"
	"io/ioutil"
	"log"

	wasm3 "github.com/matiasinsaurralde/go-wasm3"
)

const (
	wasmFilename = "boa.wasm"
)

func main() {
	env := wasm3.NewEnvironment()
	defer env.Destroy()
	runtime := wasm3.NewRuntime(&wasm3.Config{
		Environment: env,
		StackSize:   1024 * 1024,
		EnableWASI:  false,
	})
	defer runtime.Destroy()
	log.Println("Runtime loaded")

	wasmBytes, err := ioutil.ReadFile(wasmFilename)
	if err != nil {
		panic(err)
	}

	module, err := env.ParseModule(wasmBytes)
	if err != nil {
		panic(err)
	}
	_, err = runtime.LoadModule(module)
	if err != nil {
		panic(err)
	}
	log.Print("Module loaded")

	// Map exported functions:
	allocateFn, err := runtime.FindFunction("boa_alloc")
	if err != nil {
		panic(err)
	}

	execFn, err := runtime.FindFunction("boa_exec3")
	if err != nil {
		panic(err)
	}

	jsInput := "var s=\"test\"; typeof(s)"
	log.Printf(("JS Input is: %s\n"), jsInput)
	length := len(jsInput)
	ptr, err := allocateFn(length)
	if err != nil {
		panic(err)
	}
	pos := ptr
	for _, ch := range jsInput {
		runtime.Memory()[pos] = byte(ch)
		pos++
	}
	log.Printf("Allocated %d bytes in WASM memory (\"boa_alloc\"), pointer is %d\n", length, ptr)

	log.Printf("Calling \"boa_exec3\" with arguments: (ptr=%d, length=%d)\n", ptr, length)
	outPtr, err := execFn(ptr, length)
	if err != nil {
		panic(err)
	}
	log.Printf("\"boa_exec3\" returned, output pointer is %d\n", outPtr)
	i := 0
	buf := new(bytes.Buffer)
	for {
		pos = outPtr + i
		ch := runtime.Memory()[pos]
		i++
		if ch == 0 {
			break
		}
		buf.WriteByte(ch)
	}
	log.Printf("Read %d bytes from WASM memory, starting in %d\n", buf.Len(), outPtr)
	outStr := buf.String()
	log.Printf("JS output is: %s\n", outStr)
}
