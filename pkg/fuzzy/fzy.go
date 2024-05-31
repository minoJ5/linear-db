// Port of John Hawthorn's fzy: https://github.com/jhawthorn/fzy
package fuzzy

import (
	"math"
	"strings"
)

var (
	SCORE_MIN = math.Inf(-1)
	SCORE_MAX = math.Inf(1)
)

const (
	SCORE_GAP_LEADING       = -0.005
	SCORE_GAP_TRAILING      = -0.005
	SCORE_GAP_INNER         = -0.01
	SCORE_MATCH_CONSECUTIVE = 1.0
	SCORE_MATCH_SLASH       = 0.8 // 0.9
	SCORE_MATCH_WORD        = 0.9 // 0.8
	SCORE_MATCH_CAPITAL     = 0.7
	SCORE_MATCH_DOT         = 0.6
)

func isLower(s string) bool {
	return strings.EqualFold(strings.ToLower(s), s)
}

func isUpper(s string) bool {
	return strings.EqualFold(strings.ToUpper(s), s)
}

func preBonus(haystack string) []float64 {
	m := len(haystack)
	matchBonus := make([]float64, m)
	lastChar := " "
	for i := 0; i < m; i++ {
		ch := string(haystack[i])
		if lastChar == "/" {
			matchBonus[i] = SCORE_MATCH_SLASH
		} else if lastChar == "-" || lastChar == "_" || lastChar == " " {
			matchBonus[i] = SCORE_MATCH_WORD
		} else if lastChar == "." {
			matchBonus[i] = SCORE_MATCH_DOT
		} else if isLower(lastChar) && isUpper(ch) {
			matchBonus[i] = SCORE_MATCH_CAPITAL
		} else {
			matchBonus[i] = 0
		}
		lastChar = ch
	}
	return matchBonus
}

func compute(needle, haystack string) ([][]float64, [][]float64) {
	n := len(needle)
	m := len(haystack)

	lowerNeedle := strings.ToLower(needle)
	lowerHaystack := strings.ToLower(haystack)

	matchBonus := preBonus(haystack)
	M := make([][]float64, n)
	D := make([][]float64, n)
	for i := 0; i < n; i++ {
		D[i] = make([]float64, m)
		M[i] = make([]float64, m)

		prevScore := SCORE_MIN
		gapScore := SCORE_GAP_INNER
		if i == n-1 {
			gapScore = SCORE_GAP_TRAILING
		}

		for j := 0; j < m; j++ {
			if lowerNeedle[i] == lowerHaystack[j] {
				score := SCORE_MIN
				if i == 0 {
					score = float64(j)*SCORE_GAP_LEADING + matchBonus[j]
				} else if j != 0 {
					score = math.Max(
						M[i-1][j-1]+matchBonus[j],
						D[i-1][j-1]+SCORE_MATCH_CONSECUTIVE,
					)
				}
				D[i][j] = score
				M[i][j] = score
				if score > prevScore+gapScore {
					prevScore = score
				} else {
					prevScore += gapScore
				}
			} else {
				D[i][j] = SCORE_MIN
				M[i][j] = prevScore + gapScore
			}
		}
	}
	return M, D
}

func Score(needle, haystack string) float64 {
	n := len(needle)
	m := len(haystack)

	if n == 0 || m == 0 {
		return SCORE_MIN
	}

	if n == m {
		return SCORE_MAX
	}

	if m > 1024 {
		return SCORE_MIN
	}

	M, _ := compute(needle, haystack)

	return M[n-1][m-1]
}

func Positions(needle, haystack string) []int {
	n := len(needle)
	m := len(haystack)

	positions := make([]int, n)

	if n == 0 || m == 0 {
		return positions
	}

	if n == m {
		for i := range positions {
			positions[i] = i
		}
		return positions
	}

	if m > 1024 {
		return positions
	}

	M, D := compute(needle, haystack)

	matchRequired := false

	for i := n - 1; i >= 0; i-- {
		for j := m - 1; j >= 0; j-- {
			if D[i][j] != SCORE_MIN && (!matchRequired || D[i][j] == M[i][j]) {
				matchRequired = i != 0 && j != 0 && M[i][j] == D[i-1][j-1]+SCORE_MATCH_CONSECUTIVE
				positions[i] = j
				j--
				break
			}
		}
	}

	return positions
}

func HasMatch(needle, haystack string) bool {
	needle = strings.ToLower(needle)
	haystack = strings.ToLower(haystack)

	l := len(needle)
	j := 0
	for i := 0; i < l; i++ {
		j = strings.IndexRune(haystack[j:], rune(needle[i])) + 1
		if j == 0 {
			return false
		}
	}
	return true
}
