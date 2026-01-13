//go:build !solution

package keylock

import (
	"container/list"
)

type KeyLock struct {
	lockedKeys map[string]struct{}
	//Ключ в map - один из входных ключей при вызове LockKeys(), value - list *LockData
	keyMap map[string]*list.List
	//Мютекс для работы с lockedKeys
	lkMu chan struct{}
}

func New() *KeyLock {
	kl := &KeyLock{
		lockedKeys: make(map[string]struct{}),
		keyMap:     make(map[string]*list.List),
		lkMu:       make(chan struct{}, 1),
	}

	return kl
}

func (l *KeyLock) LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func()) {
	l.lkMu <- struct{}{}

	//Проверяем заблокирован ли вызов
	isLocked := false

	forEachKey(keys, func(key string) bool {
		_, ok := l.lockedKeys[key]

		if ok {
			isLocked = true
		}

		return isLocked
	})

	if isLocked {
		waitCh := make(chan struct{})

		ld := &LockData{keys, waitCh}

		//Пополнение keyMap
		forEachKey(keys, func(key string) bool {
			li, ok := l.keyMap[key]

			if !ok {
				li = list.New()
			}

			li.PushBack(ld)

			l.keyMap[key] = li

			return false
		})

		<-l.lkMu

		select {
		case <-cancel:
			canceled = true
			unlock = func() {}
		case <-waitCh:
			canceled = false
			unlock = func() {
				l.unlockKeys(keys)
			}
		}
	} else {
		//Пополняем lockedKeys
		forEachKey(keys, func(key string) bool {
			l.lockedKeys[key] = struct{}{}

			return false
		})

		<-l.lkMu
		canceled = false
		unlock = func() {
			l.unlockKeys(keys)
		}
	}
	return
}

func (l *KeyLock) unlockKeys(keys []string) {
	l.lkMu <- struct{}{}

	//Удаляем из lockedKeys
	forEachKey(keys, func(key string) bool {
		delete(l.lockedKeys, key)

		return false
	})

	ldQueue := list.New()

	//Для быстрой проверки, что LockData уже есть в ldQueue
	ldMap := make(map[*LockData]struct{})

	forEachKey(keys, func(key string) bool {
		if li, ok := l.keyMap[key]; ok {
			for el := li.Front(); el != nil; el = el.Next() {
				val := el.Value.(*LockData)

				if _, ok := ldMap[val]; !ok {
					//Проверяем, свободны ли ключи в LockData
					isLocked := false

					forEachKey(val.keys, func(k string) bool {
						_, ok := l.lockedKeys[k]

						if ok {
							isLocked = true
						}

						return isLocked
					})

					if !isLocked {
						//Добавление в ldMap
						ldMap[val] = struct{}{}

						//Добавление в ldQueue
						length := len(val.keys)

						isLast := true

						for elem := ldQueue.Front(); elem != nil; elem = elem.Next() {
							v := elem.Value.(*LockData)

							n := len(v.keys)

							if length >= n {
								li.InsertBefore(val, elem)

								isLast = false

								break
							}
						}

						if isLast {
							li.PushBack(val)
						}
					}
				}

			}
		}

		return false
	})

	for el := ldQueue.Front(); el != nil; el = el.Next() {
		val := el.Value.(*LockData)

		isLocked := false

		forEachKey(val.keys, func(k string) bool {
			_, ok := l.lockedKeys[k]

			if ok {
				isLocked = true
			}

			return isLocked
		})

		if !isLocked {
			//Разблокирование потока
			val.waitCh <- struct{}{}

			//Блокировка ключей
			forEachKey(val.keys, func(key string) bool {
				l.lockedKeys[key] = struct{}{}

				return false
			})

			l.removeLockData(val)
		}
	}

	<-l.lkMu
}

// Удаление LockData из keyMap
func (l *KeyLock) removeLockData(ld *LockData) {
	keys := ld.keys

	forEachKey(keys, func(key string) bool {
		li := l.keyMap[key]

		for el := li.Front(); el != nil; el = el.Next() {
			val := el.Value.(*LockData)

			if val == ld {
				li.Remove(el)

				break
			}
		}

		return false
	})
}

func forEachKey(keys []string, fn func(key string) bool) {
	n := len(keys)

	for i := 0; i < n; i++ {
		key := keys[i]

		res := fn(key)

		if res {
			break
		}
	}
}

// Идентифицирует вызов LockKeys
type LockData struct {
	keys   []string
	waitCh chan struct{}
}