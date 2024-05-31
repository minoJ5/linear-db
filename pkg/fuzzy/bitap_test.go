package fuzzy

import (
	"fmt"
	"testing"
)
const (
	text = "aklelpode dfeeww"
	pattern = "eew"
	textf = "app models user.rb"
	patternf = "uesr"
)
func TestBitap(t *testing.T) {

	res, err := BitapSearch(text, pattern)
	if err != nil || res != 12 {
		t.Fatalf("error: %s with result %d", err, res)
	}
	fmt.Println(res)
	resf, _ := FuzzyBitapSearch(textf, patternf, 2)
	fmt.Println(resf)
	// if err != nil || resf != 1 {
	// 	t.Fatalf("error: %s with result %d", err, res)
	// }
}

func BenchmarkBitap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FuzzyBitapSearch(textf, patternf, 2)
		//strings.Contains(text, pattern)
	}
}