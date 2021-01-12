package bitarray

import (
	"fmt"
	"testing"
)

func TestAddBit(t *testing.T) {
	tests := []struct {
		ba       *BitArray
		in       uint
		expected string
	}{
		{NewFromBytes([]byte{0xF0}, 4), 1, "[11111000]"},
		{NewFromBytes([]byte{0xFF}, 8), 1, "[11111111 10000000]"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d", tt.in), func(t *testing.T) {
			lBefore := tt.ba.Len()
			tt.ba.AddBit(tt.in)
			actual := fmt.Sprintf("%08b", tt.ba.Bytes())
			if actual != tt.expected {
				t.Errorf("got %s, want %s", actual, tt.expected)
			}
			expected := lBefore + 1
			if tt.ba.Len() != expected {
				t.Errorf("got %d, want %d", tt.ba.Len(), expected)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	tests := []struct {
		ba        *BitArray
		in        []uint
		expected  string
		expectedN int64
	}{
		{NewFromBytes([]byte{0xF0}, 4), []uint{1}, "[11111000]", 5},
		{NewFromBytes([]byte{0xF0}, 8), []uint{1}, "[11110000 10000000]", 9},
		{NewFromBytes([]byte{0xF0}, 4), []uint{0}, "[11110000]", 5},
		{NewFromBytes([]byte{0xF0}, 4), []uint{0xf0ff, 0x0f}, "[11111111 00001111 11111111]", 24},
		{NewFromBytes([]byte{0xF0}, 8), []uint{0xf0ff, 0x0f}, "[11110000 11110000 11111111 11110000]", 28},
		{NewFromBytes([]byte{0xF0}, 8), []uint{0x0f0fff}, "[11110000 11110000 11111111 11110000]", 28},
		{NewFromBytes([]byte{0xF0}, 8), []uint{0x81bf0fff}, "[11110000 10000001 10111111 00001111 11111111]", 40},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d", tt.in), func(t *testing.T) {
			for _, i := range tt.in {
				tt.ba.Add(i)
			}
			actual := fmt.Sprintf("%08b", tt.ba.Bytes())
			if actual != tt.expected {
				t.Errorf("got %s, want %s", actual, tt.expected)
			}
			if tt.ba.Len() != tt.expectedN {
				t.Errorf("got len=%d, want %d", tt.ba.Len(), tt.expectedN)
			}
		})
	}
}

func TestAddN(t *testing.T) {
	tests := []struct {
		ba       *BitArray
		in       uint
		l        int
		expected string
	}{
		{NewFromBytes([]byte{0xF0}, 4), 1, 1, "[11111000]"},
		{NewFromBytes([]byte{0xF0}, 4), 1, 2, "[11110100]"},
		{NewFromBytes([]byte{0xF0}, 4), 1, 3, "[11110010]"},
		{NewFromBytes([]byte{0xF0}, 4), 1, 4, "[11110001]"},
		{NewFromBytes([]byte{0xF0}, 4), 1, 5, "[11110000 10000000]"},
		{NewFromBytes([]byte{0xF0}, 4), 1, 6, "[11110000 01000000]"},
		{NewFromBytes([]byte{0xF0}, 4), 1, 7, "[11110000 00100000]"},
		{NewFromBytes([]byte{0xF0}, 4), 1, 8, "[11110000 00010000]"},
		{NewFromBytes([]byte{0xF0}, 4), 1, 9, "[11110000 00001000]"},
		{NewFromBytes([]byte{0xF0}, 4), 1, 10, "[11110000 00000100]"},
		{NewFromBytes([]byte{0xF0}, 4), 1, 11, "[11110000 00000010]"},
		{NewFromBytes([]byte{0xF0}, 4), 1, 12, "[11110000 00000001]"},
		{NewFromBytes([]byte{0xF0}, 4), 1, 13, "[11110000 00000000 10000000]"},
		// Add zero
		{NewFromBytes([]byte{0xF0}, 4), 0, 6, "[11110000 00000000]"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d", tt.in), func(t *testing.T) {
			lBefore := tt.ba.Len()
			n := tt.ba.AddN(tt.in, tt.l)
			actual := fmt.Sprintf("%08b", tt.ba.Bytes())
			if actual != tt.expected {
				t.Errorf("got %s, want %s", actual, tt.expected)
			}
			if n != tt.l {
				t.Errorf("got n=%d, want %d", n, tt.l)
			}
			expected := lBefore + int64(tt.l)
			if tt.ba.Len() != expected {
				t.Errorf("got len=%d, want %d", tt.ba.Len(), expected)
			}
		})
	}
}

