package lru

import "container/list"

// Cache is a LRU cache. It is not safe for concurrent access.
type Cache struct {
	maxBytes int64
	nBytes   int64
	ll       *list.List
	cache    map[string]*list.Element
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value Value)
}

// New is the Constructor of Cache
func New(maxBytes int64, OnEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: OnEvicted,
	}
}

// Get look ups a keys' value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}

	return
}

// RemoveOldest removes the oldest item
func (c *Cache) RemoveOldest() {
	if ele := c.ll.Back(); ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())

		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}

}

// Add adds a value to the cache.
func (c *Cache) Add(key string, value Value) {

	if ele, ok := c.cache[key]; ok { //1.假如key已存在则更新
		//移动最新修改的元素到队首
		c.ll.MoveToFront(ele)
		//重新计算当前缓存大小
		kv := ele.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		//更新值
		kv.value = value
	} else { //2.否则key不存在则新建
		//向队首保存新的值
		ele := c.ll.PushFront(&entry{key: key, value: value})
		//更新映射表
		c.cache[key] = ele
		//重新计算大小
		c.nBytes += int64(len(key)) + int64(value.Len())
	}

	//3.判断当前缓存是否已过大，过大需要移除旧的缓存
	for c.maxBytes != 0 && c.nBytes > c.maxBytes {
		c.RemoveOldest()
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}

type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}
