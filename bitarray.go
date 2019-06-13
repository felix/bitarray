package bitarray

import (
	//"fmt"
	//"math"
	"math/bits"
	"sync"
)

// BitArray definition.
type BitArray struct {
	raw  []byte
	bits int
	lock sync.RWMutex
}

// Reset the bit array.
func (ba *BitArray) Reset() {
	ba.raw = nil
	ba.bits = 0
}

// Add64 integers onto the bit array.
func (ba *BitArray) Add64(in uint64) {
	ba.lock.Lock()
	defer ba.lock.Unlock()

	spare := uint64((len(ba.raw) * 8) - ba.bits)
	toAdd := uint64(bits.Len64(in))

	// fmt.Printf("  in=%064b\n", in)
	// fmt.Printf("start spare=%d toadd=%d\n", spare, toAdd)

	// Deal with spare first
	if spare > 0 {
		var next uint8
		if toAdd > spare {
			//fmt.Printf("more spare=%d toadd=%d\n", spare, toAdd)
			// Shift in to right, align with spare bits
			next = uint8(in >> (toAdd - spare))
			ba.bits += int(spare)
			toAdd -= spare
		} else {
			// Should fit in spare
			//fmt.Printf("fit spare=%d toadd=%d\n", spare, toAdd)
			next = uint8((in << (spare - toAdd)) & 0xff)
			//fmt.Printf("next=%08b\n", next)
			ba.bits += int(toAdd)
			toAdd = 0
		}

		// Work on the last byte
		l := len(ba.raw) - 1
		next = uint8(ba.raw[l]) | next
		ba.raw[l] = byte(next)
		spare = uint64((len(ba.raw) * 8) - ba.bits)
		//fmt.Printf("post-spare bits=%d spare=%d toadd=%d\n", ba.bits, spare, toAdd)
	}

	// Chunk per byte
	for i := uint64(0); i < toAdd; i += 8 {
		//fmt.Printf("looping %d bits=%d spare=%d toadd=%d\n", i, ba.bits, spare, toAdd)
		next := (in << (64 - toAdd)) & 0xff00000000000000
		//fmt.Printf("next=%064b\n", next)
		next >>= 56
		//fmt.Printf("next=%08b\n", next)
		ba.raw = append(ba.raw, byte(next))
	}
	ba.bits += int(toAdd)
	//spare = uint64((len(ba.raw) * 8) - ba.bits)
	//fmt.Printf("spare=%d bits=%d\n", spare, ba.bits)
}

// Add32 integers onto the bit array.
func (ba *BitArray) Add32(in uint32) { ba.Add64(uint64(in)) }

// Add16 integers onto the bit array.
func (ba *BitArray) Add16(in uint16) { ba.Add64(uint64(in)) }

// Add8 integers onto the bit array.
func (ba *BitArray) Add8(in uint8) { ba.Add64(uint64(in)) }

// Read bits from the bit array.
/*
func (ba *BitArray) Read(start, length uint64) []byte {
	ba.lock.RLock()
	defer ba.lock.RUnlock()
	startB := int(math.Ceil(float64(start) / float64(8)))
	endB := int(math.Ceil(float64(start+length) / float64(8)))
	bs := ba.raw[startB:endB]

	// Number of bits in the first byte
	bitsInFirst := start % 8
	// Number of bits in the last byte
	bitsInLast := (length - bitsInFirst) % 8

	out := new(BitArray)
	for i, b := range bs {
		if i == 0 {
			// First byte
			mask := uint8(math.Pow(2, float64(bitsInFirst)))
			out.Add8(uint8(b) & mask)
		} else if i == len(bs)-1 {
			// Last byte
			mask := uint8(math.Pow(2, float64(bitsInLast)))
			last := uint8(b) & (mask << uint(bits.LeadingZeros8(mask)))
			last >>= uint8((len(ba.raw) * 8) - ba.bits)
			out.Add8(last)
		} else {
			// Whole byte
			out.Add8(uint8(b))
		}
	}
	return out.Bytes()
}

// Read8 bits from the bit array as a uint8.
func (ba *BitArray) Read8(start, length uint64) (uint8, error) {
	if length > 8 {
		return 0, fmt.Errorf("truncated result")
	}
	bs := ba.Read(start, length)
	if len(bs) > 1 {
		return uint8(bs[0]), fmt.Errorf("invalid byte length")
	}
	return uint8(bs[0]), nil
}
*/

// Bytes returns the compacted bit array.
func (ba *BitArray) Bytes() []byte {
	ba.lock.RLock()
	defer ba.lock.RUnlock()

	return ba.raw
}
