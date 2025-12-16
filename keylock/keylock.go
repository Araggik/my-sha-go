//go:build !solution

package keylock

import (
	"container/list"
	"fmt"
	"strings"
)

type KeyLock struct {
	lockedKeys map[string]struct{}
	freeKeys   map[string]struct{}
	//Ключ в map - один из входных ключей при вызове LockKeys(), value - list *MissingData
	keyMap map[string]*list.List
	//Ключ в map - строка, являющаяся конкатинацией недостающих ключей для разблокировки по waitCh
	//value - list *LockData
	missingKeyMap map[string]*list.List
	//Мютекс для работы с lockedKeys
	lkMu chan struct{}
}

func New() *KeyLock {
	kl := &KeyLock{
		lockedKeys:    make(map[string]struct{}),
		freeKeys:      make(map[string]struct{}),
		keyMap:        make(map[string]*list.List),
		missingKeyMap: make(map[string]*list.List),
		lkMu:          make(chan struct{}, 1),
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
		missingKeys := list.New()

		forEachKey(keys, func(key string) bool {
			_, ok := l.lockedKeys[key]

			if ok {
				insertWithSort(missingKeys, key)
			}

			return false
		})

		waitCh := make(chan struct{})

		lockData := &LockData{keys, waitCh}

		missingData := &MissingData{missingKeys, lockData}

		//Пополнение keyMap
		forEachKey(keys, func(key string) bool {
			li, ok := l.keyMap[key]

			if !ok {
				li = list.New()
			}

			li.PushBack(missingData)

			l.keyMap[key] = li

			return false
		})

		//TODO: пополнение missingKeyMap

		<-l.lkMu

		//Ожидание из каналов
		select {
		case <-waitCh:
			canceled = false
			unlock = func() {
				l.unlockKeys(keys)
			}
		case <-cancel:
			canceled = true
			unlock = func() {}
		}

		//Если нет блокировки, то просто пополняем множество lockedKeys
	} else {
		l.addLockedKeys(keys)

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

	//TODO: доделать unlockKeys

	<-l.lkMu
}

func (l *KeyLock) addLockedKeys(keys []string) {
	//value - ключ для missingKeyMap
	missingDataSet := make(map[*MissingData]string)

	forEachKey(keys, func(key string) bool {
		//Добавляем в lockedKeys
		l.lockedKeys[key] = struct{}{}

		li, ok := l.keyMap[key]

		if ok {
			for el := li.Front(); el != nil; el = el.Next() {
				missingData := el.Value.(*MissingData)

				_, ok = missingDataSet[missingData]

				//Запоминаем ключ для missingKeyMap, чтобы потом удалить lockData
				//по этому ключу
				if !ok {
					keyForMissingMap := receiveKeyForMissingMap(missingData)

					missingDataSet[missingData] = keyForMissingMap
				}

				//Пополняем список недостающих ключей в MissingData
				insertWithSort(missingData.missingKeys, key)
			}
		}

		return false
	})

	//Корректируем missingKeyMap
	for md, oldKey := range missingDataSet {
		//Добавляем в missingKeyMap
		keyForMissingMap := receiveKeyForMissingMap(md)

		li, ok := l.missingKeyMap[keyForMissingMap]

		if !ok {
			li = list.New()
		}

		li.PushBack(md.lockData)

		l.missingKeyMap[keyForMissingMap] = li

		//Удаляем из missingKeyMap
		li = l.missingKeyMap[oldKey]

		for el := li.Front(); el != nil; el = el.Next() {
			ld := el.Value.(*LockData)

			if md.lockData == ld {
				li.Remove(el)

				break
			}
		}
	}
}

func receiveKeyForMissingMap(md *MissingData) string {
	var builder strings.Builder

	for el := md.missingKeys.Front(); el != nil; el = el.Next() {
		builder.WriteString(el.Value.(string))
	}

	return builder.String()
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

// Вставляем строку в отсортированный список
func insertWithSort(l *list.List, value string) {
	el := l.Front()

	//Если пустой список
	if el == nil {
		l.PushBack(value)

		return
	}

	for ; el != nil; el = el.Next() {
		val := el.Value.(string)

		if value < val {
			l.InsertAfter(value, el)

			break
		}
	}
}

type MissingData struct {
	//Отсортированные ключи
	missingKeys *list.List
	lockData    *LockData
}

// Идентифицирует вызов LockKeys
type LockData struct {
	keys   []string
	waitCh chan struct{}
}