package varintrle

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func randIntSlice(perms int) []int64 {
	max := [9]int64{
		0x00,
		0x7f,
		0x7fff,
		0x7fffff,
		0x7fffffff,
		0x7fffffffff,
		0x7fffffffffff,
		0x7fffffffffffff,
		0x7fffffffffffffff,
	}
	var vals []int64
	for i := 0; i < perms; i++ {
		n, bytes := rand.Intn(32)+1, rand.Intn(9)
		if bytes == 0 {
			for i := 0; i < n; i++ {
				vals = append(vals, 0)
			}
			continue
		}
		for i := 0; i < n; i++ {
			val := rand.Int63n(max[bytes])
			if rand.Intn(8) > 4 {
				val = -val
			}
			vals = append(vals, val)
		}
	}
	return vals
}

func TestWriteVarintRLE(t *testing.T) {
	table := map[int64][]byte{
		0:         []byte{nbytes(1, 0)},
		-1:        []byte{nbytes(1, 1), 1},
		1:         []byte{nbytes(1, 1), 2},
		-10:       []byte{nbytes(1, 1), 19},
		1000:      []byte{nbytes(1, 2), 2000 & 0xff, 2000 & 0xff00 >> 8},
		100000:    []byte{nbytes(1, 3), 200000 & 0xff, 200000 & 0xff00 >> 8, 200000 & 0xff0000 >> 16},
		0x1000000: []byte{nbytes(1, 4), 0x2000000 & 0xff, 0x2000000 & 0xff00 >> 8, 0x2000000 & 0xff0000 >> 16, 0x2000000 & 0xff000000 >> 24},
		0x1000000 << 32: []byte{nbytes(1, 8),
			(0x2000000 << 32) & 0xff, (0x2000000 << 32) & 0xff00 >> 8,
			(0x2000000 << 32) & (0xff << 16) >> 16, (0x2000000 << 32) & (0xff << 24) >> 24,
			(0x2000000 << 32) & (0xff << 32) >> 32, (0x2000000 << 32) & (0xff << 40) >> 40,
			(0x2000000 << 32) & (0xff << 48) >> 48, (0x2000000 << 32) & (0xff << 56) >> 56,
		},
	}
	buf := &bytes.Buffer{}
	for input, expected := range table {
		buf.Reset()
		err := WriteTo(buf, []int64{input})
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(buf.Bytes(), expected) {
			log.Println(buf.Bytes())
			t.Fatalf("did not match expected output for %v", input)
		}
	}
}

func TestVarintRLE(t *testing.T) {
	var vals, actual []int64
	vals = randIntSlice(20)
	buf := &bytes.Buffer{}
	err := WriteTo(buf, vals)
	if err != nil {
		t.Fatal(err)
	}
	actual, err = ReadFrom(actual, buf)
	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(vals, actual) {
		t.Fatalf("did not match expected output")
	}
}

func BenchmarkWriteVarintRLERandom(b *testing.B) {
	rand.Seed(1)
	vals := randIntSlice(100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := WriteTo(ioutil.Discard, vals)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReadVarintRLERandom(b *testing.B) {
	rand.Seed(1)
	vals := randIntSlice(100)
	buf := bytes.NewBuffer(make([]byte, 0, len(vals)*4))
	err := WriteTo(buf, vals)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(buf.Bytes())
		vals = vals[:0]
		_, err = ReadFrom(vals, r)
		if err != nil && err != io.EOF {
			b.Fatal(err)
		}
	}
}
