package bitarray

import (
	"fmt"
	"math/big"
	"math/bits"
)

// BitArray def
type BitArray struct {
	raw   []byte
	avail uint
}

// New create an empty BitArray.
func New(b []byte, padding uint) *BitArray {
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
	idx := i / 8
	offset := i % 8
	if idx >= uint64(len(ba.raw)) {
		return false
	}
	mask := 1 << (7 - offset)
	return (ba.raw[idx] & byte(mask)) != 0
}

// Pad adds a zero padding.
func (ba *BitArray) Pad(n uint64) {
	for i := uint64(0); i < n; i++ {
		ba.Add8(uint8(0))
	}
}

// Add8 adds a uint8 to the BitArray.
func (ba *BitArray) Add8(u uint8) {
	ba.add(u, u == 0)
}

// Add16 adds a uint8 to the BitArray.
func (ba *BitArray) Add16(u uint16) {
	if u == 0 {
		ba.add(0, true)
		return
	}
	ba.add(uint8(u>>8), false)
	ba.add(uint8(u), false)
}

// Add32 adds a uint32 to the BitArray.
func (ba *BitArray) Add32(u uint32) {
	if u == 0 {
		ba.add(0, true)
		return
	}
	ba.add(uint8(u>>24), false)
	ba.add(uint8(u>>16), false)
	ba.add(uint8(u>>8), false)
	ba.add(uint8(u), false)
}

// Add64 adds a uint64 to the BitArray.
func (ba *BitArray) Add64(u uint64) {
	if u == 0 {
		ba.add(0, true)
		return
	}
	ba.add(uint8(u>>56), false)
	ba.add(uint8(u>>48), false)
	ba.add(uint8(u>>40), false)
	ba.add(uint8(u>>32), false)
	ba.add(uint8(u>>24), false)
	ba.add(uint8(u>>16), false)
	ba.add(uint8(u>>8), false)
	ba.add(uint8(u), false)
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
	shifted := out.Rsh(out, b.avail)
	fmt.Printf("shifted=%08b\n", shifted.Bytes())
	//return out.Rsh(out, b.avail), nil
	return shifted, nil
}

func (ba *BitArray) add(u uint8, zero bool) {
	var mask uint8
	if ba.avail == 0 {
		ba.raw = append(ba.raw, byte(0))
		ba.avail += 8
	}
	if zero {
		ba.avail--
		return
	}
	n := uint(bits.Len8(u))
	shift := int(ba.avail - n)
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

func abs(i int) uint {
	if i < 0 {
		return uint(-i)
	}
	return uint(i)
}

func min(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}
