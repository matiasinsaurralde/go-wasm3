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

	_, err = runtime.Load(wasmBytes)
	if err != nil {
		panic(err)
	}
	log.Print("Module loaded")

	fn, err := runtime.FindFunction(fnName)
	if err != nil {
		panic(err)
	}

	log.Println("Calling function")
	fn("1", "2")
	fn("2", "2")
}
