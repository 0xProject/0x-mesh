package codegen

import (
	"io"
	"os"
	"os/exec"
)

type goFmtWriteCloser struct {
	output io.WriteCloser
	stdin  io.WriteCloser
	gofmt  *exec.Cmd
}

// NewGofmtWriteCloser returns an io.WriteCloser that filters what is written
// to it through gofmt. It must be closed for this process to be completed, an
// error from Close can be due to syntax errors in the source that has been
// written.
func NewGofmtWriteCloser(output io.WriteCloser) (io.WriteCloser, error) {
	gofmt := exec.Command("gofmt")
	gofmt.Stdout = output
	gofmt.Stderr = os.Stderr
	stdin, err := gofmt.StdinPipe()
	if err != nil {
		return nil, err
	}
	if err = gofmt.Start(); err != nil {
		return nil, err
	}
	return &goFmtWriteCloser{
		output: output,
		stdin:  stdin,
		gofmt:  gofmt,
	}, nil
}

func (gwc *goFmtWriteCloser) Write(p []byte) (int, error) {
	return gwc.stdin.Write(p)
}

func (gwc *goFmtWriteCloser) Close() error {
	gwc.stdin.Close()
	if err := gwc.output.Close(); err != nil {
		gwc.gofmt.Wait()
		return err
	}
	return gwc.gofmt.Wait()
}
