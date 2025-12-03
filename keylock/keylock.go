//go:build !solution

package keylock

import (
	"sort"
	"strings"
)

type KeyLock struct {
	//Ключ в map - конкатинация необходимых ключей для лока
	chMap map[string][](chan struct{})
	//Мютекс для работы с chMap
	mapCh chan struct{}
}

func New() *KeyLock {
	return &KeyLock{
		chMap: make(map[string][](chan struct{})),
		mapCh: make(chan struct{}, 1),
	}
}

func (l *KeyLock) LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func()) {
	sort.Strings(keys)

	mapKey := strings.Join(keys, " ")

	mapCh <- struct{}{}

	waitCh := l.receiveWaitCh(mapKey)

	<-mapCh
}

func (l *KeyLock) receiveWaitCh(mapKey string) chan struct{} {
	chArr, ok := l.chMap[mapKey]

	if !ok {
		chArr = make([](chan struct{}))
	}

	//TODO: создание канала для LockKeys и добавление его в map, если он нужен

	l.chMap[mapKey] = chArr
}
