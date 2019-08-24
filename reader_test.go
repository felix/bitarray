package bitarray

import (
	"testing"
)

func TestReadBits(t *testing.T) {
	ba := New([]byte{0xf0, 0x01}, 16)
	r := NewReader(ba)

	var test uint
	if err := r.ReadBits(&test, 4); err != nil {
		t.Errorf("failed with %q", err)
	}
	if test != 15 {
		t.Errorf("expected 15 got %d", test)
	}

	if err := r.ReadBits(&test, 4); err != nil {
		t.Errorf("failed with %q", err)
	}
	if test != 0 {
		t.Errorf("expected 0 got %d", test)
	}

	if err := r.ReadBits(&test, 8); err != nil {
		t.Errorf("failed with %q", err)
	}
	if test != 1 {
		t.Errorf("expected 1 got %d", test)
	}
	t.Log(r.i)

	// Read past end
	if err := r.ReadBits(&test, 1); err == nil {
		t.Errorf("expected error")
	}
}
