package main

import (
	"bytes"
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

func TestCString(t *testing.T) {
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
	result, _ := fn()
	memoryLength := runtime.GetAllocatedMemoryLength()

	// Reconstruct the string from memory:
	mem := runtime.Memory()
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
	if buf.String() != "testingonly" {
		t.Fatal("Reconstructed string doesn't match")
	}
}

func BenchmarkCString(b *testing.B) {
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
		result, _ := fn()
		memoryLength := runtime.GetAllocatedMemoryLength()

		// Reconstruct the string from memory:
		mem := runtime.Memory()
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
		if buf.String() != "testingonly" {
			b.Fatal("Reconstructed string doesn't match")
		}
	}
}

func BenchmarkCStringReentrant(b *testing.B) {
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
		result, _ := fn()
		memoryLength := runtime.GetAllocatedMemoryLength()

		// Reconstruct the string from memory:
		mem := runtime.Memory()
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
		if buf.String() != "testingonly" {
			b.Fatal("Reconstructed string doesn't match")
		}
	}
}
