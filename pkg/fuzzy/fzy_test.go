package fuzzy

import (
	"fmt"
	"testing"
	"time"
)

var list []string = []string{
	"app models user1.rb",
	"app models order1.rb",
	"app models customer1.rb",
	"app models user2.rb",
	"app models order2.rb",
	"app models customer2.rb",
	"app models user3.rb",
	"app models order3.rb",
	"app models customer3.rb",
	"app models user4.rb",
	"app models order4.rb",
	"app models customer4.rb",
	"app models user5.rb",
	"app models order5.rb",
	"app models customer5.rb",
	"app models user6.rb",
	"app models order6.rb",
	"app models customer6.rb",
	"app models user7.rb",
	"app models order7.rb",
	"app models customer7.rb",
	"app models user8.rb",
	"app models order8.rb",
	"app models customer8.rb",
	"app models user9.rb",
	"app models order9.rb",
	"app models customer9.rb",
	"app models user10.rb",
	"app models order10.rb",
	"app models customer10.rb",
	"app models user11.rb",
	"app models order11.rb",
	"app models customer11.rb",
	"app models user12.rb",
	"app models order12.rb",
	"app models customer12.rb",
	"app models user13.rb",
	"app models order13.rb",
	"app models customer13.rb",
	"app models user14.rb",
	"app models order14.rb",
	"app models customer14.rb",
	"app models user15.rb",
	"app models order15.rb",
	"app models customer15.rb",
	"app models user16.rb",
	"app models order16.rb",
	"app models customer16.rb",
	"app models user17.rb",
	"app models order17.rb",
	"app models customer17.rb",
	"app models user18.rb",
	"app models order18.rb",
	"app models customer18.rb",
	"app models user19.rb",
	"app models order19.rb",
	"app models customer19.rb",
	"app models user20.rb",
	"app models order20.rb",
	"app models customer20.rb",
	"app models user21.rb",
	"app models order21.rb",
	"app models customer21.rb",
	"app models user22.rb",
	"app models order22.rb",
	"app models customer22.rb",
	"app models user23.rb",
	"app models order23.rb",
	"app models customer23.rb",
	"app models user24.rb",
	"app models order24.rb",
	"app models customer24.rb",
	"app models user25.rb",
	"app models order25.rb",
	"app models customer25.rb",
	"app models user26.rb",
	"app models order26.rb",
	"app models customer26.rb",
	"app models user27.rb",
	"app models order27.rb",
	"app models customer27.rb",
	"app models user28.rb",
	"app models order28.rb",
	"app models customer28.rb",
	"app models user29.rb",
	"app models order29.rb",
	"app models customer29.rb",
	"app models user30.rb",
	"app models order30.rb",
	"app models customer30.rb",
	"app models user31.rb",
	"app models order31.rb",
	"app models customer31.rb",
	"app models user32.rb",
	"app models order32.rb",
	"app models customer32.rb",
	"app models user33.rb",
	"app models order33.rb",
	"app models customer33.rb",
	"app models order34.rb",
	"app models customer34.rb",
	// Continue adding elements until you reach 100
}

var query string = "am user"

func TestFzy(t *testing.T) {
	start := time.Now()
	for _, s := range list {
		//fmt.Println(Positions(query, s))
		fmt.Println(HasMatch(query, s))
		fmt.Println(Score(query, s))
	}
	fmt.Println("Took: ", time.Since(start), " With Len: ", len(list))
}

func BenchmarkFzy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, s := range list {
			//Positions(query, s)
			//HasMatch(query, s)
			Score(query, s)
		}
	}
}
func TestFzyc(t *testing.T) {
	fmt.Println(Positions("hi", "hello iam"))
	FuzzyC()
}

func TestFyzCustom(t *testing.T) {
	n := "am user"
	h := "app models order34.rb"
	fmt.Println(MatchPositions(n, h))
}

func TestFyzCustomLoop(t *testing.T) {
	for _, s := range list {
		fmt.Println(MatchPositions(query, s))
	}
}
func BenchmarkFzyc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, s := range list {
			MatchPositions(query, s)
		}
	}
}

func BenchmarkFzygo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Positions("hi", "hello iam")
	}
}
