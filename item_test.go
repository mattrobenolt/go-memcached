package memcached

import (
	"testing"
	"time"
)

const TEN_SECONDS = time.Duration(10)*time.Second

func TestIsExpiredZero(t *testing.T) {
	item := &Item{}
	if item.IsExpired() {
		t.Error("Zero time shouldn't be expired")
	}
}

func TestIsExpiredFuture(t *testing.T) {
	item := &Item{
		Expires: time.Now().Add(TEN_SECONDS),
	}
	if item.IsExpired() {
		t.Error("Future shouldn't be expired")
	}
}

func TestIsExiredPast(t *testing.T) {
	item := &Item{
		Expires: time.Now().Add(-TEN_SECONDS),
	}
	if !item.IsExpired() {
		t.Error("Past should be expired")
	}
}

func TestSetExpiresZero(t *testing.T) {
	item := &Item{}
	item.SetExpires(0)
	if !item.Expires.IsZero() {
		t.Error("Zero should have a Zero time:", item.Expires)
	}
	if item.Ttl != 0 {
		t.Error("Zero should have a Ttl of 0:", item.Ttl)
	}
}

func TestSetExpiresTypical(t *testing.T) {
	item := &Item{}
	item.SetExpires(10)
	if item.Expires.IsZero() {
		t.Error("Shouldn't have a Zero Expires time")
	}
	if !item.Expires.After(time.Now()) {
		t.Error("Expires should be in the future.")
	}
	if item.Expires.Sub(time.Now()) >= TEN_SECONDS || item.Expires.Sub(time.Now()) < time.Duration(9)*time.Second {
		t.Error("Expires should be > 9, and < 10 seconds")
	}
	if item.Ttl != 10 {
		t.Error("Ttl should be 10:", item.Ttl)
	}
}

func TestSetExpiresMaxExptime(t *testing.T) {
	item := &Item{}
	expires := time.Unix(int64(60*60*24*30) + 1, 0)  // 1 second greater than 30 days
	item.SetExpires(expires.Unix())
	if !item.Expires.Equal(expires) {
		t.Error("Expires should be 1970-01-30:", item.Expires)
	}
	if item.Ttl > -1370325162 {  // well, will always be smaller than this known point in time since it's based on Now()
		t.Error("Ttl should be really really really low:", item.Ttl)
	}
}
