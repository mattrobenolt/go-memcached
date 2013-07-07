package main

import (
	"flag"
	"fmt"
	memcached "github.com/mattrobenolt/go-memcached"
	"log"
)

var (
	listen = flag.String("l", "", "Interface to listen on. Default to all addresses.")
	port   = flag.Int("p", 11211, "TCP port number to listen on (default: 11211)")
)

type Cache map[string]*memcached.Item

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
	flag.Parse()
	address := fmt.Sprintf("%s:%d", *listen, *port)
	server := memcached.NewServer(address, make(Cache))
	log.Fatal(server.ListenAndServe())
}
