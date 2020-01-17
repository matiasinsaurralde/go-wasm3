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
	env := wasm3.NewEnvironment()
	defer env.Destroy()
	runtime := wasm3.NewRuntime(env, 64*1024)
	defer runtime.Destroy()
	_, err := runtime.Load(wasmBytes)
	if err != nil {
		t.Fatal(err)
	}
	fn, err := runtime.FindFunction(fnName)
	if err != nil {
		t.Fatal(err)
	}
	result := fn()
	memoryLength := runtime.GetAllocatedMemoryLength()

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
	if buf.String() != "testingonly" {
		t.Fatal("Reconstructed string doesn't match")
	}
}

func BenchmarkCString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		env := wasm3.NewEnvironment()
		defer env.Destroy()
		runtime := wasm3.NewRuntime(env, 64*1024)
		defer runtime.Destroy()
		_, err := runtime.Load(wasmBytes)
		if err != nil {
			b.Fatal(err)
		}
		fn, err := runtime.FindFunction(fnName)
		if err != nil {
			b.Fatal(err)
		}
		result := fn()
		memoryLength := runtime.GetAllocatedMemoryLength()

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
		if buf.String() != "testingonly" {
			b.Fatal("Reconstructed string doesn't match")
		}
	}
}

func BenchmarkCStringReentrant(b *testing.B) {
	env := wasm3.NewEnvironment()
	defer env.Destroy()
	runtime := wasm3.NewRuntime(env, 64*1024)
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
		result := fn()
		memoryLength := runtime.GetAllocatedMemoryLength()

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
		if buf.String() != "testingonly" {
			b.Fatal("Reconstructed string doesn't match")
		}
	}
}
