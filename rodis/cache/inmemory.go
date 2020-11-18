package cache

import (
	"fmt"
	"sync"
)

type Inmemory struct {
	buckets map[string][]byte
	mu      sync.RWMutex
	Stat
}

func NewInmemory() Cache {
	return &Inmemory{
		buckets: make(map[string][]byte),
		mu:      sync.RWMutex{},
		Stat:    Stat{},
	}
}

func (this *Inmemory) Get(key string) ([]byte, error) {
	this.mu.RLock()
	defer this.mu.RUnlock()

	v, ok := this.buckets[key]
	if ok {
		return v, nil
	}
	return nil, fmt.Errorf("not found value,error key:%s", key)
}

func (this *Inmemory) Set(key string, value []byte) error {
	this.mu.Lock()
	defer this.mu.Unlock()

	if v, ok := this.buckets[key]; ok {
		this.del(key, v)
	}
	this.add(key, value)
	this.buckets[key] = value
	return nil
}

func (this *Inmemory) Del(key string) error {
	this.mu.Lock()
	defer this.mu.Unlock()

	if v, ok := this.buckets[key]; ok {
		this.del(key, v)
		delete(this.buckets, key)
	}
	return nil
}

func (this *Inmemory) GetStat() Stat {
	return this.Stat
}
