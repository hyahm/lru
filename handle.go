package lru

import (
	"fmt"
	"log"
	"sync"
)

const DEFAULTCOUNT = 100

//开始存一个值
func (l *List) Add(key interface{}, value interface{}) {
	if l.lru == nil {
		l = Init(DEFAULTCOUNT)
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	l.add(key, value)
}

// 获取值
func (l *List) Get(key interface{}) interface{} {
	if l.lru == nil {
		l = Init(DEFAULTCOUNT)
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	if value, ok := l.lru[key]; ok {
		return value.value
	}
	return nil
}

// 获取当前的keys, 没有值返回nil
func (l *List) Keys() []interface{} {
	if l.lru == nil {
		return nil
	}
	keys := make([]interface{}, 0)
	for k := range l.lru {
		keys = append(keys, k)
	}

	return keys
}

func (l *List) Next(key interface{}) interface{} {
	if l.lru == nil {
		l = Init(DEFAULTCOUNT)
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	if value, ok := l.lru[key]; ok {
		if value.next == nil {
			return nil
		}
		return value.next.key
	}
	return nil
}

func (l *List) Prev(key interface{}) interface{} {
	if l.lru == nil {
		l = Init(DEFAULTCOUNT)
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	if value, ok := l.lru[key]; ok {
		if value.prev == nil {
			return nil
		}
		return value.prev.key
	}
	return nil
}

func (l *List) Remove(key interface{}) {
	if l.lru == nil {
		return
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	this := l.lru[key]
	//如果是第一个元素
	if this == l.root {
		tmp := l.root.next
		tmp.prev = nil
		l.root = tmp
		delete(l.lru, key)
		return
	}
	//如果是最后一个
	if this == l.last {
		tmp := l.last.prev
		tmp.next = nil
		l.last = tmp
		delete(l.lru, key)
		return
	}
	// 不存在就直接返回
	if _, ok := l.lru[key]; !ok {
		return
	}
	// 更改上一个元素的下一个值

	l.lru[key].prev.next = l.lru[key].next
	//更新下一个元素的上一个值
	l.lru[key].next.prev = l.lru[key].prev
	//删除
	delete(l.lru, key)
	l.len--
}

func (l *List) OrderPrint() {
	if l == nil {
		return
	}
	l.lock.Lock()
	for li := l.root; li.next != nil; li = li.next {
		fmt.Println("key: ", li.key, "---- value: ", li.value, " ---- nextkey: ", li.next.key)
	}
	l.lock.Unlock()
}

func (l *List) Len() int {
	return l.len
}

func (l *List) Print() {
	if l.lru == nil {
		return
	}
	for k, v := range l.lru {
		log.Println("key: ", k, " ---- value: ", v.value)
	}
}

func (l *List) Resize(n int) {
	//如果缩小了缓存, 那么可能需要删除后面多余的索引
	l.count = n
	if n < l.count {
		for l.len > n {
			l.removeLast()
		}
	}
}

// 返回被删除的key, 如果没删除返回nil
func (l *List) add(key interface{}, value interface{}) interface{} {
	//先要判断是否存在这个key, 存在的话，就将元素移动最开始的位置,
	el := &element{
		prev: nil,
		next: nil,
	}
	if l.len == 0 {
		// 只有一个值， 那么就没有上一个元素和下一个元素
		el.value = value
		el.key = key
		//更新链表
		// 更新第一个元素
		l.root = el
		// 更新最后一个元素
		l.last = el
		// 更新长度
		l.len = 1
		// 更新lru
		l.lru[key] = el
	}
	if _, ok := l.lru[key]; ok {
		//如果是第一个元素的话, 什么也不用操作
		if l.root == l.lru[key] {
			l.root.value = value
			return nil
		} else {

			// 否则就插入到开头, 开头的元素后移
			l.moveToPrev(key, value)
		}

	} else {
		//如果不存在的话, 直接添加到开头
		//下一个元素是开头的元素
		el.next = l.root
		el.key = key
		el.value = value
		//更新第二个元素
		tmp := l.root
		tmp.prev = el

		//将开头的元素修改成新的元素
		l.root = el
		//更新lru
		l.lru[key] = l.root
		l.lru[tmp.key] = tmp

		//如果一开始只有一个元素, 更新最后一个元素的值
		if l.len == 1 {
			// 更新最后一个值
			l.last = l.root.next

		}
		l.len++
		//判断长度是否超过了缓存
		for l.len > l.count {
			//移除最后一个元素, 移除之前先更新最后一个元素
			l.removeLast()
		}
	}
	return nil
}

// 移除最后一个
func (l *List) RemoveLast() interface{} {
	return l.removeLast()
}

func (l *List) removeLast() interface{} {
	tmp := l.last.prev
	tmp.next = nil
	l.lru[tmp.key] = tmp
	removekey := l.last.key
	delete(l.lru, l.last.key)
	l.last = tmp
	l.len--
	return removekey
}

func (l *List) moveToPrev(key interface{}, value interface{}) {
	// 这里面的元素至少有2个, 否则进不来这里
	// 否则就插入到开头, 开头的元素后移
	//把当前位置元素的上一个元素的下一个元素指向本元素的下一个元素
	//el := &element{}

	if l.len == 2 {
		//如果是2个元素
		//也就是更换元素的值就好了
		//把第一个元素换到第二去
		lasttmp := l.root
		roottmp := l.last
		lasttmp.prev = roottmp
		lasttmp.next = nil
		roottmp.prev = nil
		roottmp.next = lasttmp
		roottmp.value = value
		l.root = roottmp
		l.last = lasttmp
		//更新lru
		l.lru[l.root.key] = l.root
		l.lru[l.last.key] = l.last

		//lru
		return
	}
	if l.len > 2 {
		if l.lru[key] == l.last {

			//如果这个元素是最后一个, 更新这个元素
			//如果这个值是最后一个的话, 还要更新倒数第二个元素
			l.last.prev.next = nil
			// 最后一个元素 是最后一个元素
			l.last = l.last.prev
			l.lru[l.last.key] = l.last
		}
		//如果不是, 更新这个元素 上一个和下一个元素的值
		l.lru[key].prev.next = l.lru[key].next
		l.lru[key].next.prev = l.lru[key].prev
		//抽出来这个值到开头
		l.lru[key].prev = nil
		l.lru[key].value = value
		l.lru[key].next = l.root
		// tmp 是第二个元素
		tmp := l.root
		l.root = l.lru[key]

		// 更新 第二个元素
		tmp.prev = l.root
		//更新第二个元素的Lru
		l.lru[tmp.key] = tmp

	}
}

func (l *List) FirstKey() interface{} {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.root.key
}

func (l *List) LastKey() interface{} {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.last.key
}

func (l *List) Clean(n int) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l = nil
	l.lru = nil
	l = &List{
		lru:   make(map[interface{}]*element, 0),
		len:   0,
		count: n,
		lock:  sync.RWMutex{},
		root:  &element{},
		last:  &element{},
	}
}

func (l *List) Exsit(key interface{}) bool {
	l.lock.Lock()
	defer l.lock.Unlock()
	if _, ok := l.lru[key]; ok {
		return true
	}
	return false
}
