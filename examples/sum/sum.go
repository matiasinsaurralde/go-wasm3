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
	result := fn(1, 1)
	log.Print("Result is: ", result)

	// Different call approach, retrieving functions from the module object:
	fn2, err := module.GetFunctionByName("sum")
	if err != nil {
		panic(err)
	}
	log.Printf("Found '%s' function (using module.GetFunctionByName)", fnName)
	result = fn2.Call(2, 2)
	log.Print("Result is: ", result)
}
