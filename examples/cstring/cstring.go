package main

import (
	"bytes"
	"io/ioutil"
	"log"

	wasm3 "github.com/matiasinsaurralde/go-wasm3"
)

const (
	wasmFilename = "cstring/cstring.wasm"
	fnName       = "somecall"
)

func main() {
	log.Print("Initializing WASM3")

	env := wasm3.NewEnvironment()
	defer env.Destroy()
	runtime := wasm3.NewRuntime(env, 64*1024)
	defer runtime.Destroy()
	log.Println("Runtime ok")

	wasmBytes, err := ioutil.ReadFile(wasmFilename)
	if err != nil {
		panic(err)
	}
	log.Printf("Read WASM module (%d bytes)\n", len(wasmBytes))

	module, err := env.ParseModule(wasmBytes)
	if err != nil {
		panic(err)
	}
	module, err = runtime.LoadModule(module)
	if err != nil {
		panic(err)
	}
	log.Print("Loaded module")
	fn, err := runtime.FindFunction(fnName)
	if err != nil {
		panic(err)
	}
	log.Printf("Found '%s' function (using runtime.FindFunction)", fnName)
	memoryLength := runtime.GetAllocatedMemoryLength()
	log.Printf("Allocated memory (before function call) is: %d\n", memoryLength)
	result := fn()
	memoryLength = runtime.GetAllocatedMemoryLength()
	log.Printf("Allocated memory (after function call) is: %d\n", memoryLength)

	// Reconstruct the string from memory:
	mem := runtime.GetMemory(memoryLength, 0)
	buf := new(bytes.Buffer)
	for n := 0; n < memoryLength; n++ {
		if n < result {
			continue
		}
		value := mem[n]
		if value == 0 {
			break
		}
		buf.WriteByte(value)
	}
	log.Printf("Buffer length is: %d\n", buf.Len())
	log.Printf("Buffer contains: %s\n", buf.String())
}
