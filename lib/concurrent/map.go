package concurrent

import (
	"sort"
	"sync"
)

type Map struct {
	lock sync.Mutex
	data map[string]interface{}
}

func (m *Map) Init() {
	m.lock.Lock()
	m.data = map[string]interface{}{}
	m.lock.Unlock()
}

func (m *Map) Set(k string, v interface{}) {
	m.lock.Lock()
	m.data[k] = v
	m.lock.Unlock()
}

func (m *Map) Get(k string) interface{} {
	m.lock.Lock()
	defer m.lock.Unlock()
	if v, ok := m.data[k]; ok {
		return v
	}
	return nil
}

func (m *Map) Has(k string) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	if _, ok := m.data[k]; ok {
		return true
	}
	return false
}

func (m *Map) Count() int {
	m.lock.Lock()
	defer m.lock.Unlock()
	return len(m.data)
}

func (m *Map) Data() map[string]interface{} {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.data
}

func (m *Map) Keys() []string {
	m.lock.Lock()
	defer m.lock.Unlock()
	var list []string
	for k, _ := range m.data {
		list = append(list, k)
	}
	sort.Strings(list)
	return list
}
