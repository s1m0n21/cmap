package cmap

import (
	"testing"
)

type key struct {
	Content int
}

type value struct {
	Content int
}

func TestConcurrentMap(t *testing.T) {
	cm := New(DefaultShard)

	k := &key{Content: 1}
	v := &value{Content: 2}

	err := cm.Set(k, v)
	if err != nil {
		panic(err)
	}

	r, ok, err := cm.Get(k)
	if err != nil {
		panic(err)
	}

	t.Logf("ok: %v", ok)
	t.Logf("r: %+v", r.(*value))

	dr, err := cm.Del(k)
	if err != nil {
		panic(err)
	}

	t.Logf("dr: %v", dr)

	r, ok, err = cm.Get(k)
	if err != nil {
		panic(err)
	}

	t.Logf("ok: %v", ok)
}

func TestIterator(t *testing.T) {
	m := New(DefaultShard)

	for i := 0; i < 1000; i ++ {
		if err := m.Set(i, i*2); err != nil {
			panic(err)
		}
	}

	iter := m.Iterator()
	for kv := range iter.C {
		t.Logf("%v = %v", kv.Key.(int), kv.Value.(int))
	}
}
