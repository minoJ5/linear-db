package fuzzy

import (
	"fmt"
	"testing"
)

func TestLevenstein(t *testing.T){
	fmt.Println(Distance("app models user.rb", "am uesr"))
}