package fuzzy

// #cgo CFLAGS: -I./c-fzy
// #include "match.h"
// #include "match.c"
import "C"
import (
	"math"
	"unsafe"
)

func FuzzyC() (C.double, []C.ulong) {
	n := C.CString("hi")
	h := C.CString("hello iam")
	defer C.free(unsafe.Pointer(n))
	defer C.free(unsafe.Pointer(h))

	numElements := C.strlen(n)
	size := C.size_t(numElements) * C.size_t(unsafe.Sizeof(C.ulong(0)))

	posPtr := (*C.ulong)(C.malloc(size))
	defer C.free(unsafe.Pointer(posPtr))

	f := C.match_positions(n, h, posPtr)
	s := unsafe.Slice(posPtr, numElements)

	return f, s
}

func MatchPositions(n, h string) (C.score_t, []int) {
	nc := C.CString(n)
	hc := C.CString(h)
	defer C.free(unsafe.Pointer(nc))
	defer C.free(unsafe.Pointer(hc))
	numElements := C.strlen(nc)
	size := C.size_t(numElements) * C.size_t(unsafe.Sizeof(C.size_t(0)))
	posPtr := (*C.size_t)(C.malloc(size))
	defer C.free(unsafe.Pointer(posPtr))
	f := C.match_positions(nc, hc, posPtr)
	//s := unsafe.Slice(posPtr, numElements)
	if float64(f) == math.Inf(-1) {
		return 0, nil
	}
	pos := make([]int, numElements)
	for i := 0; i < int(numElements); i++ {
		pos[i] = int(*(*C.size_t)(unsafe.Pointer(uintptr(unsafe.Pointer(posPtr)) + unsafe.Sizeof(*posPtr)*uintptr(i))))
	}
	return f, pos
}
