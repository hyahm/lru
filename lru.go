package lru

import (
	"fmt"
	"log"
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
const LESS = 5

var Lru *list

type list struct {
	lru map[interface{}]*element //  这里存key 和 元素
	//保存第一个元素
	lock sync.RWMutex
	root *element // sentinel list element, only &root, root.prev, and root.next are used
	last *element  // 最后一个元素
	len  uint64     // current list length excluding (this) sentinel element
	count uint64  // 缓存多少元素
}

func Init(n uint64) {
	// 内存足够的话, 可以设置很大, 长度不会影响效率
	if n < LESS {
		e := fmt.Sprintf("cache count must more than %d" , LESS)
		panic(e)
	}

	Lru = &list{
		lru: make(map[interface{}]*element, 0),
		count: n,
		lock: sync.RWMutex{},
		root: &element{},
		last: &element{},
	}
}
//开始存一个值
func Add(key interface{}, value interface{}) {
	if Lru == nil {
		panic("must init first")
	}
	Lru.lock.Lock()
	defer Lru.lock.Unlock()
	add(key,value)
	log.Println("last: ",Lru.last.key)
}

// 获取值
func Get(key interface{}) interface{} {
	if Lru == nil {
		panic("must init first")
	}
	Lru.lock.Lock()
	defer Lru.lock.Unlock()
	if value, ok := Lru.lru[key]; ok {
		return value.value
	}
	return nil
}

// 获取当前的keys, 没有值返回nil
func Keys() []interface{} {
	if Lru == nil {
		return nil
	}
	keys := make([]interface{}, 0)
	for k, _ := range Lru.lru {
		keys = append(keys, k)
	}

	return keys
}

func Next(key interface{}) interface{} {
	if Lru == nil {
		panic("must init first")
	}
	Lru.lock.Lock()
	defer Lru.lock.Unlock()
	if value, ok := Lru.lru[key]; ok {
		if value.next == nil {
			return  nil
		}
		return value.next.key
	}
	return nil
}

func Prev(key interface{}) interface{} {
	if Lru == nil {
		panic("must init first")
	}
	Lru.lock.Lock()
	defer Lru.lock.Unlock()
	if value, ok := Lru.lru[key]; ok {
		if value.prev == nil {
			return  nil
		}
		return value.prev.key
	}
	return nil
}

func Remove(key interface{}) {
	if Lru == nil {
		panic("must init first")
	}
	if Lru.len < LESS {
		e := fmt.Sprintf("cache count less than %d, can't remove" , LESS)
		log.Println(e)
		return
	}
	Lru.lock.Lock()
	defer Lru.lock.Unlock()
	this := Lru.lru[key]
	//如果是第一个元素
	if this == Lru.root {
		tmp := Lru.root.next
		tmp.prev = nil
		Lru.root = tmp
		delete(Lru.lru, key)
		return
	}
	//如果是最后一个
	if this == Lru.last {
		tmp := Lru.last.prev
		tmp.next = nil
		Lru.last = tmp
		delete(Lru.lru, key)
		return
	}
	// 中间的话, 直接删除,
	// 更改上一个元素的下一个值
	Lru.lru[key].prev.next = Lru.lru[key].next
	//更新下一个元素的上一个值
	Lru.lru[key].next.prev = Lru.lru[key].prev
	//删除
	delete(Lru.lru, key)
	Lru.len--
}


func OrderPrint() {
	if Lru == nil {
		panic("must init first")
	}
	Lru.lock.Lock()
	for li := Lru.root; li.next != nil; li = li.next {
		fmt.Println("key: ",li.key, "---- value: ", li.value, " ---- nextkey: ",li.next.key)
	}
	fmt.Println("key: ",Lru.last.key, "---- value: ", Lru.last.value, " ---- nextkey: ",nil)
	Lru.lock.Unlock()
}

func Len() uint64 {
	return Lru.len
}

func Print() {
	if Lru == nil {
		panic("must init first")
	}
	for k, v := range Lru.lru{
		log.Println("key: ", k, " ---- value: ",v.value)
	}
}

func Resize(n uint64) {
	//如果缩小了缓存, 那么可能需要删除后面多余的索引
	Lru.count = n
	if n < Lru.count {
		for Lru.len > n {
			removeLast()
		}
	}
}

// 返回被删除的key, 如果没删除返回nil
func add(key interface{}, value interface{}) interface{} {
	//先要判断是否存在这个key, 存在的话，就将元素移动最开始的位置,
	el := &element{
		prev: nil,
		next: nil,
	}
	if Lru.len == 0 {
		// 只有一个值， 那么就没有上一个元素和下一个元素
		el.value = value
		el.key = key
		//更新链表
		// 更新第一个元素
		Lru.root = el
		// 更新最后一个元素
		Lru.last = el
		// 更新长度
		Lru.len = 1
		// 更新lru
		Lru.lru[key] = el
	}
	if _, ok := Lru.lru[key]; ok {
		//如果是第一个元素的话, 什么也不用操作
		if Lru.root == Lru.lru[key] {
			Lru.root.value = value
			return nil
		} else {

			// 否则就插入到开头, 开头的元素后移
			moveToPrev(key , value )
		}

	} else {
		//如果不存在的话, 直接添加到开头
		//下一个元素是开头的元素
		el.next = Lru.root
		el.key = key
		el.value = value
		//更新第二个元素
		tmp := Lru.root
		tmp.prev = el

		//将开头的元素修改成新的元素
		Lru.root = el
		//更新lru
		Lru.lru[key] = Lru.root
		Lru.lru[tmp.key] = tmp

		//如果一开始只有一个元素, 更新最后一个元素的值
		if Lru.len == 1 {
			// 更新最后一个值
			Lru.last = Lru.root.next

		}
		Lru.len++
		//判断长度是否超过了缓存
		if uint64(Lru.len) > Lru.count {
			//移除最后一个元素, 移除之前先更新最后一个元素
			removeLast()
		}
	}
	return nil
}

func removeLast() interface{} {
	tmp := Lru.last.prev
	tmp.next = nil
	Lru.lru[tmp.key] = tmp
	removekey := Lru.last.key
	delete(Lru.lru, Lru.last.key)
	Lru.last = tmp
	Lru.len--
	return removekey
}

func moveToPrev(key interface{}, value interface{}) {
	// 这里面的元素至少有2个, 否则进不来这里
	// 否则就插入到开头, 开头的元素后移
	//把当前位置元素的上一个元素的下一个元素指向本元素的下一个元素
	//el := &element{}

	if Lru.len == 2 {
		//如果是2个元素
		//也就是更换元素的值就好了
		//把第一个元素换到第二去
		lasttmp := Lru.root
		roottmp := Lru.last
		lasttmp.prev = roottmp
		lasttmp.next = nil
		roottmp.prev = nil
		roottmp.next = lasttmp
		roottmp.value = value
		Lru.root = roottmp
		Lru.last = lasttmp
		//更新lru
		Lru.lru[Lru.root.key]= Lru.root
		Lru.lru[Lru.last.key]= Lru.last

		//lru
		return
	}
	if Lru.len > 2 {
		if  Lru.lru[key] == Lru.last {

			//如果这个元素是最后一个, 更新这个元素
			//如果这个值是最后一个的话, 还要更新倒数第二个元素
			Lru.last.prev.next = nil
			// 最后一个元素 是最后一个元素
			Lru.last = Lru.last.prev
			Lru.lru[Lru.last.key] = Lru.last
		}
		//如果不是, 更新这个元素 上一个和下一个元素的值
		Lru.lru[key].prev.next = Lru.lru[key].next
		Lru.lru[key].next.prev = Lru.lru[key].prev
		//抽出来这个值到开头
		Lru.lru[key].prev = nil
		Lru.lru[key].value = value
		Lru.lru[key].next = Lru.root
		// tmp 是第二个元素
		tmp := Lru.root
		Lru.root = Lru.lru[key]

		// 更新 第二个元素
		tmp.prev = Lru.root
		//更新第二个元素的Lru
		Lru.lru[tmp.key] = tmp

	}
}


func FirstKey() interface{} {
	Lru.lock.Lock()
	defer Lru.lock.Unlock()
	return Lru.root.key
}

func LastKey() interface{} {
	Lru.lock.Lock()
	defer Lru.lock.Unlock()
	return Lru.last.key
}

func Clean(n uint64) {
	Lru.lock.Lock()
	defer Lru.lock.Unlock()
	Lru.lru = nil
	Lru = nil
	Lru = &list{
		lru: make(map[interface{}]*element, 0),
		len: 0,
		count: n,
		lock: sync.RWMutex{},
		root: &element{},
		last: &element{},
	}
}

func Exsit(key interface{}) bool {
	Lru.lock.Lock()
	defer Lru.lock.Unlock()
	if _, ok := Lru.lru[key]; ok {
		return true
	}
	return false
}