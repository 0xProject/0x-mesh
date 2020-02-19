package main

import (
	"encoding/base64"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/plaid/go-envvar/envvar"
)

type EnvVars struct {
	InputPath  string `envvar:"INPUT_PATH"`
	OutputPath string `envvar:"OUTPUT_PATH"`
}

var (
	prefix = []byte("import * as base64 from 'base64-arraybuffer';\nexport const wasmBuffer = base64.decode('")
	suffix = []byte("');\n")
)

func main() {
	env := EnvVars{}
	if err := envvar.Parse(&env); err != nil {
		panic(err)
	}

	wasmBytcode, err := ioutil.ReadFile(env.InputPath)
	if err != nil {
		panic(err)
	}

	encodedLen := base64.StdEncoding.EncodedLen(len(wasmBytcode))
	encodedWasmBytcode := make([]byte, encodedLen)
	base64.StdEncoding.Encode(encodedWasmBytcode, wasmBytcode)

	outputDir := filepath.Dir(env.OutputPath)
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		panic(err)
	}

	// HACK(albrow): We generate a TypeScript file that contains the Wasm output
	// encoded as a base64 string. This is the most reliable way to load Wasm such
	// that users just see a TypeScript/JavaScript package and without relying on
	// a third-party server.
	outputFile, err := os.OpenFile(env.OutputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
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
