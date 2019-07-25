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
		s, l     int // start and length
		expected string
		avail    int
	}{
		{New([]byte{0xff}, 8), 0, 8, "[11111111]", 0},
		{New([]byte{0xff}, 8), 0, 1, "[10000000]", 7},
		{New([]byte{0xfe}, 8), 0, 8, "[11111110]", 0},
		{New([]byte{0x03}, 8), 7, 1, "[10000000]", 7},
		{New([]byte{0xd0}, 4), 0, 4, "[11010000]", 4},
		// Multiple bytes
		{New([]byte{0xd0, 0xff}, 16), 0, 9, "[11010000 10000000]", 7},
		{New([]byte{0x0f, 0xf0}, 16), 4, 8, "[11111111]", 0},
		// Cases
		// 10010110 00101100 01001001 => 1000101
		{New([]byte{0x96, 0x2c, 0x49}, 24), 6, 7, "[10001010]", 1},
		// 10010110 00101100 01001001 01110010 00101011 10000000
		//                               ^^^^^ ^^
		{New([]byte{0x96, 0x2c, 0x49, 0x72, 0x2b, 0x80}, 48), 27, 7, "[10010000]", 1},
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
		if a.avail != tt.avail {
			t.Errorf("%x %d,%d => expected avail %d got %d", tt.ba.raw, tt.s, tt.l, tt.avail, a.avail)
		}
	}
}

