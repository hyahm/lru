package main

import "lru"

func main() {
	lru.Init(10)
	lru.Add("adsf", "bbbbb")
	lru.Add("cccc", "111111")
	lru.Add("dddd", "2222222")

	lru.Add("eeeee", "33333")
	lru.Add("fffff", "4444444")
	lru.Add("gggg", "23423")
	lru.Add("hhhh", "645345")
	lru.Add("iiiii", "6789678")

	lru.Add("jjjjjjj", "123123")
	lru.Add("kkkkkk", "53536790")
	lru.Add("lllll", "0000")
	lru.Add("mmm", "226666")

	lru.OrderPrint()
}
