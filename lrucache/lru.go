//go:build !solution

package lrucache

import "container/list"

type LRUCache struct {
	//Вместимость
	cap     int
	l       *list.List
	m       map[int]int
	listMap map[int]*list.Element
}

func (cache LRUCache) Get(key int) (v int, ok bool) {
	v, ok = cache.m[key]

	if ok {
		cache.moveKeyToBack(key)
	}

	return
}

func (cache LRUCache) Set(key, value int) {
	if cache.cap > 0 {
		_, ok := cache.m[key]

		if !ok && len(cache.m) == cache.cap {
			cache.removeOldKey()
		}

		cache.m[key] = value
		cache.moveKeyToBack(key)
	}
}

func (cache LRUCache) Clear() {
	for k := range cache.m {
		delete(cache.m, k)
	}

	for k := range cache.listMap {
		delete(cache.listMap, k)
	}

	cache.l.Init()
}

func (cache LRUCache) Range(f func(key, value int) bool) {
	elem := cache.l.Front()

	for elem != nil {
		k := elem.Value.(int)

		v := cache.m[k]

		if !f(k, v) {
			return
		}

		elem = elem.Next()
	}
}

func (cache LRUCache) removeOldKey() {
	first := cache.l.Front()

	firstKey := first.Value.(int)

	delete(cache.m, firstKey)
	delete(cache.listMap, firstKey)
	cache.l.Remove(first)
}

func (cache LRUCache) moveKeyToBack(key int) {
	v, ok := cache.listMap[key]

	if ok {
		cache.l.Remove(v)
	}

	cache.listMap[key] = cache.l.PushBack(key)
}

func New(cap int) Cache {
	return LRUCache{
		cap:     cap,
		l:       list.New(),
		m:       make(map[int]int, cap),
		listMap: make(map[int]*list.Element, cap),
	}
}
