package lru

import (
	"sync"
)

type element struct {
	// 上一个元素和下一个元素
	next, prev *element
	// The list to which this element belongs.
	//元素的key
	key interface{}
	// 这个元素的值
	value interface{}
}

// 数量太少, 因为边界问题难测试, 所以定义了最小长度
//const LESS = 5

//var Lru *list

type List struct {
	lru map[interface{}]*element //  这里存key 和 元素
	//保存第一个元素
	lock  sync.RWMutex
	root  *element // sentinel list element, only &root, root.prev, and root.next are used
	last  *element // 最后一个元素
	len   int      // 元素长度
	count int      // 缓存多少元素
}

func Init(n int) *List {
	// 内存足够的话, 可以设置很大, 所有计算都是O(1)
	if n <= 0 || n > 2<<32 {
		n = 2 << 10
	}

	return &List{
		lru:   make(map[interface{}]*element, 0),
		count: n,
		lock:  sync.RWMutex{},
		root:  &element{},
		last:  &element{},
	}
}
