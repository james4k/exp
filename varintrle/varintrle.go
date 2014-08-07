package varintrle

import (
	"bytes"
	"io"
)

func nbytes(n, bytes int) uint8 {
	if bytes == 8 {
		bytes = 7
	}
	return (uint8(n-1) & 0x1f << 3) | uint8(bytes)&0x7
}

func getnbytes(b uint8) (n, bytes int) {
	n = int((b>>3)&0x1f) + 1
	bytes = int(b & 0x7)
	if bytes == 7 {
		bytes = 8
	}
	return
}

func zigzag(x int64) uint64 {
	ux := uint64(x) << 1
	if x < 0 {
		ux = ^ux
	}
	return ux
}

func unzigzag(ux uint64) int64 {
	x := int64(ux >> 1)
	if ux&1 != 0 {
		x = ^x
	}
	return x
}

// writeVarintRLE encodes a number of integer values using a modified
// varint encoding that is optimal for runs of integers of the same size
// requirements in bytes. Runs of zeros are especially optimal, taking
// O(1) space.
func writeVarintRLE(w io.Writer, vals []int64) error {
	if len(vals) == 0 {
		return nil
	}
	var err error
	n := 0
	prevbytes := -1
	buf := &bytes.Buffer{}
	for i := range vals {
		v := zigzag(vals[i])
		bytes := 0
		// buffer until we hit different byte requirements
		var varbuf [8]byte
		for v != 0 {
			varbuf[bytes] = uint8(v & 0xff)
			v >>= 8
			bytes++
		}
		if n >= 1 && (bytes != prevbytes || n == 32) {
			_, err = w.Write([]byte{
				nbytes(n, prevbytes),
			})
			if err != nil {
				return err
			}
			if buf.Len() > 0 {
				_, err = w.Write(buf.Bytes())
				if err != nil {
					return err
				}
			}
			buf.Reset()
			n = 0
		}
		if bytes > 0 {
			if bytes == 7 {
				buf.Write(varbuf[:8])
			} else {
				buf.Write(varbuf[:bytes])
			}
		}
		n++
		prevbytes = bytes
	}
	if n > 0 {
		_, err = w.Write([]byte{
			nbytes(n, prevbytes),
		})
		if err != nil {
			return err
		}
		if buf.Len() > 0 {
			_, err = w.Write(buf.Bytes())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// readVarintRLE appends all integers read from r to vals and returns
// the new slice.
func readVarintRLE(vals []int64, r io.Reader) ([]int64, error) {
	var buf [8]byte
	for {
		_, err := io.ReadFull(r, buf[:1])
		if err != nil {
			if err == io.EOF {
				return vals, err
			}
			return nil, err
		}
		n, bytes := getnbytes(buf[0])
		if bytes == 0 {
			for i := 0; i < n; i++ {
				vals = append(vals, 0)
			}
			continue
		}
		for i := 0; i < n; i++ {
			_, err := io.ReadFull(r, buf[:bytes])
			if err != nil {
				return nil, err
			}
			var val uint64
			for j := 0; j < bytes; j++ {
				val |= uint64(buf[j]) << (uint64(j) * 8)
			}
			vals = append(vals, unzigzag(val))
		}
	}
	return vals, nil
}