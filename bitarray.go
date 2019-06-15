package bitarray

import (
	"fmt"
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

// Read64 bits from the bit array.
func (ba *BitArray) Read64(start, length uint64) uint64 {
	ba.lock.RLock()
	defer ba.lock.RUnlock()
	//bs := ba.raw[startB:endB]
	fmt.Printf("raw=%08b start=%d length=%d\n", ba.raw, start, length)

	// First relevant byte
	firstByte := int(start / 8)
	lastByte := int((start + length - 1) / 8)
	fmt.Printf("firstByte=%d lastByte=%d\n", firstByte, lastByte)
	// Number of bits in the last byte
	bitsInLast := uint(start+length) - uint(lastByte*8)

	fmt.Printf("bitsinlast=%d\n", bitsInLast)

	// A mask for each byte
	masks := make([]uint8, lastByte+1-firstByte)
	for i := start; i < (start + length); i++ {
		// The bit we want is in byteIdx byte
		byteIdx := i / 8
		//b := ba.raw[byteIdx]

		// bitIdx is the bit offset in this byte
		bitIdx := i - (byteIdx * 8)
		masks[byteIdx] |= (1 << (7 - bitIdx))
		fmt.Printf("byteIdx=%d bitIdx=%d mask=%08b\n", byteIdx, bitIdx, masks[byteIdx])
	}
	var out uint64
	for i, m := range masks {
		b := uint8(ba.raw[i]) & m
		out |= uint64(b << uint(i))
	}
	fmt.Printf("out=%08b\n", out)
	out >>= (8 - bitsInLast)
	fmt.Printf("out=%08b\n", out)
	return out
	/*
		out := make([]byte, len(masks))
		for i, m := range masks {
			out[i] = uint8(ba.raw[i]) & m
		}
		return out
	*/
}

// Read8 bits from the bit array as a uint8.
func (ba *BitArray) Read8(start, length uint64) (uint8, error) {
	if length > 8 {
		return 0, fmt.Errorf("truncated result")
	}
	bs := ba.Read64(start, length)
	return uint8(bs), nil
}

// Bytes returns the compacted bit array.
func (ba *BitArray) Bytes() []byte {
	ba.lock.RLock()
	defer ba.lock.RUnlock()

	return ba.raw
}
