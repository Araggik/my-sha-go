//go:build !solution

package lrucache

import "container/list"

type LRUCache struct {
	//Вместимость
	cap int
	l *list.List
	m map[int]int
	listMap map[int]*list.Element
}

func (cache LRUCache) Get(key int) (v int, ok bool){
	v, ok = cache.m[key]

	if ok {
		cache.moveKeyToBack(key)
	}

	return 
}

func (cache LRUCache) Set(key, value int) {
	_, ok := cache.m[key]

	if !ok && len(cache.m) == cache.cap {
		cache.removeOldKey()
	}

	cache.m[key] = value
	cache.moveKeyToBack(key)

	return 
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
	for k, v := range cache.m {
		cache.moveKeyToBack(k)

		check := f(k, v)

		if !check {
			return 
		}
	}
}

func (cache LRUCache) removeOldKey() {
	first := cache.l.Front()

	firstKey := first.Value

	delete(cache.m, firstKey)
	delete(cache.listMap, firstKey)
	cache.l.Remove(first)
}

func (cache LRUCache) moveKeyToBack(key int) {
	v, ok := cache.listMap[key]

	if ok {
		cache.l.Remove(v)
	}

	cache.listMap[key] = cache.l.PushBack(key);

	return 
}

func New(cap int) Cache {
	return LRUCache{
		cap: cap,
		l: list.New(),
		m: make(map[int]int, cap),
		listMap: make(map[int]*list.Element, cap),
	}
}
