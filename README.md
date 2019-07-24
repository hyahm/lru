# lru
go 语言通用lru包, 使用异常简单
### 安装
```
go get github.com/hyahm/lru
```
### 使用

// 在使用前, 先要初始化缓存个数, 内存足够大的话, 可以设置很大, 不会因为长度影响效率, 
// 超过设定值会自动删除末尾的值, 如果存在的话会自动更新此值到开头, 更新值
 > 初始化(初始化完成后, 可以在任何地方调用方法)
  ```
  lru.Init(10)
  ```
 > 添加 key和value, 以下为doc中的例子
  ```
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

```
> 顺序打印
```
lru.OrderPrint()
```
> 无序打印
```
lru.Print()
```
> 删除key
```
lru.Remove(key)
```
> 获取缓存长度
```
lru.Len()
```
> 根据key获取值
```
lru.Get(key)
```
> 根据key获取上一个key
```
lru.Get(key)
```
> 根据key获取下一个key
```
lru.Get(key)
```
基本上这些就能满足需求
