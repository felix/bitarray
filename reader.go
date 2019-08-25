package bitarray

import "errors"

// EOF is the error returned by Read when no more input is available.
var (
	EOF              = errors.New("EOF")
	ErrInvalidWhence = errors.New("invalid whence")
	ErrInvalidOffset = errors.New("invalid offset")
)

const (
	SeekStart   = 0 // seek relative to the origin of the file
	SeekCurrent = 1 // seek relative to the current offset
	SeekEnd     = 2 // seek relative to the end
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

// ReadBit advances one bit and returns the result of a boolean AND test on it.
func (r *Reader) ReadBit() bool {
	out := r.ba.Test(r.i)
	r.i++
	return out
}

// Seek sets the internal pointer to position n.
func (r *Reader) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case 0:
		r.i = offset
	case 1:
		r.i += offset
	case 2:
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
