package cmap

import "sync"

type Iter struct {
	C chan Kv
}

type Kv struct {
	Key   interface{}
	Value interface{}
}

func (m *ConcurrentMap) Iterator() *Iter {
	iter := &Iter{
		C: make(chan Kv),
	}

	go func() {
		wg := sync.WaitGroup{}

		for _, s := range m.shard {
			wg.Add(1)

			go func(s *concurrentMapShared) {
				defer wg.Done()
				it := s.Iterator()
				for kv := range it.C {
					iter.C <-kv
				}
			}(s)
		}

		wg.Wait()
		close(iter.C)
	}()

	return iter
}

func (s *concurrentMapShared) Iterator() *Iter {
	return newIter(s)
}

func newIter(m *concurrentMapShared) *Iter {
	iter := &Iter{
		C: make(chan Kv),
	}

	go func() {
		m.RLock()
		defer m.RUnlock()

		for k, v := range m.items {
			iter.C <- Kv{
				Key: k,
				Value: v,
			}
		}

		close(iter.C)
	}()

	return iter
}