func TestSlice(t *testing.T) {
	tests := []struct {
		ba       *BitArray
		s, l     int64 // start and length
		expected string
		avail    uint
	}{
		{NewFromBytes([]byte{0xff}, 8), 0, 8, "[11111111]", 0},
		{NewFromBytes([]byte{0xff}, 8), 0, 1, "[10000000]", 7},
		{NewFromBytes([]byte{0xfe}, 8), 0, 8, "[11111110]", 0},
		{NewFromBytes([]byte{0x03}, 8), 7, 1, "[10000000]", 7},
		{NewFromBytes([]byte{0xd0}, 4), 0, 4, "[11010000]", 4},
		// Multiple bytes
		{NewFromBytes([]byte{0xd0, 0xff}, 16), 0, 9, "[11010000 10000000]", 7},
		{NewFromBytes([]byte{0x0f, 0xf0}, 16), 4, 8, "[11111111]", 0},
		// Cases
		// 10010110 00101100 01001001 => 1000101
		{NewFromBytes([]byte{0x96, 0x2c, 0x49}, 24), 6, 7, "[10001010]", 1},
		// 10010110 00101100 01001001 01110010 00101011 10000000
		//                               ^^^^^ ^^
		{NewFromBytes([]byte{0x96, 0x2c, 0x49, 0x72, 0x2b, 0x80}, 48), 27, 7, "[10010000]", 1},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("%x[%d:%d]", tt.ba.raw, tt.s, tt.l)
		t.Run(name, func(t *testing.T) {
			a, err := tt.ba.Slice(tt.s, tt.l)
			if err != nil {
				t.Errorf("failed with %q", err)
			}
			actual := fmt.Sprintf("%08b", a.Bytes())
			if actual != tt.expected {
				t.Errorf("got %q, want %q", actual, tt.expected)
			}
			if a.avail() != tt.avail {
				t.Errorf("got avail %d, want %d", a.avail(), tt.avail)
			}
		})
	}
}

func TestReadUint(t *testing.T) {
	tests := []struct {
		ba       *BitArray
		s, l     int64 // start and length
		expected uint
	}{
		{NewFromBytes([]byte{0x02}, 8), 0, 8, 0x02},
		{NewFromBytes([]byte{0xff}, 8), 0, 8, 0xff},
		{NewFromBytes([]byte{0xff}, 8), 0, 1, 0x01},
		{NewFromBytes([]byte{0xfe}, 8), 0, 8, 0xfe},
		{NewFromBytes([]byte{0x03}, 8), 7, 1, 0x01},
		{NewFromBytes([]byte{0xd0}, 8), 0, 4, 0x0d},
		// Multiple bytes
		{NewFromBytes([]byte{0xd0, 0xff}, 16), 0, 9, 0x1a1},
		{NewFromBytes([]byte{0x0f, 0xf0}, 16), 4, 8, 0xff},
	}

	for _, tt := range tests {
		name := fmt.Sprintf("%x[%d:%d]", tt.ba.raw, tt.s, tt.l)
		t.Run(name, func(t *testing.T) {
			actual, err := tt.ba.ReadUint(tt.s, tt.l)
			if err != nil {
				t.Errorf("failed with %q", err)
			}
			if actual != tt.expected {
				t.Errorf("got %d, want %d", actual, tt.expected)
			}
		})
	}
}

