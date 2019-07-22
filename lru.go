package lru

import "fmt"

type element struct {
	// 上一个元素和下一个元素
	next, prev *element
	// The list to which this element belongs.
	//list *List
	// 这个元素的值
	Value interface{}
}
//  这里存key 和 元素
var Lru map[interface{}]*element

type List struct {
	//保存第一个元素
	root *element // sentinel list element, only &root, root.prev, and root.next are used
	last *element  // 最后一个元素
	len  int     // current list length excluding (this) sentinel element
}

func New() *List {
	Lru = make(map[interface{}]*element, 0)
	return &List{
		root: &element{},
	}
}
//开始存一个值
func (l *List) Add(key interface{}, value interface{}) *List {
	return l.add(key,value)
}

func (l *List) Print() {
	for li := l.root; li != nil; li = li.next {
		fmt.Println(111111)
		fmt.Println(li.Value)
	}
}

func (l *List) add(key interface{}, value interface{}) *List  {
	//先要判断是否存在这个key, 存在的话，就将元素移动最开始的位置,
	el := &element{
		prev: nil,
		next: nil,
	}
	if l.len == 0 {
		// 只有一个值， 那么就没有上一个元素和下一个元素
		el.Value = value
		//更新链表
		// 更新第一个元素
		l.root = el
		// 更新最后一个元素
		l.last = el
		// 更新长度
		l.len = 1
		// 更新lru
		Lru[key] = el
		return l
	}
	if _, ok := Lru[key]; ok {
		//如果是第一个元素的话, 什么也不用操作
		if l.root == Lru[key] {
			return l
		} else {
			// 否则就插入到开头, 开头的元素后移
			l.moveToPrev(key , value )
		}

	} else {
		//如果不存在的话, 直接添加到开头
		//下一个元素是开头的元素
		el.next = l.root
		el.Value = value
		// 更新 旧的root的值
		l.root.prev = el
		//将开头的元素修改成新的元素
		l.root = el
		Lru[key] = l.root
		//如果一开始只有一个元素, 更新最后一个元素的值
		if l.len == 1 {
			// 更新最后一个值
			l.last = l.root.next

		}
		l.len ++
	}
	//判断长度
	return l
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
		l.root, l.last = l.last , l.root
	}
}