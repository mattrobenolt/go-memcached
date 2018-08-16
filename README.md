# Memcached library for Go

## What even is this?
Modeled similarly to the stdlib `net/http` package, `memcached` gives you a simple interface to building your own memcached protocol compatible applications.

## Install
```
$ go get github.com/mattrobenolt/go-memcached
```

## Interfaces
Implement as little or as much as you'd like.
```go
type Getter interface {
	RequestHandler
	Get(string) (*Item, error)
}

type Setter interface {
	RequestHandler
	Set(*Item) error
}

type Deleter interface {
	RequestHandler
	Delete(string) error
}
```

## Hello World
```go
package main

import (
	memcached "github.com/mattrobenolt/go-memcached"
)

type Cache map[string]*memcached.Item

func (c Cache) Get(key string) memcached.MemcachedResponse {
	return &memcached.ItemResponse{&memcached.Item{Key: key, Value: []byte("world")}}
}

func main() {
	server := memcached.NewServer(":11211", make(Cache))
	server.ListenAndServe()
}
```

## Examples
 * [Simple Memcached](examples/memcached.go)  *Don't actually use this*

## Documentation
 * [http://godoc.org/github.com/mattrobenolt/go-memcached](http://godoc.org/github.com/mattrobenolt/go-memcached)