func TestTest(t *testing.T) {
	tests := []struct {
		ba       *BitArray
		in       int64
		expected bool
	}{
		{NewFromBytes([]byte{0x01}, 8), 0, false},
		{NewFromBytes([]byte{0x01}, 8), 1, false},
		{NewFromBytes([]byte{0x01}, 8), 7, true},
		{NewFromBytes([]byte{0x01}, 8), 8, false},
		{NewFromBytes([]byte{0x00, 0x80}, 16), 8, true},
		{NewFromBytes([]byte{0x00, 0x02}, 16), 14, true},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d", tt.in), func(t *testing.T) {
			actual := tt.ba.Test(tt.in)
			if actual != tt.expected {
				t.Errorf("got %t, want %t", actual, tt.expected)
			}
		})
	}
}

func TestSet(t *testing.T) {
	tests := []struct {
		ba       *BitArray
		in       int64
		expected string
	}{
		{NewFromBytes([]byte{0x00}, 1), 0, "[10000000]"},
		{NewFromBytes([]byte{0x80}, 8), 1, "[11000000]"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d", tt.in), func(t *testing.T) {
			tt.ba.Set(tt.in)
			actual := fmt.Sprintf("%08b", tt.ba.Bytes())
			if actual != tt.expected {
				t.Errorf("got %s, want %s", actual, tt.expected)
			}
		})
	}
}

func TestShiftL(t *testing.T) {
	tests := []struct {
		ba       *BitArray
		in       uint
		expected string
	}{
		{NewFromBytes([]byte{0x01}, 7), 1, "[00000010]"},
		{NewFromBytes([]byte{0x01}, 7), 2, "[00000100]"},
		{NewFromBytes([]byte{0x01}, 7), 3, "[00001000]"},
		{NewFromBytes([]byte{0x01}, 7), 4, "[00010000]"},
		{NewFromBytes([]byte{0x01}, 7), 5, "[00100000]"},
		{NewFromBytes([]byte{0x01}, 7), 6, "[01000000]"},
		{NewFromBytes([]byte{0x01}, 7), 7, "[10000000]"},
		{NewFromBytes([]byte{0x01}, 7), 8, "[00000000]"},
		// Across a byte
		{NewFromBytes([]byte{0x01, 0x01}, 7), 5, "[00100000 00100000]"},
		{NewFromBytes([]byte{0x01, 0x01}, 7), 8, "[00000001 00000000]"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d", tt.in), func(t *testing.T) {
			tt.ba.ShiftL(tt.in)
			actual := fmt.Sprintf("%08b", tt.ba.Bytes())
			if actual != tt.expected {
				t.Errorf("got %s, want %s", actual, tt.expected)
			}
		})
	}
}

func TestShiftR(t *testing.T) {
	tests := []struct {
		ba       *BitArray
		in       uint
		expected string
	}{
		{NewFromBytes([]byte{0x80}, 8), 1, "[01000000]"},
		{NewFromBytes([]byte{0x80}, 8), 2, "[00100000]"},
		{NewFromBytes([]byte{0x80}, 8), 3, "[00010000]"},
		{NewFromBytes([]byte{0x80}, 8), 4, "[00001000]"},
		{NewFromBytes([]byte{0x80}, 8), 5, "[00000100]"},
		{NewFromBytes([]byte{0x80}, 8), 6, "[00000010]"},
		{NewFromBytes([]byte{0x80}, 8), 7, "[00000001]"},
		{NewFromBytes([]byte{0x80}, 8), 8, "[00000000]"},
		// Across a byte
		{NewFromBytes([]byte{0x80, 0x80}, 16), 5, "[00000100 00000100]"},
		{NewFromBytes([]byte{0x80, 0x80}, 16), 8, "[00000000 10000000]"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d", tt.in), func(t *testing.T) {
			tt.ba.ShiftR(tt.in)
			actual := fmt.Sprintf("%08b", tt.ba.Bytes())
			if actual != tt.expected {
				t.Errorf("got %s, want %s", actual, tt.expected)
			}
		})
	}
}

