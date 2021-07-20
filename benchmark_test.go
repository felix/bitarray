package bitarray

import "testing"

var bmResult []byte

func BenchmarkSlice(b *testing.B) {
	ba := NewFromBytes([]byte{0x96, 0x2c, 0x49, 0x72, 0x2b, 0x80}, 48)
	var r *BitArray
	for n := 0; n < b.N; n++ {
		r, _ = ba.Slice(27, 7)
	}
	bmResult = r.Bytes()
}
