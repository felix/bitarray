package bitarray

import (
	"fmt"
	"math/big"
	"math/bits"
)

// BitArray holds an array of bits.
type BitArray struct {
	raw   []byte
	avail int
	size  int
}

// New creates a BitArray with the given []byte and count bits.
// Count is from position 0 of b.
func New(b []byte, count int) *BitArray {
	out := &BitArray{}
	out.raw = b
	out.size = count
	out.avail = len(b)*8 - count
	return out
}

// BitWriter is the interface that wraps the WriteBit method.
type BitWriter interface {
	WriteBits(uint64, uint64) error
}

// Bytes returns the BitArray as bytes.
func (ba BitArray) Bytes() []byte {
	return ba.raw
}

// Len returns the BitArray length.
func (ba BitArray) Len() int {
	return ba.size
}

// Bytes returns the BitArray as bytes.
// The remaining bits of the last byte are set to zero.
func (ba BitArray) String() string {
	return fmt.Sprintf("%08b", ba.raw)
}

// Test returns true if bit at offset i is 1, false otherwise.
func (ba BitArray) Test(i int) bool {
	if i > ba.size {
		return false
	}
	idx := i / 8
	offset := i % 8
	if idx >= len(ba.raw) {
		return false
	}
	mask := 1 << uint(7-offset)
	return (ba.raw[idx] & byte(mask)) != 0
}

// Set a single bit to 1 at position n.
func (ba *BitArray) Set(n int) {
	if n > ba.size {
		return
	}
	idx := n / 8
	offset := n % 8
	mask := 1 << uint(7-offset)
	ba.raw[idx] |= byte(mask)
}

// Unset a single bit to 0 at position n.
func (ba *BitArray) Unset(n int) {
	if n > ba.size {
		return
	}
	idx := n / 8
	offset := n % 8
	mask := 1 << uint(7-offset)
	ba.raw[idx] &^= byte(mask)
}

// Pad adds n zeros as padding.
func (ba *BitArray) Pad(n uint64) {
	for i := uint64(0); i < n; i++ {
		ba.Add8(uint8(0))
	}
}

// AddBit adds a single bit to the array.
func (ba *BitArray) AddBit(u uint8) {
	ba.grow()
	if u == 0 {
		ba.Unset(ba.size)
	} else {
		ba.Set(ba.size)
	}
	ba.size++
	ba.avail--
}

// Add8 adds a uint8 to the BitArray with leading zeros removed.
func (ba *BitArray) Add8(u uint8) {
	if u == 0 {
		ba.AddBit(0)
		return
	}
	ba.add(u)
}

// Add8N adds a uint8 with a fixed width of n, left padded with zeros.
func (ba *BitArray) Add8N(u, s uint8) {
	n := uint8(bits.Len8(u))
	for i := uint8(0); i < (s - n); i++ {
		ba.AddBit(0)
	}
	if n != 0 {
		ba.Add8(u)
	}
}

// Add16 adds a uint8 to the BitArray.
func (ba *BitArray) Add16(u uint16) {
	if u == 0 {
		ba.AddBit(0)
		return
	}
	ba.add(uint8(u >> 8))
	ba.add(uint8(u))
}

// Add16N adds a uint16 with a fixed width of n.
func (ba *BitArray) Add16N(u, s uint16) {
	n := uint16(bits.Len16(u))
	for i := uint16(0); i < (s - n); i++ {
		ba.AddBit(0)
	}
	if n != 0 {
		ba.Add16(u)
	}
}

// Add32 adds a uint32 to the BitArray.
func (ba *BitArray) Add32(u uint32) {
	if u == 0 {
		ba.AddBit(0)
		return
	}
	ba.add(uint8(u >> 24))
	ba.add(uint8(u >> 16))
	ba.add(uint8(u >> 8))
	ba.add(uint8(u))
}

// Add32N adds a uint32 with a fixed width of n.
func (ba *BitArray) Add32N(u, s uint32) {
	n := uint32(bits.Len32(u))
	for i := uint32(0); i < (s - n); i++ {
		ba.AddBit(0)
	}
	if n != 0 {
		ba.Add32(u)
	}
}

// Add64 adds a uint64 to the BitArray.
func (ba *BitArray) Add64(u uint64) {
	if u == 0 {
		ba.AddBit(0)
		return
	}
	ba.add(uint8(u >> 56))
	ba.add(uint8(u >> 48))
	ba.add(uint8(u >> 40))
	ba.add(uint8(u >> 32))
	ba.add(uint8(u >> 24))
	ba.add(uint8(u >> 16))
	ba.add(uint8(u >> 8))
	ba.add(uint8(u))
}

