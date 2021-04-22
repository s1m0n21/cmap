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
}
