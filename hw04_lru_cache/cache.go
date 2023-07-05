package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value any) bool
	Get(key Key) (any, bool)
	Clear()
}

type cacheItem struct {
	k Key
	v any
}

type lruCache struct {
	capacity int
	mu       sync.Mutex
	queue    List
	items    map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lc *lruCache) Set(key Key, value any) bool {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	if i, ok := lc.items[key]; ok {
		ci := i.Value.(*cacheItem)
		ci.v = value
		lc.queue.MoveToFront(i)
		return ok
	}

	li := lc.queue.PushFront(&cacheItem{key, value})
	if lc.queue.Len() > lc.capacity {
		b := lc.queue.Back()
		lc.queue.Remove(b)
		delete(lc.items, b.Value.(*cacheItem).k)
	}
	lc.items[key] = li

	return false
}

func (lc *lruCache) Get(key Key) (any, bool) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	if i, ok := lc.items[key]; ok {
		lc.queue.MoveToFront(i)
		return i.Value.(*cacheItem).v, ok
	}
	return nil, false
}

func (lc *lruCache) Clear() {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	lc.queue = NewList()
	lc.items = make(map[Key]*ListItem, lc.capacity)
}
