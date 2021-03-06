package zoo

import "sync"

type SafeMap struct {
	sync.RWMutex
	m map[string]interface{}
}

func NewSafeMap() *SafeMap {
	return &SafeMap{m: map[string]interface{}{}}
}

func (m *SafeMap) Set(k string, v interface{}) {
	m.Lock()
	defer m.Unlock()

	m.m[k] = v
}

func (m *SafeMap) Get(k string) interface{} {
	m.RLock()
	defer m.RUnlock()

	if v, ok := m.m[k]; ok {
		return v
	}

	return nil
}
