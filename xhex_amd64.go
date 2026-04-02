// Copyright (c) 2020. Temple3x (temple3x@gmail.com)
// Copyright (c) 2017 Zach Bjornson
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

import "github.com/templexxx/cpu"

func init() {
	if cpu.X86.HasAVX2 {
		encode = func(dst, src []byte) {
			n := len(src)
			if n == 0 {
				return
			}
			encodeAVX2(&dst[0], &src[0], n)
			done := n >> 4 << 4
			if done == n {
				return
			}

			// Deal with unaligned part.
			dst = dst[done*2:]
			src = src[done:]
			encodeBase(dst, src)
		}

		decode = func(dst, src []byte) error {
			n := len(src)
			if n == 0 {
				return nil
			}
			decodeAVX2(&dst[0], &src[0], n)
			done := n >> 5 << 5
			if done == n {
				return nil
			}
			// Deal with unaligned part.
			dst = dst[done/2:]
			src = src[done:]
			return decodeBase(dst, src)
		}
	}
}

// replaceHighMask shifts each high byte down by one position within every
// 16-bit lane and zeros the original high-byte slots.
// Example:
// Before packed shuffle bytes:
// [82 0 253 0 252 0 7 0 33 0 130 0 101 0 79 0 22 0 63 0 95 0 15 0 154 0 98 0 29 0 114 0]
// After packed shuffle bytes with this mask:
// [0 82 0 253 0 252 0 7 0 33 0 130 0 101 0 79 0 22 0 63 0 95 0 15 0 154 0 98 0 29 0 114]
//
// Equivalent unsigned form:
// []uint8{129, 0, 129, 2, 129, 4, 129, 6, 129, 8, 129, 10, 129, 12, 129, 14,
//
//	129, 0, 129, 2, 129, 4, 129, 6, 129, 8, 129, 10, 129, 12, 129, 14}
//
// Both forms produce identical results.
var replaceHighMask = []int8{-1, 0, -1, 2, -1, 4, -1, 6, -1, 8, -1, 10, -1, 12, -1, 14,
	-1, 0, -1, 2, -1, 4, -1, 6, -1, 8, -1, 10, -1, 12, -1, 14}

// decodeMask1 and decodeMask2 separate high/low nibbles and widen to 16-bit lanes.
var decodeMask1 = []int8{0, -1, 2, -1, 4, -1, 6, -1, 8, -1, 10, -1, 12, -1, 14, -1}
var decodeMask2 = []int8{1, -1, 3, -1, 5, -1, 7, -1, 9, -1, 11, -1, 13, -1, 15, -1}

// encodeAVX2 encodes src in blocks of 16 bytes using AVX2.
//
// The algorithm is based on:
// https://github.com/zbjornson/fast-hex
//
// This implementation keeps the assembly path reviewable and avoids AVX/SSE
// transition penalties in the hot loop.
//
//go:noescape
func encodeAVX2(dst, src *byte, n int)

// decodeAVX2 decodes src in blocks of 32 bytes using AVX2.
// The core idea also comes from https://github.com/zbjornson/fast-hex,
// with loop granularity adapted for this implementation.
//
//go:noescape
func decodeAVX2(dst, src *byte, n int)
