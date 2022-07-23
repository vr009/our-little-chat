package delivery

import (
	"fmt"
	"testing"
)

func TestA(t *testing.T) {
	a := 100
	b := &a
	fmt.Println(a, *b)
	a = 120
	fmt.Println(a, *b)
}