// Add64N adds a uint64 with a fixed width of n.
func (ba *BitArray) Add64N(u, s uint64) {
	n := uint64(bits.Len64(u))
	for i := uint64(0); i < (s - n); i++ {
		ba.AddBit(0)
	}
	if n != 0 {
		ba.Add64(u)
	}
}

// Pack stuff together into existing BitArray.
func (ba *BitArray) Pack(fields ...interface{}) error {
	for _, f := range fields {
		switch c := f.(type) {
		case uint:
			ba.Add32(uint32(c))
		case uint8:
			ba.Add8(c)
		case uint16:
			ba.Add16(c)
		case uint32:
			ba.Add32(c)
		case uint64:
			ba.Add64(c)
		case []uint:
			for _, i := range c {
				ba.Add32(uint32(i))
			}
		case []uint8:
			for _, i := range c {
				ba.Add8(i)
			}
		case []uint16:
			for _, i := range c {
				ba.Add16(i)
			}
		case []uint32:
			for _, i := range c {
				ba.Add32(i)
			}
		case []uint64:
			for _, i := range c {
				ba.Add64(i)
			}
		case int:
			ba.Add32(uint32(c))
		case int8:
			ba.Add8(uint8(c))
		case int16:
			ba.Add16(uint16(c))
		case int32:
			ba.Add32(uint32(c))
		case int64:
			ba.Add64(uint64(c))
		case []int:
			for _, i := range c {
				ba.Add32(uint32(i))
			}
		case []int8:
			for _, i := range c {
				ba.Add8(uint8(i))
			}
		case []int16:
			for _, i := range c {
				ba.Add16(uint16(i))
			}
		case []int32:
			for _, i := range c {
				ba.Add32(uint32(i))
			}
		case []int64:
			for _, i := range c {
				ba.Add64(uint64(i))
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
		for i := 0; i < o.Len(); i++ {
			if o.Test(i) {
				ba.AddBit(1)
			} else {
				ba.AddBit(0)
			}
		}
	}
}

// Slice reads a range from the BitArray.
func (ba *BitArray) Slice(start, length int) (*BitArray, error) {
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

// ReadBig reads a big.Int from the BitArray.
func (ba *BitArray) ReadBig(start, length int) (*big.Int, error) {
	b, err := ba.Slice(start, length)
	if err != nil {
		return nil, err
	}
	out := new(big.Int).SetBytes(b.Bytes())
	return out.Rsh(out, uint(b.avail)), nil
}

func (ba *BitArray) add(u uint8) {
	n := uint8(bits.Len8(u))
	ba.addN(u, n)
}

func (ba *BitArray) addN(u, n uint8) {
	ba.grow()
	var mask uint8
	shift := int8(uint8(ba.avail) - n)
	if shift < 0 {
		// It doesn't fit
		mask = u >> abs(shift)
		ba.raw[len(ba.raw)-1] |= byte(mask)
		ba.raw = append(ba.raw, byte(0))
		shift = 8 + shift
	}
	mask = u << abs(shift)
	ba.raw[len(ba.raw)-1] |= byte(mask)
	ba.avail = int(abs(shift))
	ba.size += int(n)
	ba.norm()
}

func (ba *BitArray) grow() {
	if ba.avail <= 0 {
		ba.raw = append(ba.raw, byte(0))
		ba.avail += 8
	}
}

// Remove extraneous bits.
func (ba *BitArray) norm() {
	n := len(ba.raw)
	if ba.avail > 8 {
		ba.raw = ba.raw[:n]
	}
	ba.avail = n*8 - ba.size
}

// ShiftL shifts all bits to the left and returns those
// shifted off. s cannot be larger than 8.
func (ba *BitArray) ShiftL(s uint8) (r byte) {
	if s > 8 {
		return
	}
	if n := len(ba.raw); n > 0 {
		_s := 8 - s
		b1 := ba.raw[n-1]
		r = b1 >> _s
		for i := 0; i < n-1; i++ {
			b := b1
			b1 = ba.raw[i+1]
			ba.raw[i] = b<<s | b1>>_s
		}
		ba.raw[n-1] = b1 << s
	}
	return
}

// ShiftR shifts all bits to the right and returns those
// shifted off. s cannot be larger than 8.
func (ba *BitArray) ShiftR(s uint8) (r byte) {
	if s > 8 {
		return
	}
	if n := len(ba.raw); n > 0 {
		_s := 8 - s
		b1 := ba.raw[0]
		r = b1 << _s
		for i := n - 1; i > 0; i-- {
			b := b1
			b1 = ba.raw[i-1]
			ba.raw[i] = b>>s | b1<<_s
		}
		ba.raw[0] = b1 >> s
	}
	return
}

func abs(i int8) uint {
	if i < 0 {
		return uint(-i)
	}
	return uint(i)
}
