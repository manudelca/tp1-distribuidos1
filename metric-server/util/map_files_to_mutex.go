package util

import "sync"

type MapStringToMutex struct {
	stringToMutex map[string]*sync.Mutex
	mapMutex      sync.RWMutex
}

func (m *MapStringToMutex) Add(key string, mut *sync.Mutex) *sync.Mutex {
	m.mapMutex.Lock()
	savedMutex, ok := m.stringToMutex[key]
	if !ok {
		m.stringToMutex[key] = mut
		savedMutex = mut
	}
	m.mapMutex.Unlock()
	return savedMutex
}

func (m *MapStringToMutex) Get(key string) (*sync.Mutex, bool) {
	m.mapMutex.RLock()
	mut, ok := m.stringToMutex[key]
	m.mapMutex.RUnlock()
	return mut, ok
}
