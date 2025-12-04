//go:build !solution

package keylock

import "container/heap"

type KeyLock struct {
	chMap map[string]LockDataHeap
	//Мютекс для работы с chMap
	muCh chan struct{}
	//Заблокированные ключи
	lockedKeys map[string]struct{}
}

func New() *KeyLock {
	return &KeyLock{
		chMap:      make(map[string]LockDataHeap),
		mapCh:      make(chan struct{}, 1),
		lockedKeys: make(map[string]struct{}),
	}
}

func (l *KeyLock) LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func()) {
	isLocked := false

	n := len(keys)

	l.muCh <- struct{}{}

	//Проверяем, заблокированы ли уже ключи
	for i := 0; i < n && !isLocked; i++ {
		_, ok := l.lockedKeys[keys[i]]

		if ok {
			isLocked := true
		}
	}

	if isLocked {
		waitCh := make(chan struct{})

		data := LockData{keys, waitCh}

		for _, key := range keys {
			h := l.chMap[key]

			heap.Push(h, data)

		}

		<-l.muCh

		<-waitCh

	} else {
		for _, key := range keys {
			l.lockedKeys[key] = struct{}{}
		}

		<-l.muCh
	}

}

type LockData struct {
	keys   []string
	waitCh chan struct{}
}

type LockDataHeap []LockData

func (h LockDataHeap) Len() int {
	return len(h)
}

func (h LockDataHeap) Less(i, j int) bool {
	return len(h[i].keys) < len(h[j].keys)
}

func (h LockDataHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *LockDataHeap) Push(x any) {
	*h = append(*h, x.(LockData))
}

func (h *LockDataHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}