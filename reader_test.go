package bitarray

import (
	"testing"
)

func TestReadBits(t *testing.T) {
	ba := NewFromBytes([]byte{0xf0, 0x01}, 16)
	r := NewReader(ba)

	if r.Pos() != 0 {
		t.Errorf("got pos %d, want %d", r.Pos(), 0)
	}

	var test uint
	if err := r.ReadBits(&test, 4); err != nil {
		t.Errorf("failed with %q", err)
	}
	if test != 15 {
		t.Errorf("got 15, want %d", test)
	}
	if r.Pos() != 4 {
		t.Errorf("got pos %d, want %d", r.Pos(), 4)
	}
	r.Seek(-4, SeekCurrent)
	if r.Pos() != 0 {
		t.Errorf("got pos %d, want %d", r.Pos(), 0)
	}

	r.Seek(4, SeekStart)
	if r.Pos() != 4 {
		t.Errorf("got pos %d, want %d", r.Pos(), 4)
	}
	if err := r.ReadBits(&test, 4); err != nil {
		t.Errorf("failed with %q", err)
	}
	if test != 0 {
		t.Errorf("got 0, want %d", test)
	}
	if r.Pos() != 8 {
		t.Errorf("got pos %d, want %d", r.Pos(), 8)
	}

	if err := r.ReadBits(&test, 8); err != nil {
		t.Errorf("failed with %q", err)
	}
	if test != 1 {
		t.Errorf("got 1, want %d", test)
	}
	if r.Pos() != 16 {
		t.Errorf("got pos %d, want %d", r.Pos(), 16)
	}

	// Read past end
	if err := r.ReadBits(&test, 1); err == nil {
		t.Errorf("expected error")
	}

	r.Seek(1, SeekEnd)
	if r.Pos() != 15 {
		t.Errorf("got pos %d, want %d", r.Pos(), 15)
	}
	if !r.ReadBit() {
		t.Errorf("got 1, want 0")
	}
	if r.Pos() != 16 {
		t.Errorf("got pos %d, want %d", r.Pos(), 16)
	}

	// Some errors
	if _, err := r.Seek(1, 4); err == nil {
		t.Errorf("expected err got none")
	}
	if _, err := r.Seek(17, SeekStart); err == nil {
		t.Errorf("got none, want err")
	}
	if _, err := r.Seek(-1, SeekStart); err == nil {
		t.Errorf("got none, want err")
	}
}
