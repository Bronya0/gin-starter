package set

import (
	"fmt"
	"testing"
)

func TestSet(t *testing.T) {
	s1 := New(1, 2, 3)
	fmt.Println(s1.Contain(4))
}
