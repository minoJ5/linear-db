package fuzzy

import (
	"fmt"
)

const (
	CHAR_MAX = 255
)

// Use strings.Contains instead
func BitapSearch(text, pattern string) (int, error) {

	m := len(pattern)
	if m == 0 {
		return 0, nil
	}
	if m > 64 {
		return -1, fmt.Errorf("pattern %s is too long: length %d, max is 64 bits", pattern, m)
	}
	patternMask := make([]uint64, CHAR_MAX+1)
	for i := 0; i < m; i++ {
		patternMask[(pattern)[i]] |= 1 << uint(i)
	}
	var R uint64 = 1
	for i := 0; i < len(text); i++ {
		R = (R << 1) | 1
		R &= patternMask[(text)[i]]
		if R&(1<<uint(m-1)) != 0 {
			return i - m + 1, nil
		}
	}
	return -1, nil
}

func FuzzyBitapSearch(text, pattern string, k int) (int, error) {

	m := len(pattern)
	if m == 0 {
		return 0, nil
	}
	if m > 64 {
		return -1, fmt.Errorf("pattern %s is too long: length %d, max is 64 bits", pattern, m)
	}
	patternMask := make([]uint64, CHAR_MAX+1)

	R := make([]uint64, k+1)
	for i := 0; i <= k; i++ {
		R[i] = ^uint64(1)
	}
	for i := 0; i <= CHAR_MAX; i++ {
		patternMask[i] = ^uint64(0)
	}
	for i := 0; i < m; i++ {
		patternMask[(pattern)[i]] &= ^(uint64(1) << i)
	}
	for i := 0; i < len(text); i++ {
		var oldRd1 uint64 = R[0]
		R[0] |= patternMask[(text)[i]]
		R[0] <<= 1
		for d := 1; d <= k; d++ {
			var tmp uint64 = R[d]
			R[d] = (oldRd1 & (R[d] | patternMask[(text)[i]])) << 1
			oldRd1 = tmp
		}
		if (R[k] & (uint64(1) << m)) == 0 {
			return i - m + 1, nil
		}
	}
	return -1, nil
}
