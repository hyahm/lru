package lru

import (
	"fmt"
	"testing"
)

func Test_lru(t *testing.T) {
	l := Init(3)
	l.Add("apple", 1)
	l.Add("orange", 2)
	l.Add("apple", 3)
	l.Add("orange", 378)
	l.Add("orange", 313)
	l.Add("apple", 262)
	l.Remove("apple")
	l.Print()
	x := l.Keys()
	fmt.Println(x)
	fmt.Println(l.Get("orange"))
}
