# BitArray

This package provides a bit array structure with sub-byte methods like packing
and shifting for arrays of bits of arbitrary length.

This is not useful for set operations. It facilitates encoding data that is not
necessarily aligned to 8 bits. It was originally designed to cater to the
packed encoding required for the ISO/IEC 20248 digital signature format.

## Usage

```go
import (
	"fmt"

	"src.userspace.com.au/bitarray"
)

func main() {
	ba := &bitarray.BitArray{}

	// Add single bits, never truncated
	ba.AddBit(1)
	ba.AddBit(0)
	ba.AddBit(1)
	fmt.Println(ba.Len()) // => 3
	fmt.Printf("%08b\n", ba.Bytes()) // => [10100000]

	// Add with truncated leading padding
	ba.Add(5)
	fmt.Printf("%08b\n", ba.Bytes()) // => [10110100]

	// Add with fixed length
	ba.AddN(5, 10)
	fmt.Printf("%08b\n", ba.Bytes()) // => [10110100 00000101]

	r := bitarray.NewReader(ba)
	var out uint
	r.ReadBits(&out, 16)
	fmt.Printf("%08b\n", out) // => 1011010000000101

	ba.ShiftL(2)
	fmt.Println(ba.String()) // => [11010000 000101--]

	ba2, _ := ba.Slice(9, 6)
	fmt.Println(ba2.String()) // => [000101--]
}
```
