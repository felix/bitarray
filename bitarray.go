package bitarray

import (
	"fmt"
	"math/big"
	"math/bits"
)

// BitArray def
type BitArray struct {
	raw   []byte
	avail uint8
}

// New create an empty BitArray.
func New(b []byte, padding uint8) *BitArray {
	out := &BitArray{}
	out.raw = b
	//fmt.Printf("new raw=%08b\n", out.raw)
	out.avail = padding
	return out
}

// Bytes returns the BitArray as bytes.
func (ba BitArray) Bytes() []byte {
	return ba.raw
}

// Len returns the BitArray length.
func (ba BitArray) Len() uint64 {
	return uint64(len(ba.raw))*8 - uint64(ba.avail)
}

// Bytes returns the BitArray as bytes.
func (ba BitArray) String() string {
	return fmt.Sprintf("%08b", ba.raw)
}

// Test returns true/false on bit at offset i.
func (ba BitArray) Test(i uint64) bool {
	// if i >= ba.Len() {
	// 	return false
	// }
	idx := i / 8
	offset := i % 8
	if idx >= uint64(len(ba.raw)) {
		return false
	}
	mask := 1 << (7 - offset)
	return (ba.raw[idx] & byte(mask)) != 0
}

// Set a single bit at position n.
func (ba *BitArray) Set(n uint64) {
	if n >= ba.Len() {
		return
	}
	idx := n / 8
	offset := n % 8
	mask := 1 << (7 - offset)
	ba.raw[idx] |= byte(mask)
}

// Unset a single bit at position n.
func (ba *BitArray) Unset(n uint64) {
	if n >= ba.Len() {
		return
	}
	idx := n / 8
	offset := n % 8
	mask := 1 << (7 - offset)
	ba.raw[idx] &^= byte(mask)
}

// Pad adds a zero padding.
func (ba *BitArray) Pad(n uint64) {
	for i := uint64(0); i < n; i++ {
		ba.Add8(uint8(0))
	}
}

// AddBit adds a single bit to the array.
func (ba *BitArray) AddBit(u uint8) {
	ba.grow()
	if u == 0 {
		ba.Unset(ba.Len())
	} else {
		ba.Set(ba.Len())
	}
	ba.avail--
}

// Add8 adds a uint8 to the BitArray.
func (ba *BitArray) Add8(u uint8) {
	if u == 0 {
		ba.AddBit(0)
		return
	}
	ba.add(u)
}

// Add8N adds a uint8 with a fixed width of n.
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

// Slice reads a range from the BitArray.
func (ba *BitArray) Slice(start, length uint64) (*BitArray, error) {
	out := new(BitArray)
	for i := start; i < (start + length); i++ {
		if ba.Test(i) {
			out.Add8(1)
		} else {
			out.Add8(0)
		}
	}
	return out, nil
}

// ReadBig reads a big.Int from the BitArray.
func (ba *BitArray) ReadBig(start, length uint64) (*big.Int, error) {
	fmt.Printf("about to read from=%08b start=%d l=%d\n", ba.raw, start, length)
	b, err := ba.Slice(start, length)
	if err != nil {
		return nil, err
	}
	fmt.Printf("read from=%x start=%d l=%d %08b avail=%d\n", ba.raw, start, length, b.raw, b.avail)
	out := new(big.Int).SetBytes(b.Bytes())
	shifted := out.Rsh(out, uint(b.avail))
	fmt.Printf("shifted=%08b\n", shifted.Bytes())
	//return out.Rsh(out, b.avail), nil
	return shifted, nil
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
	ba.avail = abs(shift)
}

func (ba *BitArray) grow() {
	if ba.avail <= 0 {
		ba.raw = append(ba.raw, byte(0))
		ba.avail += 8
	}
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

func abs(i int8) uint8 {
	if i < 0 {
		return uint8(-i)
	}
	return uint8(i)
}

func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}
