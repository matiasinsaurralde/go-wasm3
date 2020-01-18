package main

import (
	"io/ioutil"
	golog "log"

	wasm3 "github.com/matiasinsaurralde/go-wasm3"
)

var (
	allocateFn        wasm3.FunctionWrapper
	newSchemaParserFn wasm3.FunctionWrapper
	validateFn        wasm3.FunctionWrapper
	runtime           *wasm3.Runtime
	print             = golog.Print
	printf            = golog.Printf
)

const (
	wasmFilename = "libxml2.wasm"
)

func initRuntimeAndModule() error {
	runtime = wasm3.NewRuntime(&wasm3.Config{
		Environment: wasm3.NewEnvironment(),
		StackSize:   1024 * 1024,
		EnableWASI:  true,
	})

	wasmBytes, err := ioutil.ReadFile(wasmFilename)
	if err != nil {
		return err
	}

	module, err := runtime.ParseModule(wasmBytes)
	if err != nil {
		return err
	}
	_, err = runtime.LoadModule(module)
	if err != nil {
		return err
	}
	return nil
}

func mapCalls() error {
	var err error
	allocateFn, err = runtime.FindFunction("wasm_allocate")
	if err != nil {
		return err
	}

	newSchemaParserFn, err = runtime.FindFunction("wasm_new_schema_parser2")
	if err != nil {
		return err
	}

	validateFn, err = runtime.FindFunction("wasm_validate_xml")
	if err != nil {
		return err
	}
	return nil
}

func allocate(input []byte) (int, error) {
	ptr, err := allocateFn(len(input))
	if err != nil {
		return 0, nil
	}
	pos := ptr
	for _, ch := range input {
		runtime.Memory()[pos] = byte(ch)
		pos++
	}
	return ptr, nil
}

func newSchemaParser(ptr, length int) (int, error) {
	outPtr, err := newSchemaParserFn(ptr, length)
	if err != nil {
		return 0, err
	}
	return outPtr, nil
}

func validate(xmlPtr, xmlLength, schemaParserPtr int) (int, error) {
	out, err := validateFn(xmlPtr, xmlLength, schemaParserPtr)
	return out, err
}

func main() {
	err := initRuntimeAndModule()
	if err != nil {
		panic(err)
	}
	defer runtime.Destroy()
	err = mapCalls()
	if err != nil {
		panic(err)
	}
	xsdFile, err := ioutil.ReadFile("input.xsd")
	if err != nil {
		panic(err)
	}
	xsdPtr, err := allocate(xsdFile)
	if err != nil {
		panic(err)
	}
	printf("Allocated %d bytes in WASM memory (\"wasm_allocate\"), pointer is %d (XSD file)\n", len(xsdFile), xsdPtr)
	xmlFile, err := ioutil.ReadFile("input.xml")
	if err != nil {
		panic(err)
	}
	xmlPtr, err := allocate(xmlFile)
	if err != nil {
		panic(err)
	}
	printf("Allocated %d bytes in WASM memory (\"wasm_allocate\"), pointer is %d (XSD file)\n", len(xmlFile), xmlPtr)
	schemaParserPtr, err := newSchemaParser(xsdPtr, len(xsdFile))
	if err != nil {
		panic(err)
	}
	printf("Schema parser was initialized (\"wasm_new_schema_parser2\"), pointer is %d\n", schemaParserPtr)

	out, err := validate(xmlPtr, len(xmlFile), schemaParserPtr)
	if err != nil {
		panic(err)
	}
	printf("\"wasm_validate_xml\" output is: %d", out)
}
