package interview

import (
	"sync"
)

type UserAges struct {
	ages map[string]int
	mu   sync.Mutex
	rw   sync.RWMutex
}

func (ua *UserAges) Add(name string, age int) {
	ua.mu.Lock()
	defer ua.mu.Unlock()
	ua.ages[name] = age
}
func (ua *UserAges) Get(name string) int {
	if age, ok := ua.ages[name]; ok { // 并发时会爆异常，fatal error: concurrent map read and map write，应该该用读写锁
		return age
	}
	return -1
}
func (ua *UserAges) Add2(name string, age int) {
	ua.rw.Lock()
	defer ua.rw.Unlock()
	ua.ages[name] = age
}

func (ua *UserAges) Get2(name string) int {
	ua.rw.RLock()
	defer ua.rw.RUnlock()
	if age, ok := ua.ages[name]; ok {
		return age
	}
	return -1
}

type Once struct {
	m    sync.Mutex
	done uint32
}

func (o *Once) Do(f func()) {
	if o.done == 1 {
		return
	}
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		o.done = 1
		f()
	}
}
