// +build !js

package keys

import (
	"io/ioutil"
	"os"
)

func readFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func mkdirAll(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

func writeFile(path string, data []byte) error {
	return ioutil.WriteFile(path, data, os.ModePerm)
}
