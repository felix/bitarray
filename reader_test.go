package bitarray

import (
	"testing"
)

func TestReadBits(t *testing.T) {
	tests := map[string]struct {
		ba *BitArray
		// Seeking
		offset     int64
		dir        SeekFrom
		seekOffset int64
		// Reading
		count      int
		readOffset int64
		want       uint
	}{
		"fromStart": {
			ba:         NewFromBytes([]byte{0xf0, 0x01}, 16),
			offset:     0,
			dir:        SeekStart,
			seekOffset: 0,
			count:      4,
			readOffset: 4,
			want:       15,
		},
		"fromEnd": {
			ba:         NewFromBytes([]byte{0xf0, 0x01}, 16),
			offset:     4,
			dir:        SeekEnd,
			seekOffset: 12,
			count:      4,
			readOffset: 16,
			want:       1,
		},
		"fromCurrent": {
			ba:         NewFromBytes([]byte{0xf0, 0x01}, 16),
			offset:     12,
			dir:        SeekCurrent,
			seekOffset: 12,
			count:      4,
			readOffset: 16,
			want:       1,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r := NewReader(tt.ba)
			actualPos, err := r.Seek(tt.offset, tt.dir)
			if err != nil {
				t.Fatalf("failed to seek: %s", err)
			}
			if actualPos != tt.seekOffset {
				t.Errorf("got %d, want %d", actualPos, tt.seekOffset)
			}
			var actual uint
			if err := r.ReadBits(&actual, tt.count); err != nil {
				t.Fatalf("failed to read: %s", err)
			}
			if actual != tt.want {
				t.Errorf("got %d, want %d", actual, tt.want)
			}
			if r.Pos() != tt.readOffset {
				t.Errorf("got %d, want %d", r.Pos(), tt.readOffset)
			}

		})
	}
}
