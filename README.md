# xhex

[![pkg.go.dev](https://pkg.go.dev/badge/github.com/templexxx/xhex.svg)](https://pkg.go.dev/github.com/templexxx/xhex)
[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Unit Test](https://github.com/templexxx/xhex/actions/workflows/unit-test.yml/badge.svg)](https://github.com/templexxx/xhex/actions/workflows/unit-test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/templexxx/xhex)](https://goreportcard.com/report/github.com/templexxx/xhex)

`xhex` provides high-performance hexadecimal encoding/decoding for Go.
On AVX2-capable x86_64 CPUs, it uses a SIMD fast path and falls back to a portable implementation otherwise.

## Highlights

- API-compatible behavior with Go's `encoding/hex` for byte-slice encode/decode use cases.
- Runtime CPU feature detection; AVX2 acceleration is enabled automatically when available.
- Significant throughput gains in large-buffer workloads.

## Usage

```go
package main

import "github.com/templexxx/xhex"

func main() {
	src := []byte("xhex")
	enc := make([]byte, len(src)*2)
	xhex.Encode(enc, src)

	dec := make([]byte, len(enc)/2)
	_ = xhex.Decode(dec, enc)
}
```

## Benchmark (vs `encoding/hex`)

Typical results:

```
Compare with standard lib:
benchmark                  old ns/op     new ns/op     delta
BenchmarkEncode/16-8       30.7          5.86          -80.91%
BenchmarkEncode/24-8       43.4          17.8          -58.99%
BenchmarkEncode/1024-8     1793          62.8          -96.50%

benchmark                  old MB/s     new MB/s     speedup
BenchmarkEncode/16-8       520.44       2732.67      5.25x
BenchmarkEncode/24-8       552.44       1349.15      2.44x
BenchmarkEncode/1024-8     571.10       16298.50     28.54x

benchmark                  old ns/op     new ns/op     delta
BenchmarkDecode/32-8       59.8          10.4          -82.61%
BenchmarkDecode/48-8       87.5          35.3          -59.66%
BenchmarkDecode/2048-8     3634          182           -94.99%

benchmark                  old MB/s     new MB/s     speedup
BenchmarkDecode/32-8       534.90       3074.74      5.75x
BenchmarkDecode/48-8       548.75       1359.05      2.48x
BenchmarkDecode/2048-8     563.56       11227.56     19.92x
```

Measured with project benchmarks under default settings.
Actual speedup depends on CPU model, Go version, and input size distribution.
