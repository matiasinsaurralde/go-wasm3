package main

import (
	"io/ioutil"
	"testing"
)

var (
	badXMLData, _  = ioutil.ReadFile("bad_input.xml")
	goodXMLData, _ = ioutil.ReadFile("input.xml")
	xsdData, _     = ioutil.ReadFile("input.xsd")

	badXMLPtr       int
	goodXMLPtr      int
	xsdPtr          int
	schemaParserPtr int
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

	badXMLPtr, err = allocate(badXMLData)
	if err != nil {
		panic(err)
	}
	goodXMLPtr, err = allocate(goodXMLData)
	if err != nil {
		panic(err)
	}
	xsdPtr, err = allocate(xsdData)
	if err != nil {
		panic(err)
	}
	schemaParserPtr, err = newSchemaParser(xsdPtr, len(xsdData))
	if err != nil {
		panic(err)
	}
}

func validateGoodXML() (int, error) {
	return validate(goodXMLPtr, len(goodXMLData), schemaParserPtr)
}

func validateBadXML() (int, error) {
	return validate(badXMLPtr, len(badXMLData), schemaParserPtr)
}

func TestXMLValidation(t *testing.T) {
	i, err := validateGoodXML()
	if i > 0 && err != nil {
		t.Fatal("Unexpected output, result should be a positive number, err should be nil")
	}
	i, _ = validateBadXML()
	if i != -1 {
		t.Fatal("Unexpected output, result should be -1")
	}
}
func BenchmarkXMLValidation(b *testing.B) {
	b.Run("Good XML", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			validateGoodXML()
		}
	})
	b.Run("Bad XML", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			validateBadXML()
		}
	})
}
