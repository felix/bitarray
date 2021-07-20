package bitarray

import (
	"fmt"
	"testing"
)

func TestAddBit(t *testing.T) {
	tests := map[string]struct {
		ba       *BitArray
		in       uint
		expected string
	}{
		"normal":     {New(SetBytes([]byte{0xF0}), SetSize(4)), 1, "[11111---]"},
		"newByte":    {New(SetBytes([]byte{0xFF})), 1, "[11111111 1-------]"},
		"resizedNew": {New(SetBytes([]byte{0xFF}), SetSize(16)), 1, "[11111111 00000000 1-------]"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			lBefore := tt.ba.Len()
			tt.ba.AddBit(tt.in)
			actual := tt.ba.String()
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
	tests := map[string]struct {
		ba        *BitArray
		in        []uint
		expected  string
		expectedN int64
	}{
		"addBit":        {NewFromBytes([]byte{0xF0}, 4), []uint{1}, "[11111---]", 5},
		"acrossByte":    {NewFromBytes([]byte{0xF0}, 8), []uint{1}, "[11110000 1-------]", 9},
		"addZero":       {NewFromBytes([]byte{0xF0}, 4), []uint{0}, "[11110---]", 5},
		"multiple":      {NewFromBytes([]byte{0xF0}, 4), []uint{0xf0ff, 0x0f}, "[11111111 00001111 11111111]", 24},
		"multipleLarge": {NewFromBytes([]byte{0xF0}, 8), []uint{0xf0ff, 0x0f}, "[11110000 11110000 11111111 1111----]", 28},
		"singleLarge":   {NewFromBytes([]byte{0xF0}, 8), []uint{0x0f0fff}, "[11110000 11110000 11111111 1111----]", 28},
		"singleLarge2":  {NewFromBytes([]byte{0xF0}, 8), []uint{0x81bf0fff}, "[11110000 10000001 10111111 00001111 11111111]", 40},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			for _, i := range tt.in {
				tt.ba.Add(i)
			}
			actual := tt.ba.String()
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
	tests := map[string]struct {
		ba       *BitArray
		in       uint
		l        int
		expected string
	}{
		"single":           {NewFromBytes([]byte{0xF0}, 4), 1, 1, "[11111---]"},
		"fullByte":         {NewFromBytes([]byte{0xF0}, 4), 1, 4, "[11110001]"},
		"acrossByte":       {NewFromBytes([]byte{0xF0}, 4), 1, 5, "[11110000 1-------]"},
		"twoFullBytes":     {NewFromBytes([]byte{0xF0}, 4), 1, 12, "[11110000 00000001]"},
		"acrossTwo":        {NewFromBytes([]byte{0xF0}, 4), 1, 13, "[11110000 00000000 1-------]"},
		"addZero":          {NewFromBytes([]byte{0xF0}, 4), 0, 6, "[11110000 00------]"},
		"truncate":         {NewFromBytes([]byte{0xF0}, 4), 255, 1, "[11111---]"},
		"truncatedOnEmpty": {NewFromBytes([]byte{0x00}, 0), 255, 2, "[11------]"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			lBefore := tt.ba.Len()
			n := tt.ba.AddN(tt.in, tt.l)
			actual := tt.ba.String()
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
	tests := map[string]struct {
		ba       *BitArray
		s, l     int64 // start and length
		expected string
		avail    uint
	}{
		"fullFromStart":    {NewFromBytes([]byte{0xff}, 8), 0, 8, "[11111111]", 0},
		"single":           {NewFromBytes([]byte{0xff}, 8), 0, 1, "[1-------]", 7},
		"toEnd":            {NewFromBytes([]byte{0x03}, 8), 7, 1, "[1-------]", 7},
		"middle":           {NewFromBytes([]byte{0xd0}, 4), 0, 4, "[1101----]", 4},
		"multipleBytes":    {NewFromBytes([]byte{0xd0, 0xff}, 16), 0, 9, "[11010000 1-------]", 7},
		"multipleToSingle": {NewFromBytes([]byte{0x0f, 0xf0}, 16), 4, 8, "[11111111]", 0},
		// 10010110 00101100 01001001 => 1000101
		"testCase1": {NewFromBytes([]byte{0x96, 0x2c, 0x49}, 24), 6, 7, "[1000101-]", 1},
		// 10010110 00101100 01001001 01110010 00101011 10000000
		//                               ^^^^^ ^^
		"testCase2":  {NewFromBytes([]byte{0x96, 0x2c, 0x49, 0x72, 0x2b, 0x80}, 48), 27, 7, "[1001000-]", 1},
		"zeroLength": {NewFromBytes([]byte{0, 0xFF, 0xFF}, 24), 8, 0, "[]", 0},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			a, err := tt.ba.Slice(tt.s, tt.l)
			if err != nil {
				t.Errorf("failed with %q", err)
			}
			actual := a.String()
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
	tests := map[string]struct {
		ba       *BitArray
		s, l     int64 // start and length
		expected uint
	}{
		"fullByte":      {NewFromBytes([]byte{0x02}, 8), 0, 8, 0x02},
		"singleBit":     {NewFromBytes([]byte{0xff}, 8), 0, 1, 0x01},
		"endOfByte":     {NewFromBytes([]byte{0x03}, 8), 7, 1, 0x01},
		"halfByte":      {NewFromBytes([]byte{0xd0}, 8), 0, 4, 0x0d},
		"multipleBytes": {NewFromBytes([]byte{0xd0, 0xff}, 16), 0, 9, 0x1a1},
		"acrossBytes":   {NewFromBytes([]byte{0x0f, 0xf0}, 16), 4, 8, 0xff},
	}

	for name, tt := range tests {
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
	tests := map[string]struct {
		ba       *BitArray
		in       int64
		expected bool
	}{
		"offset1":        {NewFromBytes([]byte{0x01}, 8), 0, false},
		"offset2":        {NewFromBytes([]byte{0x01}, 8), 1, false},
		"endOfByte":      {NewFromBytes([]byte{0x01}, 8), 7, true},
		"outOfRange":     {NewFromBytes([]byte{0x01}, 8), 8, false},
		"multipleBytes":  {NewFromBytes([]byte{0x00, 0x80}, 16), 8, true},
		"multipleBytes2": {NewFromBytes([]byte{0x00, 0x02}, 16), 14, true},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.ba.Test(tt.in)
			if actual != tt.expected {
				t.Errorf("got %t, want %t", actual, tt.expected)
			}
		})
	}
}

func TestSet(t *testing.T) {
	tests := map[string]struct {
		ba       *BitArray
		in       int64
		expected string
	}{
		"start":        {NewFromBytes([]byte{0x00}, 1), 0, "[10000000]"},
		"withExisting": {NewFromBytes([]byte{0x80}, 8), 1, "[11000000]"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.ba.Set(tt.in)
			actual := fmt.Sprintf("%08b", tt.ba.Bytes())
			if actual != tt.expected {
				t.Errorf("got %s, want %s", actual, tt.expected)
			}
		})
	}
}

func TestShiftL(t *testing.T) {
	tests := map[string]struct {
		ba       *BitArray
		shift    int64
		expected string
	}{
		"short":              {NewFromBytes([]byte{0x01}, 8), 1, "[0000001-]"},
		"longer":             {NewFromBytes([]byte{0x01}, 8), 7, "[1-------]"},
		"zero":               {NewFromBytes([]byte{0x01}, 8), 0, "[00000001]"},
		"empty":              {New(), 0, "[]"},
		"fullByte":           {NewFromBytes([]byte{0x01}, 8), 8, "[]"},
		"acrossByte":         {NewFromBytes([]byte{0x01, 0x01}, 16), 5, "[00100000 001-----]"},
		"acrossBytesTrimmed": {NewFromBytes([]byte{0x01, 0x01}, 16), 8, "[00000001]"},
		"moreThan8":          {NewFromBytes([]byte{0x00, 0x01}, 16), 11, "[00001---]"},
		"moreThan16":         {NewFromBytes([]byte{0x00, 0x00, 0x01}, 24), 22, "[01------]"},
		"shiftMoreThanSize":  {NewFromBytes([]byte{0x00, 0x01}, 16), 22, "[]"},
		"shiftZero":          {NewFromBytes([]byte{0x00, 0x01}, 16), 0, "[00000000 00000001]"},
		"trimmed":            {NewFromBytes([]byte{0x00, 0x01}, 16), 8, "[00000001]"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.ba.ShiftL(tt.shift)
			actual := tt.ba.String()
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
	tests := map[string]struct {
		ba1      *BitArray
		ba2      *BitArray
		expected string
	}{
		"singleBites":   {NewFromBytes([]byte{0x80}, 1), NewFromBytes([]byte{0x80}, 1), "[11000000]"},
		"fullAndSingle": {NewFromBytes([]byte{0x80}, 8), NewFromBytes([]byte{0x80}, 1), "[10000000 10000000]"},
		"twoHalves":     {NewFromBytes([]byte{0xF0}, 4), NewFromBytes([]byte{0xF0}, 4), "[11111111]"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			tt.ba1.Append(*tt.ba2)
			actual := fmt.Sprintf("%08b", tt.ba1.Bytes())
			if actual != tt.expected {
				t.Errorf("got %s, want %s", actual, tt.expected)
			}
		})
	}
}
