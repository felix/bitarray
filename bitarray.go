package bitarray

import (
	"fmt"
	"math/big"
	"math/bits"
	"strings"
)

// BitArray holds an array of bits.
type BitArray struct {
	raw  []byte
	size int64
}

// NewFromBytes creates a BitArray from the given []byte and count bits.
// Count is from position 0 of b.
func NewFromBytes(b []byte, count int64) *BitArray {
	out := &BitArray{}
	out.raw = b
	out.size = count
	//out.avail = uint(len(b)*8 - count)
	return out
}

const uintSize = 32 << (^uint(0) >> 63)

// Bytes returns the BitArray as bytes.
func (ba BitArray) Bytes() []byte {
	return ba.raw
}

// Len returns the BitArray length.
func (ba BitArray) Len() int64 {
	return ba.size
}

// Bytes returns the BitArray as bytes.
// The remaining bits of the last byte are set to zero.
func (ba BitArray) String() string {
	var s strings.Builder
	s.WriteByte('[')
	for i := int64(0); i < ba.size; i++ {
		if i > 0 && i%8 == 0 {
			s.WriteByte(' ')
		}
		if ba.Test(i) {
			s.WriteByte('1')
		} else {
			s.WriteByte('0')
		}
	}
	for i := uint(0); i < ba.avail(); i++ {
		s.WriteByte('-')
	}
	s.WriteByte(']')
	return s.String()
}

// Test returns true if bit at (zero based) offset i is 1, false otherwise.
func (ba BitArray) Test(i int64) bool {
	if i > ba.size-1 {
		return false
	}
	idx := i / 8
	offset := i % 8
	// if idx >= int64(len(ba.raw)) {
	// 	return false
	// }
	mask := 1 << uint(7-offset)
	return (ba.raw[idx] & byte(mask)) != 0
}

// Set a single bit to 1 at position n.
func (ba *BitArray) Set(n int64) {
	if n > ba.size {
		return
	}
	idx := n / 8
	offset := n % 8
	mask := 1 << uint(7-offset)
	ba.raw[idx] |= byte(mask)
}

// Unset a single bit to 0 at position n.
func (ba *BitArray) Unset(n int64) {
	if n > ba.size {
		return
	}
	idx := n / 8
	offset := n % 8
	mask := 1 << uint(7-offset)
	ba.raw[idx] &^= byte(mask)
}

// Pad array with n zeros.
func (ba *BitArray) Pad(n uint) int {
	c := 0
	for i := uint(0); i < n; i++ {
		c += ba.Add(0)
	}
	return c
}

// AddBit adds a single bit to the array.
func (ba *BitArray) AddBit(u uint) int {
	ba.grow()
	if u == 0 {
		ba.Unset(ba.size)
	} else {
		ba.Set(ba.size)
	}
	ba.size++
	return 1
}

// Add an uint to the array with leading zeros removed,
// returns the number of bits added.
func (ba *BitArray) Add(u uint) int {
	if u == 0 {
		ba.AddBit(0)
		return 1
	}
	used := bits.Len(u)
	for i := used - 1; i >= 0; i-- {
		set := uint8(u>>uint(i)) & 0x01
		ba.AddBit(uint(set))
	}
	return used
}

// AddN adds a uint with a fixed width of n, left padded to width with zeros,
// returns the number of bits added.
func (ba *BitArray) AddN(u uint, width int) int {
	n := bits.Len(u)
	if n > width {
		// Truncate
		return ba.Add(u >> (n - width))
	}
	c := ba.Pad(uint(width - n))
	if n != 0 {
		c += ba.Add(u)
	}
	return c
}

