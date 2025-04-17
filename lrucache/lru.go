//go:build !solution

package lrucache

import (
	"container/list"
)

type pairInt struct {
	key   int
	value int
}

type lruCache struct {
	cap  int
	list *list.List
	m    map[int]*list.Element
}

func (lc lruCache) pushTop(key, value int) {
	if ptr, exist := lc.m[key]; exist {
		lc.list.MoveToBack(ptr)
		ptr.Value = pairInt{key, value}
	} else if len(lc.m) < lc.cap {
		lc.m[key] = lc.list.PushBack(pairInt{key, value})
	} else if lc.cap > 0 {
		delete(lc.m, lc.list.Front().Value.(pairInt).key)
		lc.list.Remove(lc.list.Front())

		lc.m[key] = lc.list.PushBack(pairInt{key, value})
	}
}

func (lc lruCache) Get(key int) (val int, exist bool) {
	if ptr, ex := lc.m[key]; ex {
		val, exist = ptr.Value.(pairInt).value, true
		lc.pushTop(key, val)
		return val, exist
	}
	return val, false
}

func (lc lruCache) Set(key, val int) {
	lc.pushTop(key, val)
}

func (lc lruCache) Clear() {
	lc.list.Init()
	for k := range lc.m {
		delete(lc.m, k)
	}
}

func (lc lruCache) Range(f func(int, int) bool) {
	for ptr := lc.list.Front(); ptr != nil; ptr = ptr.Next() {
		tmp := ptr.Value.(pairInt)
		if !f(tmp.key, tmp.value) {
			return
		}
	}
}

func New(cap int) Cache {
	lc := lruCache{cap, list.New(), make(map[int]*list.Element)}
	var cache Cache
	cache = lc
	return cache
}
