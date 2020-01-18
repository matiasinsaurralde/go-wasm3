package wasm3

/*
#include "go-wasm3.h"
*/
import "C"

var(
	lastError string
)

//export set_error
func set_error(str ResultT) {
	lastError = C.GoString(str)
}

// LastErrorString returns the last runtime error
func LastErrorString() string {
	return lastError
}
