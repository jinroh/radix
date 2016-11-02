package radix

// Inspired and adapted from crypto/cipher/xor.go

import (
	"runtime"
	"unsafe"
)

const wordSize = int(unsafe.Sizeof(uintptr(0)))
const supportsUnaligned = runtime.GOARCH == "386" || runtime.GOARCH == "amd64" || runtime.GOARCH == "ppc64" || runtime.GOARCH == "ppc64le" || runtime.GOARCH == "s390x"

// fastXORStrings xors in bulk. It only works on architectures that
// support read.
func fastXORStrings(a, b string) (bool, int, byte) {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}

	w := n / wordSize
	if w > 0 {
		aw := *(*[]uintptr)(unsafe.Pointer(&a))
		bw := *(*[]uintptr)(unsafe.Pointer(&b))
		for i := 0; i < w; i++ {
			if dw := aw[i] ^ bw[i]; dw != 0 {
				var j int
				for j = 0; dw&0xFF == 0; j++ {
					dw = dw >> 8
				}
				return false, i*wordSize + j, byte(dw)
			}
		}
	}

	for i := (n - n%wordSize); i < n; i++ {
		if d := a[i] ^ b[i]; d != 0 {
			return false, i, d
		}
	}

	return true, n, 0
}

func safeXORStrings(a, b string) (bool, int, byte) {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}

	for i := 0; i < n; i++ {
		if d := a[i] ^ b[i]; d != 0 {
			return false, i, d
		}
	}

	return true, n, 0
}

// xorBytes xors the bytes in a and b. The destination is assumed to
// have enough space. Returns the number of bytes xor'd.
func xorStrings(a, b string) (bool, int, byte) {
	if supportsUnaligned {
		return fastXORStrings(a, b)
	}
	return safeXORStrings(a, b)
}
