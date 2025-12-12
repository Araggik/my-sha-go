//go:build !solution

package keylock

import (
	"container/list"
	"fmt"
)

type KeyLock struct {
	lockedKeys map[string]struct{}
	freeKeys   map[string]struct{}
	//Ключ в map - один из входных ключей при вызове LockKeys(), value - list *MissingData
	keyMap map[string]list.List
	//Ключ в map - строка, являющаяся конкатинацией недостающих ключей для разблокировки по waitCh
	//value - list *LockData
	missingKeyMap map[string]list.List
	//Мютекс для работы с lockedKeys
	lkMu chan struct{}
}

func New() *KeyLock {
	kl := &KeyLock{
		lockedKeys:    make(map[string]struct{}),
		freeKeys:      make(map[string]struct{}),
		keyMap:        make(map[string]list.List),
		missingKeyMap: make(map[string]list.List),
		lkMu:          make(chan struct{}, 1),
	}

	return kl
}

func (l *KeyLock) LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func()) {
	l.lkMu <- struct{}{}

	//Проверяем заблокирован ли вызов
	isLocked := false

	n := len(keys)

	for i := 0; i < n && !isLocked; i++ {
		key := keys[i]

		_, ok := l.lockedKeys[key]

		if ok {
			isLocked = true
		}
	}

	if isLocked {
		missingKeys := list.New()

		for i := 0; i < n; i++ {
			key := keys[i]

			_, ok := l.lockedKeys[key]

			if ok {
				//TODO: добавить missingKey в missingKeys, так чтобы сохранялась сортировка
			}
		}

		<-l.lkMu

		waitCh := make(chan struct{})

		lockData := LockData{keys, waitCh}

		//Если нет блокировки, то просто пополняем множество lockedKeys
	} else {
		for i := 0; i < n; i++ {
			key := keys[i]
			l.lockedKeys[key] = struct{}{}
		}

		<-l.lkMu

		canceled = false
		unlock = func() {
			l.unlockKeys(keys)
		}
	}
	return
}

func (l *KeyLock) unlockKeys(keys []string) {

}

type MissingData struct {
	//Отсортированные ключи
	missingKeys []string
	lockData    *LockData
}

// Идентифицирует вызов LockKeys
type LockData struct {
	keys   []string
	waitCh chan struct{}
}