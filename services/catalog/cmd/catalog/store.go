package main

import "sync"

type productStore struct {
	mu sync.RWMutex
	list []Product
}

func storeNew(seed []Product) *productStore {
	s := &productStore{}
	if seed != nil {
		s.list = append(s.list, seed...)
	}
	return s
}

func (s *productStore) List() []Product {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]Product, len(s.list))
	copy(out, s.list)
	return out
}
