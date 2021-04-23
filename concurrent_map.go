package cmap

import (
	"bytes"
	"encoding/gob"
	"hash/fnv"
	"sync"
)

const DefaultShard uint64 = 32

type ConcurrentMap struct {
	shard []*concurrentMapShared
	count uint64
}

type concurrentMapShared struct {
	items map[interface{}]interface{}
	sync.RWMutex
}

func (m *ConcurrentMap) getShard(key interface{}) (*concurrentMapShared, error) {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(key); err != nil {
		return nil, err
	}

	h := fnv.New64()
	if _, err := h.Write(buf.Bytes()); err != nil {
		return nil, err
	}
	return m.shard[h.Sum64() % m.count], nil
}

func (m *ConcurrentMap) Set(key, value interface{}) error {
	shard, err := m.getShard(key)
	if err != nil {
		return err
	}

	shard.Lock()
	defer shard.Unlock()

	shard.items[key] = value

	return nil
}

func (m *ConcurrentMap) Get(key interface{}) (interface{}, bool, error) {
	shard, err := m.getShard(key)
	if err != nil {
		return nil, false, err
	}

	shard.RLock()
	defer shard.RUnlock()

	value, ok := shard.items[key]

	return value, ok, nil
}

func (m *ConcurrentMap) Has(key interface{}) (bool, error) {
	shard, err := m.getShard(key)
	if err != nil {
		return false, err
	}

	shard.RLock()
	defer shard.RUnlock()

	_, has := shard.items[key]

	return has, nil
}

func (m *ConcurrentMap) Del(key interface{}) (bool, error) {
	shard, err := m.getShard(key)
	if err != nil {
		return false, err
	}

	shard.Lock()
	defer shard.Unlock()

	delete(shard.items, key)

	return true, nil
}

func New(shard uint64) *ConcurrentMap {
	var items = make([]*concurrentMapShared, 0)

	var i uint64 = 0
	for ; i < shard; i++ {
		items = append(items, &concurrentMapShared{items: make(map[interface{}]interface{})})
	}

	return &ConcurrentMap{
		shard: items,
		count: shard,
	}
}