func TestReadBig(t *testing.T) {
	tests := []struct {
		ba       *BitArray
		s, l     int // start and length
		expected uint64
	}{
		{New([]byte{0x02}, 8), 0, 8, 0x02},
		{New([]byte{0xff}, 8), 0, 8, 0xff},
		{New([]byte{0xff}, 8), 0, 1, 0x01},
		{New([]byte{0xfe}, 8), 0, 8, 0xfe},
		{New([]byte{0x03}, 8), 7, 1, 0x01},
		{New([]byte{0xd0}, 8), 0, 4, 0x0d},
		// Multiple bytes
		{New([]byte{0xd0, 0xff}, 16), 0, 9, 0x1a1},
		{New([]byte{0x0f, 0xf0}, 16), 4, 8, 0xff},
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

func TestTest(t *testing.T) {
	tests := []struct {
		ba       *BitArray
		in       int
		expected bool
	}{
		{New([]byte{0x01}, 8), 0, false},
		{New([]byte{0x01}, 8), 1, false},
		{New([]byte{0x01}, 8), 7, true},
		{New([]byte{0x00, 0x80}, 16), 8, true},
		{New([]byte{0x00, 0x02}, 16), 14, true},
	}

	for _, tt := range tests {
		actual := tt.ba.Test(tt.in)
		if actual != tt.expected {
			t.Errorf("%d => expected %t got %t", tt.in, tt.expected, actual)
		}
	}
}

func TestSet(t *testing.T) {
	tests := []struct {
		ba       *BitArray
		in       int
		expected string
	}{
		{New([]byte{0x00}, 1), 0, "[10000000]"},
		{New([]byte{0x80}, 8), 1, "[11000000]"},
	}

	for _, tt := range tests {
		tt.ba.Set(tt.in)
		actual := fmt.Sprintf("%08b", tt.ba.Bytes())
		if actual != tt.expected {
			t.Errorf("%v => expected %s got %s", tt.in, tt.expected, actual)
		}
	}
}

func TestAddBit(t *testing.T) {
	tests := []struct {
		ba       *BitArray
		in       uint8
		expected string
	}{
		{New([]byte{0xF0}, 4), 1, "[11111000]"},
		{New([]byte{0xFF}, 8), 1, "[11111111 10000000]"},
	}

	for _, tt := range tests {
		lBefore := tt.ba.Len()
		tt.ba.AddBit(tt.in)
		actual := fmt.Sprintf("%08b", tt.ba.Bytes())
		if actual != tt.expected {
			t.Errorf("%v => expected %s got %s", tt.in, tt.expected, actual)
		}
		expected := lBefore + 1
		if tt.ba.Len() != expected {
			t.Errorf("%v => expected %d got %d", tt.in, expected, tt.ba.Len())
		}
	}
}

func TestShiftL(t *testing.T) {
	tests := []struct {
		ba       *BitArray
		s        uint8
		expected string
	}{
		{New([]byte{0x01}, 7), 1, "[00000010]"},
		{New([]byte{0x01}, 7), 2, "[00000100]"},
		{New([]byte{0x01}, 7), 3, "[00001000]"},
		{New([]byte{0x01}, 7), 4, "[00010000]"},
		{New([]byte{0x01}, 7), 5, "[00100000]"},
		{New([]byte{0x01}, 7), 6, "[01000000]"},
		{New([]byte{0x01}, 7), 7, "[10000000]"},
		{New([]byte{0x01}, 7), 8, "[00000000]"},
		// Across a byte
		{New([]byte{0x01, 0x01}, 7), 5, "[00100000 00100000]"},
		{New([]byte{0x01, 0x01}, 7), 8, "[00000001 00000000]"},
	}

	for _, tt := range tests {
		tt.ba.ShiftL(tt.s)
		actual := fmt.Sprintf("%08b", tt.ba.Bytes())
		if actual != tt.expected {
			t.Errorf("%d => expected %s got %s", tt.s, tt.expected, actual)
		}
	}
}

func TestShiftR(t *testing.T) {
	tests := []struct {
		ba       *BitArray
		s        uint8
		expected string
	}{
		{New([]byte{0x80}, 8), 1, "[01000000]"},
		{New([]byte{0x80}, 8), 2, "[00100000]"},
		{New([]byte{0x80}, 8), 3, "[00010000]"},
		{New([]byte{0x80}, 8), 4, "[00001000]"},
		{New([]byte{0x80}, 8), 5, "[00000100]"},
		{New([]byte{0x80}, 8), 6, "[00000010]"},
		{New([]byte{0x80}, 8), 7, "[00000001]"},
		{New([]byte{0x80}, 8), 8, "[00000000]"},
		// Across a byte
		{New([]byte{0x80, 0x80}, 16), 5, "[00000100 00000100]"},
		{New([]byte{0x80, 0x80}, 16), 8, "[00000000 10000000]"},
	}

	for _, tt := range tests {
		tt.ba.ShiftR(tt.s)
		actual := fmt.Sprintf("%08b", tt.ba.Bytes())
		if actual != tt.expected {
			t.Errorf("%d => expected %s got %s", tt.s, tt.expected, actual)
		}
	}
}

func TestAdd8N(t *testing.T) {
	tests := []struct {
		ba       *BitArray
		in       uint8
		l        uint8
		expected string
	}{
		{New([]byte{0xF0}, 4), 1, 1, "[11111000]"},
		{New([]byte{0xF0}, 4), 1, 2, "[11110100]"},
		{New([]byte{0xF0}, 4), 1, 3, "[11110010]"},
		{New([]byte{0xF0}, 4), 1, 4, "[11110001]"},
		{New([]byte{0xF0}, 4), 1, 5, "[11110000 10000000]"},
	}

	for _, tt := range tests {
		lBefore := tt.ba.Len()
		tt.ba.Add8N(tt.in, tt.l)
		actual := fmt.Sprintf("%08b", tt.ba.Bytes())
		if actual != tt.expected {
			t.Errorf("%v => expected %s got %s", tt.in, tt.expected, actual)
		}
		expected := lBefore + int(tt.l)
		if tt.ba.Len() != expected {
			t.Errorf("%v => expected %d got %d", tt.in, expected, tt.ba.Len())
		}
	}
}

func TestAdd16N(t *testing.T) {
	tests := []struct {
		ba       *BitArray
		in       uint16
		l        uint16
		expected string
	}{
		{New([]byte{0xF0}, 4), 1, 1, "[11111000]"},
		{New([]byte{0xF0}, 4), 1, 2, "[11110100]"},
		{New([]byte{0xF0}, 4), 1, 3, "[11110010]"},
		{New([]byte{0xF0}, 4), 1, 4, "[11110001]"},
		{New([]byte{0xF0}, 4), 1, 5, "[11110000 10000000]"},
		{New([]byte{0xF0}, 4), 1, 6, "[11110000 01000000]"},
		{New([]byte{0xF0}, 4), 1, 7, "[11110000 00100000]"},
		{New([]byte{0xF0}, 4), 1, 8, "[11110000 00010000]"},
		{New([]byte{0xF0}, 4), 1, 9, "[11110000 00001000]"},
		{New([]byte{0xF0}, 4), 1, 10, "[11110000 00000100]"},
		{New([]byte{0xF0}, 4), 1, 11, "[11110000 00000010]"},
		{New([]byte{0xF0}, 4), 1, 12, "[11110000 00000001]"},
		{New([]byte{0xF0}, 4), 1, 13, "[11110000 00000000 10000000]"},
		// Add zero
		{New([]byte{0xF0}, 4), 0, 6, "[11110000 00000000]"},
	}

	for _, tt := range tests {
		lBefore := tt.ba.Len()
		tt.ba.Add16N(tt.in, tt.l)
		actual := fmt.Sprintf("%08b", tt.ba.Bytes())
		if actual != tt.expected {
			t.Errorf("%v => expected %s got %s", tt.in, tt.expected, actual)
		}
		expected := lBefore + int(tt.l)
		if tt.ba.Len() != expected {
			t.Errorf("%v => expected %d got %d", tt.in, expected, tt.ba.Len())
		}
	}
}

func TestAdd32N(t *testing.T) {
	tests := []struct {
		ba       *BitArray
		in       uint32
		l        uint32
		expected string
	}{
		{New([]byte{0xF0}, 4), 1, 1, "[11111000]"},
		{New([]byte{0xF0}, 4), 1, 2, "[11110100]"},
		{New([]byte{0xF0}, 4), 1, 3, "[11110010]"},
		{New([]byte{0xF0}, 4), 1, 4, "[11110001]"},
		{New([]byte{0xF0}, 4), 1, 5, "[11110000 10000000]"},
		{New([]byte{0xF0}, 4), 1, 6, "[11110000 01000000]"},
		{New([]byte{0xF0}, 4), 1, 7, "[11110000 00100000]"},
		{New([]byte{0xF0}, 4), 1, 8, "[11110000 00010000]"},
		{New([]byte{0xF0}, 4), 1, 9, "[11110000 00001000]"},
		{New([]byte{0xF0}, 4), 1, 10, "[11110000 00000100]"},
		{New([]byte{0xF0}, 4), 1, 11, "[11110000 00000010]"},
		{New([]byte{0xF0}, 4), 1, 12, "[11110000 00000001]"},
		{New([]byte{0xF0}, 4), 1, 13, "[11110000 00000000 10000000]"},
		// Add zero
		{New([]byte{0xF0}, 4), 0, 6, "[11110000 00000000]"},
	}

	for _, tt := range tests {
		lBefore := tt.ba.Len()
		tt.ba.Add32N(tt.in, tt.l)
		actual := fmt.Sprintf("%08b", tt.ba.Bytes())
		if actual != tt.expected {
			t.Errorf("%v => expected %s got %s", tt.in, tt.expected, actual)
		}
		expected := lBefore + int(tt.l)
		if tt.ba.Len() != expected {
			t.Errorf("%v => expected %d got %d", tt.in, expected, tt.ba.Len())
		}
	}
}

func TestAdd64N(t *testing.T) {
	tests := []struct {
		ba       *BitArray
		in       uint64
		l        uint64
		expected string
	}{
		{New([]byte{0xF0}, 4), 1, 1, "[11111000]"},
		{New([]byte{0xF0}, 4), 1, 2, "[11110100]"},
		{New([]byte{0xF0}, 4), 1, 3, "[11110010]"},
		{New([]byte{0xF0}, 4), 1, 4, "[11110001]"},
		{New([]byte{0xF0}, 4), 1, 5, "[11110000 10000000]"},
		{New([]byte{0xF0}, 4), 1, 6, "[11110000 01000000]"},
		{New([]byte{0xF0}, 4), 1, 7, "[11110000 00100000]"},
		{New([]byte{0xF0}, 4), 1, 8, "[11110000 00010000]"},
		{New([]byte{0xF0}, 4), 1, 9, "[11110000 00001000]"},
		{New([]byte{0xF0}, 4), 1, 10, "[11110000 00000100]"},
		{New([]byte{0xF0}, 4), 1, 11, "[11110000 00000010]"},
		{New([]byte{0xF0}, 4), 1, 12, "[11110000 00000001]"},
		{New([]byte{0xF0}, 4), 1, 13, "[11110000 00000000 10000000]"},
		// Add zero
		{New([]byte{0xF0}, 4), 0, 6, "[11110000 00000000]"},
	}

	for _, tt := range tests {
		lBefore := tt.ba.Len()
		tt.ba.Add64N(tt.in, tt.l)
		actual := fmt.Sprintf("%08b", tt.ba.Bytes())
		if actual != tt.expected {
			t.Errorf("%v => expected %s got %s", tt.in, tt.expected, actual)
		}
		expected := lBefore + int(tt.l)
		if tt.ba.Len() != expected {
			t.Errorf("%v => expected %d got %d", tt.in, expected, tt.ba.Len())
		}
	}
}
