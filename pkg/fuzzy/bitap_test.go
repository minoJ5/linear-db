package fuzzy

import (
	"fmt"
	"testing"
)

func TestBitap(t *testing.T) {
	text := "aklelpode dfeeww"
	pattern := "eew"
	res, err := BitapSearch(&text, &pattern)
	if err != nil || res != 12 {
		t.Fatalf("error: %s with result %d", err, res)
	}
	fmt.Println(res)
	textf := "abcdefghijklmnopqrstuvwxyz"
	patternf := "bxdef"
	resf, err := FuzzyBitapSearch(&textf, &patternf, 1)
		if err != nil || resf != 1 {
		t.Fatalf("error: %s with result %d", err, res)
	}
}
