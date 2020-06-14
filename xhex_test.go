// Copyright (c) 2020. Temple3x (temple3x@gmail.com)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package xhex

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

type encDecTest struct {
	enc string
	dec []byte
}

var encDecTests = []encDecTest{
	{"", []byte{}}, // empty
	{"0001020304050607", []byte{0, 1, 2, 3, 4, 5, 6, 7}}, // unaligned
	{"000102030405060708090a0b0c0d0e0f",
		[]byte{0, 1, 2, 3, 4, 5, 6, 7,
			8, 9, 10, 11, 12, 13, 14, 15}}, // aligned
	{"000102030405060708090a0b0c0d0e0f010e",
		[]byte{0, 1, 2, 3, 4, 5, 6, 7,
			8, 9, 10, 11, 12, 13, 14, 15, 1, 14}}, // aligned + unaligned
	{"08090a0b0c0d0e0f", []byte{8, 9, 10, 11, 12, 13, 14, 15}},
	{"f0f1f2f3f4f5f6f7", []byte{0xf0, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7}},
	{"f8f9fafbfcfdfeff", []byte{0xf8, 0xf9, 0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff}},
	{"67", []byte{'g'}},
	{"e3a1", []byte{0xe3, 0xa1}},
}

func TestEncode(t *testing.T) {
	for i, test := range encDecTests {
		dst := make([]byte, len(test.dec)*2)
		Encode(dst, test.dec)

		if string(dst) != test.enc {
			t.Errorf("#%d: got: %#v want: %#v", i, string(dst), test.enc)
		}
	}

	// Compare with Go standard library.
	rand.Seed(time.Now().UnixNano())
	for i := 1; i < 1024; i++ {
		src := make([]byte, i)
		act := make([]byte, i*2)
		rand.Read(src)

		Encode(act, src)

		exp := make([]byte, i*2)
		hex.Encode(exp, src)

		if !bytes.Equal(exp, act) {
			t.Fatal("encode misamtch")
		}
	}
}

func TestDecode(t *testing.T) {
	// Case for decoding uppercase hex characters, since
	// Encode always uses lowercase.
	decTests := append(encDecTests, encDecTest{"F8F9FAFBFCFDFEFF", []byte{0xf8, 0xf9, 0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff}})
	for i, test := range decTests {
		dst := make([]byte, len(test.enc)/2)
		err := Decode(dst, []byte(test.enc))
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(dst, test.dec) {
			t.Errorf("#%d: got: %#v want: %#v", i, dst, test.enc)
		}
	}

	// Compare with Go standard library.
	rand.Seed(time.Now().UnixNano())
	for i := 1; i < 1024; i++ {
		p := make([]byte, i*2)
		act := make([]byte, i)
		rand.Read(act)
		Encode(p, act)
		err := Decode(act, p)
		if err != nil {
			t.Fatal(err)
		}

		exp := make([]byte, i)
		_, err = hex.Decode(exp, p)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(exp, act) {
			t.Fatal("decode misamtch")
		}
	}
}

var errTests = []struct {
	in  string
	out string
	err error
}{
	{"", "", nil},
	{"0", "", ErrLength},
	{"zd4aa", "", InvalidByteError('z')},
	{"d4aaz", "\xd4\xaa", InvalidByteError('z')},
	{"30313", "01", ErrLength},
	{"0g", "", InvalidByteError('g')},
	{"00gg", "\x00", InvalidByteError('g')},
	{"0\x01", "", InvalidByteError('\x01')},
	{"ffeed", "\xff\xee", ErrLength},
}

func TestDecodeErr(t *testing.T) {
	for _, tt := range errTests {
		out := make([]byte, len(tt.in)+10)
		err := Decode(out, []byte(tt.in))
		if err != tt.err {
			t.Errorf("Decode(%q) = %q, %v, want %q, %v", tt.in, string(out[:len(tt.out)]), err, tt.out, tt.err)
		}
	}
}

var sink []byte

func BenchmarkEncode(b *testing.B) {
	for _, size := range []int{16, 24, 1024} { // 16 for aligned, 24 for aligned+unaligned, 1024 for showing performance.
		src := bytes.Repeat([]byte{2, 3, 5, 7, 9, 11, 13, 17}, size/8)
		sink = make([]byte, 2*size)

		b.Run(fmt.Sprintf("%v", size), func(b *testing.B) {
			b.SetBytes(int64(size))
			for i := 0; i < b.N; i++ {
				Encode(sink, src)
			}
		})
	}
}

func BenchmarkDecode(b *testing.B) {
	for _, size := range []int{32, 48, 2048} {
		src := bytes.Repeat([]byte{'2', 'b', '7', '4', '4', 'f', 'a', 'a'}, size/8)
		sink = make([]byte, size/2)

		b.Run(fmt.Sprintf("%v", size), func(b *testing.B) {
			b.SetBytes(int64(size))
			for i := 0; i < b.N; i++ {
				hex.Decode(sink, src)
			}
		})
	}
}
