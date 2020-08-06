package meta

import (
	"testing"
)

func TestNewGUID(t *testing.T) {
	guid, err := NewGUID([]byte{0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf})
	if err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}

	actual := guid.String()
	expected := "000102030405060708090a0b0c0d0e0f"
	if actual != expected {
		t.Errorf("want %q, got %q", expected, actual)
		return
	}
}

func TestRandomGUIDGenerator(t *testing.T) {
	guidGen := RandomGUIDGenerator()
	_, err := guidGen()
	if err != nil {
		t.Errorf("want nil, got %#v", err)
		return
	}
}
