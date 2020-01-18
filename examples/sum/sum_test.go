package main

import (
	"io/ioutil"
	"testing"

	wasm3 "github.com/matiasinsaurralde/go-wasm3"
)

var (
	wasmBytes []byte
)

func init() {
	var err error
	wasmBytes, err = ioutil.ReadFile(wasmFilename)
	if err != nil {
		panic(err)
	}
}

func TestSum(t *testing.T) {
	runtime := wasm3.NewRuntime(&wasm3.Config{
		Environment: wasm3.NewEnvironment(),
		StackSize:   64 * 1024,
	})
	defer runtime.Destroy()
	_, err := runtime.Load(wasmBytes)
	if err != nil {
		t.Fatal(err)
	}
	fn, err := runtime.FindFunction(fnName)
	if err != nil {
		t.Fatal(err)
	}
	result, _ := fn(1, 1)
	if result != 2 {
		t.Fatal("Result doesn't match")
	}
}

func BenchmarkSum(b *testing.B) {
	for n := 0; n < b.N; n++ {
		runtime := wasm3.NewRuntime(&wasm3.Config{
			Environment: wasm3.NewEnvironment(),
			StackSize:   64 * 1024,
		})
		defer runtime.Destroy()
		_, err := runtime.Load(wasmBytes)
		if err != nil {
			b.Fatal(err)
		}
		fn, err := runtime.FindFunction(fnName)
		if err != nil {
			b.Fatal(err)
		}
		fn(1, 2)
	}
}

func BenchmarkSumReentrant(b *testing.B) {
	runtime := wasm3.NewRuntime(&wasm3.Config{
		Environment: wasm3.NewEnvironment(),
		StackSize:   64 * 1024,
	})
	defer runtime.Destroy()
	_, err := runtime.Load(wasmBytes)
	if err != nil {
		b.Fatal(err)
	}
	fn, err := runtime.FindFunction(fnName)
	if err != nil {
		b.Fatal(err)
	}
	for n := 0; n < b.N; n++ {
		fn(1, 2)
	}
}
