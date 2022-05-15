package cache

import (
	"strconv"
	"strings"
	"sync"
	"time"
)

type Cache interface {
	// size 支持以下参数：1KB，100KB，1MB，1GB 等
	SetMaxMemory(size string) bool
	// 设置一个缓存项，并且在 expire 时间之后过期
	Set(key string, val interface{}, expire time.Duration)
	// 获取一个项
	Get(key string) (interface{}, bool)
	// 删除一个项
	Del(key string) bool
	// 检测一个项是否存在
	Exists(key string) bool
	// 清空所有项
	Flush() bool
	// 返回 key 的总数
	Keys() int64
}

const (
	KB int64 = 1024
	MB int64 = 1024 * 1024
	GB int64 = 1024 * 1024 * 1024
)

type cache struct {
	mu        sync.RWMutex
	items     map[string]Item
	maxMemory int64
}

func NewCache() *cache {
	c := &cache{
		items: make(map[string]Item),
	}
	go autoRecycling(c)

	return c
}

type Item struct {
	Val    interface{}
	Expire int64
}

// 判断是否已过期
func (i Item) Expired() bool {
	return i.Expire < time.Now().Unix()
}

// 设置内存大小
func (c *cache) SetMaxMemory(size string) bool {
	if size == "" {
		return false
	}

	length := len(size)
	numStr, unitStr := size[0:length-2], size[length-2:length]

	num, err := strconv.ParseInt(numStr, 10, 64)
	if err != nil {
		return false
	}
	unit := strings.ToUpper(unitStr)

	var max int64
	switch unit {
	case "KB":
		max = num * KB
	case "MB":
		max = num * MB
	case "GB":
		max = num * GB
	default:
		return false
	}
	c.maxMemory = max

	return true
}

// 是否超出内存限制
func (c *cache) isOutOfMemeory() bool {
	n := len(c.items)
	if n == 0 {
		return false
	}

	// @TODO 计算内存大小

	return false
}

// 设置缓存项
func (c *cache) Set(key string, val interface{}, expire time.Duration) {
	// @TODO 超出内存限制，淘汰机制
	if c.isOutOfMemeory() {
		return
	}

	e := time.Now().Add(expire).Unix()
	item := Item{
		Val:    val,
		Expire: e,
	}

	c.mu.Lock()
	c.items[key] = item
	c.mu.Unlock()
}

// 获取缓存项
func (c *cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if !ok {
		return nil, false
	}

	// 若过期，主动回收
	if item.Expired() {
		delete(c.items, key)
		return nil, false
	}

	return item.Val, true
}

// 删除缓存项
func (c *cache) Del(key string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.items[key]; ok {
		delete(c.items, key)
		return true
	}

	return false
}

// 检测项是否存在
func (c *cache) Exists(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if _, ok := c.items[key]; ok {
		return true
	}

	return false
}

// 清空所有项
func (c *cache) Flush() bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]Item)
	return true
}

// 返回 key 的总数
func (c *cache) Keys() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return int64(len(c.items))
}

// 每5分钟自动回收过期的缓存项
func autoRecycling(c *cache) {
	ticker := time.NewTicker(5 * time.Minute)
	for {
		select {
		case <-ticker.C:
			c.deleteExpired()
		}
	}
}

// 删除过期的缓存项
func (c *cache) deleteExpired() {
	c.mu.Lock()
	for key, item := range c.items {
		if item.Expired() {
			delete(c.items, key)
		}
	}
	c.mu.Unlock()
}
