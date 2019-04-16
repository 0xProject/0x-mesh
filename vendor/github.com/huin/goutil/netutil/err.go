package netutil

import (
	"strings"
)

type MultiError []error

func (e *MultiError) RecordError(err error) {
	if err != nil {
		*e = append(*e, err)
	}
}

func (e MultiError) ToError() error {
	if len(e) > 0 {
		return e
	} else if len(e) == 1 {
		return e[0]
	}
	return nil
}

func (e MultiError) Error() string {
	strs := make([]string, len(e))
	for i, err := range e {
		strs[i] = err.Error()
	}
	return "multiple errors: " + strings.Join(strs, ", ")
}
