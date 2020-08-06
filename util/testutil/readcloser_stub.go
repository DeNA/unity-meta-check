package testutil

import (
	"errors"
)

type ErrorReadCloserStub struct{}

func (*ErrorReadCloserStub) Close() error {
	return errors.New("close: EXPECTED_FAILURE")
}

func (*ErrorReadCloserStub) Read(_ []byte) (n int, err error) {
	return 0, errors.New("read: EXPECTED_FAILURE")
}
