package lru

import (
	"fmt"
	"testing"
)

func Test_Add(t *testing.T) {
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

func BenchmarkWrite(b *testing.B) {
	//times :=
	l := Init(10000000)
	for i := 0;i<b.N;i++ {
		l.Add("apple", i)
	}
}

func BenchmarkRead(b *testing.B) {
	//times :=
	l := Init(10000000)
	for i := 0;i<10000000;i++ {
		l.Add("apple", i)
	}
	l.Print()
}

func ExampleList_Add() {
	l := Init(3)
	l.Add("apple", 1)
	l.Add("orange", 2)
	l.Add("apple", 3)
	l.Add("orange", 378)
	l.Add("orange", 313)
	l.Add("apple", 262)
	x := l.Keys()
	fmt.Println(x)
	fmt.Println(l.Get("orange"))
}

func ExampleList_Get() {
	l := Init(3)
	l.Add("apple", 1)
	l.Add("orange", 2)
	l.Add("apple", 3)
	l.Add("orange", 378)
	l.Add("orange", 313)
	l.Add("apple", 262)
	x := l.Keys()
	fmt.Println(x)
	fmt.Println(l.Get("orange"))
}

func ExampleList_Remove() {
	l := Init(3)
	l.Add("orange", 313)
	l.Add("apple", 262)
	l.Remove("apple")
	l.Print()
}