package main

import (
	"testing"
)

func init() {
	print = func(...interface{}) {}
	printf = func(string, ...interface{}) {}
	err := initRuntimeAndModule()
	if err != nil {
		panic(err)
	}
	err = mapCalls()
	if err != nil {
		panic(err)
	}
}

func boaCall(t testing.TB) {
	jsInput := "var s=\"test\"; typeof(s)"
	inputLength := len(jsInput)

	ptr, err := allocate(jsInput)
	out, err := exec(ptr, inputLength)
	if err != nil {
		t.Fatal(err)
	}
	if out != "string" {
		t.Fatalf("Unexpected output: %s", out)
	}
}

func TestBoaCall(t *testing.T) {
	boaCall(t)
}

func BenchmarkBoaCall(b *testing.B) {
	for n := 0; n < b.N; n++ {
		boaCall(b)
	}
}
