package main

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	inputPath  = "./wasm/main.wasm"
	outputPath = "./ts/generated/wasm_buffer.ts"
)

var (
	prefix = []byte("import * as base64 from 'base64-arraybuffer';\nexport const wasmBuffer = base64.decode('")
	suffix = []byte("');\n")
)

func main() {
	wasmBytcode, err := ioutil.ReadFile(inputPath)
	if err != nil {
		panic(err)
	}

	encodedLen := base64.StdEncoding.EncodedLen(len(wasmBytcode))
	encodedWasmBytcode := make([]byte, encodedLen)
	base64.StdEncoding.Encode(encodedWasmBytcode, wasmBytcode)

	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		panic(err)
	}

	// HACK(albrow): We generate a TypeScript file that contains the Wasm output
	// encoded as a base64 string. This is the most reliable way to load Wasm such
	// that users just see a TypeScript/JavaScript package and without relying on
	// a third-party server.
	outputFile, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	mustWrite(outputFile, prefix)
	mustWrite(outputFile, encodedWasmBytcode)
	mustWrite(outputFile, suffix)
}

func mustWrite(writer io.Writer, data []byte) {
	if _, err := writer.Write(data); err != nil {
		panic(err)
	}
}
