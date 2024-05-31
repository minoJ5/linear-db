package fuzzy

import (
	"fmt"
	"testing"
)

var list []string = []string{
	"app models user.rb",
	"app models order.rb",
	"app models customer.rb",
}

var query string = "am user"

func TestFzy(t *testing.T) {
	for _, s := range list {
		fmt.Println(Positions(query, s))
		fmt.Println(HasMatch(query, s))
		fmt.Println(Score(query, s))
	}
}

func BenchmarkFzy(b *testing.B) {
	for i :=0; i < b.N; i++ {
		for _, s := range list {
			Positions(query, s)
			HasMatch(query, s)
			//score(query, s)
		}
	}

}