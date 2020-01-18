package main

import (
	"io/ioutil"
	"log"

	wasm3 "github.com/matiasinsaurralde/go-wasm3"
)

const (
	wasmFilename = "sum.wasm"
	fnName       = "sum"
)

func main() {
	log.Print("Initializing WASM3")

	runtime := wasm3.NewRuntime(&wasm3.Config{
		Environment: wasm3.NewEnvironment(),
		StackSize:   64 * 1024,
	})
	log.Println("Runtime ok")

	wasmBytes, err := ioutil.ReadFile(wasmFilename)
	if err != nil {
		panic(err)
	}
	log.Printf("Read WASM module (%d bytes)\n", len(wasmBytes))

	module, err := runtime.ParseModule(wasmBytes)
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
	result, _ := fn(1, 1)
	log.Print("Result is: ", result)

	// Different call approach, retrieving functions from the module object:
	fn2, err := module.GetFunctionByName("sum")
	if err != nil {
		panic(err)
	}
	log.Printf("Found '%s' function (using module.GetFunctionByName)", fnName)
	result, _ = fn2.Call(2, 2)
	log.Print("Result is: ", result)
}
