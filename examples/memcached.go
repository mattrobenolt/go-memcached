package main

import (
	"log"
	memcached "github.com/mattrobenolt/go-memcached"
)

type Store map[string] *memcached.Item

type Memcached struct {
	store Store
}

func (m *Memcached) Get(key string) (item *memcached.Item, err error) {
	if item, ok := m.store[key]; ok {
		if item.IsExpired() {
			m.delete(key)
		} else {
			return item, nil
		}
	}
	return nil, memcached.NotFound
}

func (m *Memcached) Set(item *memcached.Item) error {
	m.store[item.Key] = item
	return nil
}

func (m *Memcached) Delete(key string) error {
	m.delete(key)
	return nil
}

func (m *Memcached) delete(key string) {
	delete(m.store, key)
}

func main() {
	cache := &Memcached{
		store: make(Store),
	}
	server := memcached.NewServer(":11211", cache)
	log.Fatal(server.ListenAndServe())
}
