package errutil

import (
	"bytes"
	"fmt"
)

type Errors []error

func NewErrors(errs []error) error {
	if len(errs) == 0 {
		panic("empty errors must be not an error")
	}
	return Errors(errs)
}

var _ error = Errors{}

func (s Errors) Error() string {
	buf := &bytes.Buffer{}
	for _, e := range s {
		_, _ = fmt.Fprintln(buf, e.Error())
	}
	return buf.String()
}
