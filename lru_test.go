package lru

import (
	"fmt"
	"testing"
)

func Test_lru(t *testing.T) {
	Init(3)
	add("apple", 1)
	add("orange", 2)
	add("apple", 3)
	add("orange", 378)
	add("orange", 313)
	add("apple", 262)
	Remove("apple")
	Print()
	x := Keys()
	fmt.Println(x)
	fmt.Println(Get("orange"))
}
