package bitarray

import (
	"testing"
)

func TestReadBits(t *testing.T) {
	ba := New([]byte{0xf0, 0x01}, 16)
	r := NewReader(ba)

	if r.Pos() != 0 {
		t.Errorf("expected pos %d got %d", 0, r.Pos())
	}

	var test uint
	if err := r.ReadBits(&test, 4); err != nil {
		t.Errorf("failed with %q", err)
	}
	if test != 15 {
		t.Errorf("expected 15 got %d", test)
	}
	if r.Pos() != 4 {
		t.Errorf("expected pos %d got %d", 4, r.Pos())
	}
	r.Seek(-4, SeekCurrent)
	if r.Pos() != 0 {
		t.Errorf("expected pos %d got %d", 0, r.Pos())
	}

	r.Seek(4, SeekStart)
	if r.Pos() != 4 {
		t.Errorf("expected pos %d got %d", 4, r.Pos())
	}
	if err := r.ReadBits(&test, 4); err != nil {
		t.Errorf("failed with %q", err)
	}
	if test != 0 {
		t.Errorf("expected 0 got %d", test)
	}
	if r.Pos() != 8 {
		t.Errorf("expected pos %d got %d", 8, r.Pos())
	}

	if err := r.ReadBits(&test, 8); err != nil {
		t.Errorf("failed with %q", err)
	}
	if test != 1 {
		t.Errorf("expected 1 got %d", test)
	}
	if r.Pos() != 16 {
		t.Errorf("expected pos %d got %d", 16, r.Pos())
	}

	// Read past end
	if err := r.ReadBits(&test, 1); err == nil {
		t.Errorf("expected error")
	}

	r.Seek(1, SeekEnd)
	if r.Pos() != 15 {
		t.Errorf("expected pos %d got %d", 15, r.Pos())
	}
	if !r.ReadBit() {
		t.Errorf("expected 1 got 0")
	}
	if r.Pos() != 16 {
		t.Errorf("expected pos %d got %d", 16, r.Pos())
	}

	// Some errors
	if _, err := r.Seek(1, 4); err == nil {
		t.Errorf("expected err got none")
	}
	if _, err := r.Seek(17, SeekStart); err == nil {
		t.Errorf("expected err got none")
	}
	if _, err := r.Seek(-1, SeekStart); err == nil {
		t.Errorf("expected err got none")
	}
}
