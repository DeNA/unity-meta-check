package meta

import (
	"encoding/hex"
	"fmt"
)

type GUID struct {
	bytes []byte
}

const GUIDByteLength = 16

func NewGUID(bytes []byte) (*GUID, error) {
	if len(bytes) != GUIDByteLength {
		return nil, fmt.Errorf("length of GUID must be %d bytes", GUIDByteLength)
	}
	return &GUID{bytes}, nil
}

func (g GUID) String() string {
	s := make([]byte, hex.EncodedLen(GUIDByteLength))
	hex.Encode(s, g.bytes)
	return string(s)
}
