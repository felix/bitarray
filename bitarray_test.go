package bitarray

import (
	"fmt"
	"testing"
)

func TestAdd8(t *testing.T) {
	tests := []struct {
		in       []uint8
		expected string
	}{
		{[]uint8{0xFF, 0xFF}, "[11111111 11111111]"},
		{[]uint8{1, 128, 23}, "[11000000 01011100]"},
		{[]uint8{0xFF, 1, 2, 1, 4, 1, 1}, "[11111111 11011001 10000000]"},
	}

	ba := BitArray{}

	for _, tt := range tests {
		for _, i := range tt.in {
			ba.Add8(i)
		}
		actual := fmt.Sprintf("%08b", ba.Bytes())
		if actual != tt.expected {
			t.Errorf("expected %q got %q", tt.expected, actual)
		}
		ba.Reset()
	}
}

func TestAdd16(t *testing.T) {
	tests := []struct {
		in       []uint16
		expected string
	}{
		{[]uint16{0xFF, 0xFF}, "[11111111 11111111]"},
		{[]uint16{1, 128, 23}, "[11000000 01011100]"},
		{[]uint16{0xFF, 1, 2, 1, 4, 1, 1}, "[11111111 11011001 10000000]"},
		{[]uint16{0xFFFF, 1}, "[11111111 11111111 10000000]"},
	}

	ba := BitArray{}

	for _, tt := range tests {
		for _, i := range tt.in {
			ba.Add16(i)
		}
		actual := fmt.Sprintf("%08b", ba.Bytes())
		if actual != tt.expected {
			t.Errorf("expected %q got %q", tt.expected, actual)
		}
		ba.Reset()
	}
}

func TestAdd32(t *testing.T) {
	tests := []struct {
		in       []uint32
		expected string
	}{
		{[]uint32{0xFF, 0xFF}, "[11111111 11111111]"},
		{[]uint32{1, 128, 23}, "[11000000 01011100]"},
		{[]uint32{0xFF, 1, 2, 1, 4, 1, 1}, "[11111111 11011001 10000000]"},
		{[]uint32{0xFFFFFF, 1}, "[11111111 11111111 11111111 10000000]"},
	}

	ba := BitArray{}

	for _, tt := range tests {
		for _, i := range tt.in {
			ba.Add32(i)
		}
		actual := fmt.Sprintf("%08b", ba.Bytes())
		if actual != tt.expected {
			t.Errorf("expected %q got %q", tt.expected, actual)
		}
		ba.Reset()
	}
}

func TestAdd64(t *testing.T) {
	tests := []struct {
		in       []uint64
		expected string
	}{
		{[]uint64{0xFF, 0xFF}, "[11111111 11111111]"},
		{[]uint64{1, 128, 23}, "[11000000 01011100]"},
		{[]uint64{0xFF, 1, 2, 1, 4, 1, 1}, "[11111111 11011001 10000000]"},
		{[]uint64{0xFFFFFFFF, 1}, "[11111111 11111111 11111111 11111111 10000000]"},
	}

	ba := BitArray{}

	for _, tt := range tests {
		for _, i := range tt.in {
			ba.Add64(i)
		}
		actual := fmt.Sprintf("%08b", ba.Bytes())
		if actual != tt.expected {
			t.Errorf("expected %q got %q", tt.expected, actual)
		}
		ba.Reset()
	}
}

func TestRead8(t *testing.T) {
	tests := []struct {
		ba       *BitArray
		s, l     uint64 // start and length
		expected string
	}{
		{&BitArray{raw: []byte{0xff}, bits: 8}, 0, 8, "11111111"},
		{&BitArray{raw: []byte{0xff}, bits: 8}, 0, 1, "00000001"},
		{&BitArray{raw: []byte{0xfe}, bits: 8}, 0, 8, "11111110"},
		{&BitArray{raw: []byte{0x03}, bits: 8}, 7, 1, "00000001"},
		{&BitArray{raw: []byte{0xd0}, bits: 4}, 0, 4, "00001101"},
	}

	for _, tt := range tests {
		a, err := tt.ba.Read8(tt.s, tt.l)
		if err != nil {
			t.Errorf("failed with %q", err)
		}
		actual := fmt.Sprintf("%08b", a)
		fmt.Printf("expected=%s actual=%s\n", tt.expected, actual)
		if actual != tt.expected {
			t.Errorf("expected %q got %q", tt.expected, actual)
		}
	}
}
