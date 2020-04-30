go-wasm3
==

[![GoDoc](https://godoc.org/github.com/matiasinsaurralde/go-wasm3?status.svg)](https://godoc.org/github.com/matiasinsaurralde/go-wasm3)
[![Build status](https://circleci.com/gh/matiasinsaurralde/go-wasm3/tree/master.svg?style=shield)](https://circleci.com/gh/matiasinsaurralde/go-wasm3/tree/master)

Golang wrapper for [WASM3](https://github.com/wasm3/wasm3), WIP.

This is part of a series of WASM-related experiments: [go-wavm](https://github.com/matiasinsaurralde/go-wavm) and [go-wasm-benchmark](https://github.com/matiasinsaurralde/go-wasm-benchmark).

## Install/build

This package ships with pre-built [`WASM3`](https://github.com/wasm3/wasm3) libraries (static builds) for OS X and Linux. If you want to hack around it, check [the original repository](https://github.com/wasm3/wasm3).

If you're using one of the mentioned platforms, you may install the package using `go get`:

```
$ go get -u github.com/matiasinsaurralde/go-wasm3
```

To inspect or run the little sample use:

```
$ cd $GOPATH/src/github.com/matiasinsaurralde/go-wasm3/examples/sum
$ go build # or "go run sum.go"
$ ./sum
```

The output will look as follows:

```
2020/01/15 09:51:24 Initializing WASM3
2020/01/15 09:51:24 Runtime ok
2020/01/15 09:51:24 Read WASM module (139 bytes)
2020/01/15 09:51:24 Module loaded
2020/01/15 09:51:24 Calling function
Result: 3
Result: 4
```

You will find additional sample projects in the next section.

## Sample projects

### boa

This program uses the [boa](https://github.com/jasonwilliams/boa) engine to evaluate JS code. Boa is an embeddable JS engine written in Rust, for this sample it was compiled targeting WASM.

Link [here](https://github.com/matiasinsaurralde/go-wasm3/tree/master/examples/boa).

### libxml

This program loads [libxml2](https://github.com/GNOME/libxml2) as a WASM module (it's a custom build, full instructions [here](https://github.com/matiasinsaurralde/wasm-libxml2)). The library is used to validate a XML file against a XSD (both loaded from the Go side).

Link [here](https://github.com/matiasinsaurralde/go-wasm3/tree/master/examples/libxml).


## Memory access

Take the following sample program:

```c
#include <stdlib.h>
#include <string.h>

char* somecall() {
    // Allocate a few bytes on the heap:
    char* test = (char*) malloc(12*sizeof(char));

    // Copy a string into the previously defined address:
    strcpy(test, "testingonly");

    // Return the pointer:
    return test;
}
```

Build it using `wasicc`, this will generate a `cstring.wasm` file (WASM module):

```
wasicc cstring.c -Wl,--export-all -o cstring.wasm
```

The following Go code will load the WASM module and retrieve the data after calling `somecall`:

```go
    // Initialize the runtime and load the module:
    env := wasm3.NewEnvironment()
	defer env.Destroy()
	runtime := wasm3.NewRuntime(env, 64*1024)
	defer runtime.Destroy()
    wasmBytes, err := ioutil.ReadFile("program.wasm")
	module, _ := env.ParseModule(wasmBytes)
	runtime.LoadModule(module)
    fn, _ := runtime.FindFunction(fnName)

    // Call somecall and get the pointer to our data:
    result := fn()
    
    // Reconstruct the string from memory:
    memoryLength = runtime.GetAllocatedMemoryLength()
    mem := runtime.GetMemory(memoryLength, 0)
    
    // Initialize a Go buffer:
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

    // Print the string: "testingonly"
    str := buf.String()
    fmt.Println(str)
```

For more details check [this](https://github.com/matiasinsaurralde/go-wasm3/tree/master/examples/cstring).

## Limitations and future

This is a WIP. Stay tuned!

## Related projects

A Rust wrapper is available [here](https://github.com/Veykril/wasm3-rs).

## License

[MIT](https://github.com/matiasinsaurralde/go-wasm3/blob/master/LICENSE).

[wasm3](https://github.com/wasm3/wasm3/blob/master/LICENSE) is also under this license.
