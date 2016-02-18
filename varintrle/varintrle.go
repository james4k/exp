// Package varintrle implements an algorithm similar to varint encoding,
// but for runs of integers.
//
// TODO: get into detail on the format. described briefly in WriteTo's
// comment.
//
// TODO: re-design to be friendly to vector and parallel computation.
package varintrle

import (
	"bytes"
	"errors"
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

// WriteRun writes to w a number of integer values using a modified
// varint encoding that is optimal for runs of integers of the same size
// requirements in bytes. Runs of zeros are especially optimal, taking
// O(1) space.
func WriteRun(w io.Writer, vals []int64) error {
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

// ReadRun reads all integers read from r into vals, up to len(vals),
// and returns the number of values read. If number of values read is
// less than len(vals), an error is returned. If the bytes decode into
// more values than len(vals), an error is returned.
func ReadRun(vals []int64, r io.Reader) (int, error) {
	var buf [8]byte
	pos := 0
	for {
		if pos >= len(vals) {
			return pos, nil
		}
		_, err := io.ReadFull(r, buf[:1])
		if err != nil {
			if err == io.EOF {
				return pos, nil
			}
			return pos, err
		}
		n, bytes := getnbytes(buf[0])
		if pos+n > len(vals) {
			return pos, errors.New("varintrle: unexpected values to read")
		}
		if bytes == 0 {
			for i := 0; i < n; i++ {
				vals[pos+i] = 0
			}
			pos += n
			continue
		}
		for i := 0; i < n; i++ {
			_, err := io.ReadFull(r, buf[:bytes])
			if err != nil {
				return pos, err
			}
			var val uint64
			for j := 0; j < bytes; j++ {
				val |= uint64(buf[j]) << (uint64(j) * 8)
			}
			vals[pos] = unzigzag(val)
			pos++
		}
	}
}
