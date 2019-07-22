package main

import "lru"

func main() {
	l := lru.New()
	l.Add("adsf", "bbbbb")
	l.Print()
}