func TestPack(t *testing.T) {
	tests := []struct {
		in       interface{}
		expected string
	}{
		// Types
		{uint(108), "[1101100-]"},
		{uint8(108), "[1101100-]"},
		{uint16(108), "[1101100-]"},
		{uint32(108), "[1101100-]"},
		{uint64(108), "[1101100-]"},
		{int(108), "[1101100-]"},
		{int8(0xf), "[1111----]"},
		{int16(108), "[1101100-]"},
		{int32(108), "[1101100-]"},
		{int64(108), "[1101100-]"},
		{[]byte{0x6c}, "[1101100-]"},
		{[]uint{108}, "[1101100-]"},
		{[]int{108}, "[1101100-]"},
		{[]int8{108}, "[1101100-]"},
		{[]int16{108}, "[1101100-]"},
		{[]int32{108}, "[1101100-]"},
		{[]int64{108}, "[1101100-]"},
		// Zero
		{uint(0), "[0-------]"},
		{uint8(0), "[0-------]"},
		{uint16(0), "[0-------]"},
		{uint32(0), "[0-------]"},
		{uint64(0), "[0-------]"},
		{int(0), "[0-------]"},
		{int8(0), "[0-------]"},
		{int16(0), "[0-------]"},
		{int32(0), "[0-------]"},
		{int64(0), "[0-------]"},
		{[]byte{0}, "[0-------]"},
		{[]uint{0}, "[0-------]"},
		{[]int{0}, "[0-------]"},
		{[]int8{0}, "[0-------]"},
		{[]int16{0}, "[0-------]"},
		{[]int32{0}, "[0-------]"},
		{[]int64{0}, "[0-------]"},
		// Adding a zero
		{[]int{1, 0, 23}, "[1010111-]"},
		{[]uint16{1, 0, 23}, "[1010111-]"},
		{[]uint32{1, 0, 23}, "[1010111-]"},
		{[]uint64{1, 0, 23}, "[1010111-]"},
		// Cases
		{[]uint8{0xff, 0xff}, "[11111111 11111111]"},
		{[]uint8{0xff, 0xf0}, "[11111111 11110000]"},
		{[]uint8{0xf0, 0xf0}, "[11110000 11110000]"},
		{[]uint8{0xf0, 0xf0, 1}, "[11110000 11110000 1-------]"},
		{[]int{1, 128, 23}, "[11000000 010111--]"},
		{[]int{1, 129, 23}, "[11000000 110111--]"},
		{[]interface{}{uint8(0xff), 1, 2, 1, 4, 1, 1}, "[11111111 11011001 1-------]"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.in), func(t *testing.T) {
			ba, err := Pack(tt.in)
			if err != nil {
				t.Fatalf("failed %q", err)
			}
			actual := ba.String()
			if actual != tt.expected {
				t.Errorf("got %s, want %s", actual, tt.expected)
			}
			if actual != ba.String() {
				t.Errorf("got %q, want %q", ba.String(), actual)
			}
		})
	}
}

func TestAppend(t *testing.T) {
	tests := []struct {
		ba1      *BitArray
		ba2      *BitArray
		expected string
	}{
		{NewFromBytes([]byte{0x80}, 1), NewFromBytes([]byte{0x80}, 1), "[11000000]"},
		{NewFromBytes([]byte{0x80}, 8), NewFromBytes([]byte{0x80}, 1), "[10000000 10000000]"},
		{NewFromBytes([]byte{0xF0}, 4), NewFromBytes([]byte{0xF0}, 4), "[11111111]"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s+%s", tt.ba1, tt.ba2), func(t *testing.T) {
			fmt.Printf("%08b", tt.ba1.Bytes())
			fmt.Printf("%08b", tt.ba2.Bytes())
			tt.ba1.Append(*tt.ba2)
			actual := fmt.Sprintf("%08b", tt.ba1.Bytes())
			if actual != tt.expected {
				t.Errorf("got %s, want %s", actual, tt.expected)
			}
		})
	}
}
