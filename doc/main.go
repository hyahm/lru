package main

import (
	"fmt"
	"lru"
)

func main() {
	lru.Init(10)

	lru.Add("adsf", "bbbbb")
	lru.Add("cccc", "111111")

	lru.Add("adsf", "333")
	fmt.Println(lru.Len())
	lru.OrderPrint()
}
