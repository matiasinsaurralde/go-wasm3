libxml
==

More details [here](https://github.com/matiasinsaurralde/wasm-libxml2).

Benchmark output:

```
% go test -bench=. -benchmem -v
# github.com/matiasinsaurralde/go-wasm3/examples/libxml.test
=== RUN   TestXMLValidation
--- PASS: TestXMLValidation (0.01s)
goos: darwin
goarch: amd64
pkg: github.com/matiasinsaurralde/go-wasm3/examples/libxml
BenchmarkXMLValidation/Good_XML-8         	    1000	   1498826 ns/op	      88 B/op	       5 allocs/op
BenchmarkXMLValidation/Bad_XML-8          	    1000	   2158523 ns/op	     104 B/op	       6 allocs/op
PASS
ok  	github.com/matiasinsaurralde/go-wasm3/examples/libxml	4.123s
```