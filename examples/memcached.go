package main

import (
	"log"
	memcached "github.com/mattrobenolt/go-memcached"
)

type Cache map[string] *memcached.Item

func (c Cache) Get(key string) (item *memcached.Item, err error) {
	if item, ok := c[key]; ok {
		if item.IsExpired() {
			delete(c, key)
		} else {
			return item, nil
		}
	}
	return nil, memcached.NotFound
}

func (c Cache) Set(item *memcached.Item) error {
	c[item.Key] = item
	return nil
}

func (c Cache) Delete(key string) error {
	delete(c, key)
	return nil
}

func main() {
	server := memcached.NewServer(":11211", make(Cache))
	log.Fatal(server.ListenAndServe())
}
