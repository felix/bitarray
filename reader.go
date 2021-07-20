package bitarray

import (
	"errors"
)

var (
	// EOF is returned when no more input is available.
	EOF = errors.New("EOF") // revive:disable:error-naming
	// ErrInvalidWhence is return for invalid seeking.
	ErrInvalidWhence = errors.New("invalid whence")
	// ErrInvalidOffset is returned for invalid seek offsets.
	ErrInvalidOffset = errors.New("invalid offset")
)

type SeekFrom int

const (
	// SeekStart seeks relative to the origin of the file
	SeekStart SeekFrom = iota
	// SeekCurrent seeks relative to the current offset
	SeekCurrent = 1
	// SeekEnd seeks relative to the end
	SeekEnd = 2
)

// A Reader implements the BitReader interface.
type Reader struct {
	ba *BitArray
	i  int64
}

// NewReader creates a new Reader.
func NewReader(ba *BitArray) *Reader {
	return &Reader{ba: ba, i: 0}
}

// ReadBits reads n bits from the BitArray into out.
func (r *Reader) ReadBits(out *uint, n int) error {
	if r.i >= int64(r.ba.size) {
		return EOF
	}
	// TODO this will have issues on 32bit systems
	i, err := r.ba.ReadUint(r.i, int64(n))
	if err != nil {
		return err
	}
	*out = i
	r.i += int64(n)
	return nil
}

// Pos returns the current position of the reader.
func (r *Reader) Pos() int64 {
	return r.i
}

// ReadBit advances one bit and returns the result of a boolean AND test on it.
func (r *Reader) ReadBit() bool {
	out := r.ba.Test(r.i)
	r.i++
	return out
}

// Seek sets the internal pointer to position n. It returns the resulting offset.
func (r *Reader) Seek(offset int64, whence SeekFrom) (int64, error) {
	switch whence {
	case SeekStart:
		r.i = offset
	case SeekCurrent:
		r.i += offset
	case SeekEnd:
		r.i = int64(r.ba.size) - offset
	default:
		return r.i, ErrInvalidWhence
	}
	if r.i < 0 {
		return r.i, ErrInvalidOffset
	}
	if r.i > r.ba.size {
		return r.i, EOF
	}
	return r.i, nil
}
