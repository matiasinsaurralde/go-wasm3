go-wasm3
==

Golang wrapper for [WASM3](https://github.com/wasm3/wasm3), WIP.

This is part of a series of WASM-related experiments: [go-wavm](https://github.com/matiasinsaurralde/go-wavm) and [go-wasm-benchmark](https://github.com/matiasinsaurralde/go-wasm-benchmark).

## Install/build

For now I've attached two static builds of [WASM3](https://github.com/wasm3/wasm3) for OS X and Linux (64 bits). If you want to hack around the library, you will need [this CMake tweak](https://github.com/matiasinsaurralde/wasm3/commit/824cb245617ad9888e1b36c47c164d5c687cd272).

If you don't want to build [WASM3](https://github.com/wasm3/wasm3) and one of the mentioned platforms is in use, `go get` should be enough:

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

## Limitations and future

This is a WIP, the sample that was described shows only a very basic flow. It's not yet possible to access the instance memory or play around module imports/exports but I'm working on it. Stay tuned!