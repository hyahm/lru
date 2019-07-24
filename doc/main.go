package main

import (
	"fmt"
	"lru"
)

type el struct {
	Id int
	Name string
}

func main() {
	lru.Init(10)

	lru.Add("adsf", "bbbbb")
	lru.Add("cccc", "111111")
	e := &el{
		Id: 1,
		Name: "68",
	}
	lru.Add("adsf", e)
	fmt.Println(lru.Len())
	lru.OrderPrint()
}
