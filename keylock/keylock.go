//go:build !solution

package keylock

import "container/heap"

type KeyLock struct {
	chMap map[string]*LockDataHeap
	//Мютекс для работы с chMap
	muCh chan struct{}
	//Заблокированные ключи
	lockedKeys map[string]struct{}
}

func New() *KeyLock {
	kl := &KeyLock{
		chMap:      make(map[string]*LockDataHeap),
		muCh:       make(chan struct{}, 1),
		lockedKeys: make(map[string]struct{}),
	}

	return kl
}

func (l *KeyLock) LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func()) {
	isLocked := false

	n := len(keys)

	l.muCh <- struct{}{}

	//Проверяем, заблокированы ли уже ключи
	for i := 0; i < n && !isLocked; i++ {
		_, ok := l.lockedKeys[keys[i]]

		if ok {
			isLocked = true
		}
	}

	if isLocked {
		waitCh := make(chan struct{})

		data := LockData{keys: keys, waitCh: waitCh}

		for _, key := range keys {
			h, ok := l.chMap[key]

			if !ok {
				h = &LockDataHeap{}

				l.chMap[key] = h
			}

			heap.Push(h, data)
		}

		<-l.muCh

		select {
		case <-waitCh:
			unlock = func() {
				l.unlockKeys(keys)
			}
		case <-cancel:
			canceled = true
			unlock = func() {}
		}
		close(waitCh)

	} else {
		for _, key := range keys {
			l.lockedKeys[key] = struct{}{}
		}

		<-l.muCh

		unlock = func() {
			l.unlockKeys(keys)
		}
	}

	return
}

func (l *KeyLock) unlockKeys(keys []string) {
	l.muCh <- struct{}{}

	for _, key := range keys {
		delete(l.lockedKeys, key)
	}

	//Нужен для поиска каналов, через которые возможно можно разблокировать горутины
	unlockHeap := &LockDataHeap{}

	for _, key := range keys {
		h, ok := l.chMap[key]

		if ok {
			var data *LockData

			//TODO: подумать над этим циклом, вроде бы нужно брать
			//не первый элемент из кучи, а элемент с меньшим набором ключей, которые все свободны

			for data == nil && h.Len() > 0 {
				d := heap.Pop(h).(LockData)

				select {
				//Если канал закрыт, то берем следующий элемент из кучи
				case <-d.waitCh:
				default:
					data = &d
				}
			}

			if data != nil {
				data.key = key
				heap.Push(unlockHeap, *data)
			}
		}
	}

	//TODO: берем элементы из unlockHeap и посылаем сигнал в канал, если ключи свободны.
	//Если ключи не свободны, то нужно вернуть элемент в l.chMap

	<-l.muCh
}

type LockData struct {
	keys   []string
	waitCh chan struct{}
	//Используется только в unlockKeys()
	key any
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