// Pack stuff together into existing array.
func (ba *BitArray) Pack(fields ...interface{}) error {
	for _, f := range fields {
		switch c := f.(type) {
		case uint:
			ba.Add(c)
		case uint8:
			ba.Add(uint(c))
		case uint16:
			ba.Add(uint(c))
		case uint32:
			ba.Add(uint(c))
		case uint64:
			// TODO check uintSize?
			ba.Add(uint(c))
		case []uint:
			for _, i := range c {
				ba.Add(i)
			}
		case []uint8:
			for _, i := range c {
				ba.Add(uint(i))
			}
		case []uint16:
			for _, i := range c {
				ba.Add(uint(i))
			}
		case []uint32:
			for _, i := range c {
				ba.Add(uint(i))
			}
		case []uint64:
			// TODO check uintSize?
			for _, i := range c {
				ba.Add(uint(i))
			}
		case int:
			ba.Add(uint(c))
		case int8:
			ba.Add(uint(c))
		case int16:
			ba.Add(uint(c))
		case int32:
			ba.Add(uint(c))
		case int64:
			ba.Add(uint(c))
		case []int:
			for _, i := range c {
				ba.Add(uint(i))
			}
		case []int8:
			for _, i := range c {
				ba.Add(uint(i))
			}
		case []int16:
			for _, i := range c {
				ba.Add(uint(i))
			}
		case []int32:
			for _, i := range c {
				ba.Add(uint(i))
			}
		case []int64:
			for _, i := range c {
				ba.Add(uint(i))
			}
		case []interface{}:
			return ba.Pack(c...)
		default:
			return fmt.Errorf("unable to pack %T", c)
		}
	}
	return nil
}

// Pack stuff together into a BitArray.
func Pack(fields ...interface{}) (*BitArray, error) {
	out := new(BitArray)
	err := out.Pack(fields...)
	return out, err
}

// Append packs another BitArray on the end.
func (ba *BitArray) Append(others ...BitArray) {
	for _, o := range others {
		for i := int64(0); i < o.Len(); i++ {
			if o.Test(i) {
				ba.AddBit(1)
			} else {
				ba.AddBit(0)
			}
		}
	}
}

// Slice reads a range from the BitArray.
func (ba *BitArray) Slice(start, length int64) (*BitArray, error) {
	out := new(BitArray)
	for i := start; i < (start + length); i++ {
		if ba.Test(i) {
			out.AddBit(1)
		} else {
			out.AddBit(0)
		}
	}
	return out, nil
}

// ReadUint reads a uint from the BitArray.
func (ba *BitArray) ReadUint(start, length int64) (uint, error) {
	b, err := ba.Slice(start, length)
	if err != nil {
		return 0, err
	}
	out := new(big.Int).SetBytes(b.Bytes())
	return uint(out.Rsh(out, b.avail()).Uint64()), nil
}

// grow the underlying storage until we have available bits.
func (ba *BitArray) grow() {
	if ba.avail() <= 0 {
		ba.raw = append(ba.raw, byte(0))
	}
}

// The number of bits available in the underlying storage.
func (ba BitArray) avail() uint {
	return uint(int64(len(ba.raw)*8) - ba.size)
}

// ShiftL shifts all bits to the left and returns those
// shifted off. s cannot be larger than 8.
// TODO shift more than 8!
func (ba *BitArray) ShiftL(s uint) (r byte) {
	if s > 8 {
		return
	}
	n := len(ba.raw)
	if n == 0 {
		return
	}

	_s := 8 - s
	b1 := ba.raw[n-1]
	r = b1 >> _s
	for i := 0; i < n-1; i++ {
		b := b1
		b1 = ba.raw[i+1]
		ba.raw[i] = b<<s | b1>>_s
	}
	ba.raw[n-1] = b1 << s
	return r
}

// ShiftR shifts all bits to the right and returns those
// shifted off. s cannot be larger than 8.
// TODO shift more than 8!
func (ba *BitArray) ShiftR(s uint) (r byte) {
	if s > 8 {
		return
	}
	n := len(ba.raw)
	if n == 0 {
		return
	}

	_s := 8 - s
	b1 := ba.raw[0]
	r = b1 << _s
	for i := n - 1; i > 0; i-- {
		b := b1
		b1 = ba.raw[i-1]
		ba.raw[i] = b>>s | b1<<_s
	}
	ba.raw[0] = b1 >> s
	return r
}
