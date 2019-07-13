package bitarray

import (
	"fmt"
	"testing"
)

func TestPack(t *testing.T) {
	tests := []struct {
		in       interface{}
		expected string
	}{
		// Types
		{uint(108), "[11011000]"},
		{uint8(108), "[11011000]"},
		{uint16(108), "[11011000]"},
		{uint32(108), "[11011000]"},
		{uint64(108), "[11011000]"},
		{int(108), "[11011000]"},
		{int8(0xf), "[11110000]"},
		{int16(108), "[11011000]"},
		{int32(108), "[11011000]"},
		{int64(108), "[11011000]"},
		{[]byte{0x6c}, "[11011000]"},
		{[]uint{108}, "[11011000]"},
		{[]int{108}, "[11011000]"},
		{[]int8{108}, "[11011000]"},
		{[]int16{108}, "[11011000]"},
		{[]int32{108}, "[11011000]"},
		{[]int64{108}, "[11011000]"},
		// Zero
		{uint(0), "[00000000]"},
		{uint8(0), "[00000000]"},
		{uint16(0), "[00000000]"},
		{uint32(0), "[00000000]"},
		{uint64(0), "[00000000]"},
		{int(0), "[00000000]"},
		{int8(0), "[00000000]"},
		{int16(0), "[00000000]"},
		{int32(0), "[00000000]"},
		{int64(0), "[00000000]"},
		{[]byte{0}, "[00000000]"},
		{[]uint{0}, "[00000000]"},
		{[]int{0}, "[00000000]"},
		{[]int8{0}, "[00000000]"},
		{[]int16{0}, "[00000000]"},
		{[]int32{0}, "[00000000]"},
		{[]int64{0}, "[00000000]"},
		// Adding a zero
		{[]int{1, 0, 23}, "[10101110]"},
		{[]uint16{1, 0, 23}, "[10101110]"},
		{[]uint32{1, 0, 23}, "[10101110]"},
		{[]uint64{1, 0, 23}, "[10101110]"},
		// Cases
		{[]uint8{0xff, 0xff}, "[11111111 11111111]"},
		{[]uint8{0xff, 0xf0}, "[11111111 11110000]"},
		{[]uint8{0xf0, 0xf0}, "[11110000 11110000]"},
		{[]uint8{0xf0, 0xf0, 1}, "[11110000 11110000 10000000]"},
		{[]int{1, 128, 23}, "[11000000 01011100]"},
		{[]int{1, 129, 23}, "[11000000 11011100]"},
		{[]interface{}{uint8(0xff), 1, 2, 1, 4, 1, 1}, "[11111111 11011001 10000000]"},
	}

	for _, tt := range tests {
		ba, err := Pack(tt.in)
		if err != nil {
			t.Fatalf("%v => failed %q", tt.in, err)
		}
		actual := fmt.Sprintf("%08b", ba.Bytes())
		if actual != tt.expected {
			t.Errorf("%v => expected %s got %s", tt.in, tt.expected, actual)
		}
		if actual != ba.String() {
			t.Errorf("%v => expected %q got %q", tt.in, actual, ba.String())
		}
	}
}

func TestSlice(t *testing.T) {
	tests := []struct {
		ba       *BitArray
		s, l     uint64 // start and length
		expected string
	}{
		{&BitArray{raw: []byte{0xff}, avail: 0}, 0, 8, "[11111111]"},
		{&BitArray{raw: []byte{0xff}, avail: 0}, 0, 1, "[10000000]"},
		{&BitArray{raw: []byte{0xfe}, avail: 0}, 0, 8, "[11111110]"},
		{&BitArray{raw: []byte{0x03}, avail: 0}, 7, 1, "[10000000]"},
		{&BitArray{raw: []byte{0xd0}, avail: 4}, 0, 4, "[11010000]"},
		// Multiple bytes
		{&BitArray{raw: []byte{0xd0, 0xff}, avail: 0}, 0, 9, "[11010000 10000000]"},
		{&BitArray{raw: []byte{0x0f, 0xf0}, avail: 0}, 4, 8, "[11111111]"},
	}

	for _, tt := range tests {
		a, err := tt.ba.Slice(tt.s, tt.l)
		if err != nil {
			t.Errorf("failed with %q", err)
		}
		actual := fmt.Sprintf("%08b", a.Bytes())
		if actual != tt.expected {
			t.Errorf("%x %d,%d => expected %q got %q", tt.ba.raw, tt.s, tt.l, tt.expected, actual)
		}
	}
}

func TestReadBig(t *testing.T) {
	tests := []struct {
		ba       *BitArray
		s, l     uint64 // start and length
		expected uint64
	}{
		{&BitArray{raw: []byte{0x02}, avail: 0}, 0, 8, 0x02},
		{&BitArray{raw: []byte{0xff}, avail: 0}, 0, 8, 0xff},
		{&BitArray{raw: []byte{0xff}, avail: 0}, 0, 1, 0x01},
		{&BitArray{raw: []byte{0xfe}, avail: 0}, 0, 8, 0xfe},
		{&BitArray{raw: []byte{0x03}, avail: 0}, 7, 1, 0x01},
		{&BitArray{raw: []byte{0xd0}, avail: 4}, 0, 4, 0x0d},
		// Multiple bytes
		{&BitArray{raw: []byte{0xd0, 0xff}, avail: 0}, 0, 9, 0x1a1},
		{&BitArray{raw: []byte{0x0f, 0xf0}, avail: 0}, 4, 8, 0xff},
	}

	for _, tt := range tests {
		b, err := tt.ba.ReadBig(tt.s, tt.l)
		if err != nil {
			t.Errorf("failed with %q", err)
		}
		actual := b.Uint64()
		if actual != tt.expected {
			t.Errorf("%x %d,%d => expected %d got %d", tt.ba.raw, tt.s, tt.l, tt.expected, actual)
		}
	}
}
