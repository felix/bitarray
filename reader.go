package bitarray

import "errors"

// EOF is the error returned by Read when no more input is available.
var EOF = errors.New("EOF")

// BitReader is the interface that wraps the ReadBits method.
type BitReader interface {
	ReadBits(*uint, uint) error
}

// A Reader implements the BitReader interface.
type Reader struct {
	ba *BitArray
	i  uint
}

// NewReader creates a new Reader.
func NewReader(ba *BitArray) *Reader {
	return &Reader{ba: ba, i: 0}
}

// ReadBits reads n bits from the BitArray into out.
func (r *Reader) ReadBits(out *uint, n uint) error {
	if r.i >= uint(r.ba.size) {
		return EOF
	}
	// TODO this will have issues on 32bit systems
	i, err := r.ba.ReadUint(int(r.i), int(n))
	if err != nil {
		return err
	}
	*out = i
	r.i += n
	return nil
}

// ReadBit advances one bit and returns the result of a boolean AND test on it.
func (r *Reader) ReadBit() bool {
	out := r.ba.Test(int(r.i))
	r.i++
	return out
}

// Seek sets the internal pointer to position n.
func (r *Reader) Seek(n uint) {
	r.i = n
}
