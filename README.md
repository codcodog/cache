一个简易的内存缓存系统
=======================

### 使用示例
```golang
c := NewCache()
c.Set("hello", "world", 5*time.Second)

h, found := c.Get("hello")
if found {
    fmt.Println(h)
}
```

详细参考：cache_test.go
