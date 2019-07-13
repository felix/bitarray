package bitarray

import (
	"fmt"
	"math/bits"
)

// BitArray def
type BitArray struct {
	raw   []byte
	avail uint
}

// New create an empty BitArray.
func New() *BitArray {
	return &BitArray{
		raw:   make([]byte, 1),
		avail: 8,
	}
}

// Bytes returns the BitArray as bytes.
func (ba BitArray) Bytes() []byte {
	return ba.raw
}

// Bytes returns the BitArray as bytes.
func (ba BitArray) String() string {
	return fmt.Sprintf("%08b", ba.raw)
}

// Add8 adds a uint8 to the BitArray.
func (ba *BitArray) Add8(u uint8) {
	if u == 0 {
		ba.avail--
		return
	}
	ba.add(u)
}

func (ba *BitArray) add(u uint8) {
	var mask uint8
	if ba.avail == 0 {
		ba.raw = append(ba.raw, byte(0))
		ba.avail += 8
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

// Add16 adds a uint8 to the BitArray.
func (ba *BitArray) Add16(u uint16) {
	if u == 0 {
		ba.avail--
		return
	}
	ba.add(uint8(u >> 8))
	ba.add(uint8(u))
}

// Add32 adds a uint32 to the BitArray.
func (ba *BitArray) Add32(u uint32) {
	if u == 0 {
		ba.avail--
		return
	}
	ba.add(uint8(u >> 24))
	ba.add(uint8(u >> 16))
	ba.add(uint8(u >> 8))
	ba.add(uint8(u))
}

// Add64 adds a uint64 to the BitArray.
func (ba *BitArray) Add64(u uint64) {
	if u == 0 {
		ba.avail--
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
	out := New()
	err := out.Pack(fields...)
	return out, err
}

func abs(i int) uint {
	if i < 0 {
		return uint(-i)
	}
	return uint(i)
}
