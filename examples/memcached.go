package main

import (
	"flag"
	"fmt"
	memcached "github.com/mattrobenolt/go-memcached"
	"log"
	"runtime"
)

var (
	listen  = flag.String("l", "", "Interface to listen on. Default to all addresses.")
	port    = flag.Int("p", 11211, "TCP port number to listen on (default: 11211)")
	threads = flag.Int("t", runtime.NumCPU(), fmt.Sprintf("number of threads to use (default: %d)", runtime.NumCPU()))
)

type Cache map[string]*memcached.Item

func (c Cache) Get(key string) memcached.MemcachedResponse {
	if item, ok := c[key]; ok {
		if item.IsExpired() {
			delete(c, key)
		} else {
			return &memcached.ItemResponse{item}
		}
	}
	return nil
}

func (c Cache) Set(item *memcached.Item) memcached.MemcachedResponse {
	c[item.Key] = item
	return nil
}

func (c Cache) Delete(key string) memcached.MemcachedResponse {
	delete(c, key)
	return nil
}

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(*threads)

	address := fmt.Sprintf("%s:%d", *listen, *port)
	server := memcached.NewServer(address, make(Cache))
	log.Fatal(server.ListenAndServe())
}
