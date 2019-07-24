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
 > 添加 key和value
  ```
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
