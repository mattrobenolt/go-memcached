package memcached

import (
	"testing"
	"time"
)

func TestStaticStat(t *testing.T) {
	stat := &StaticStat{"lol"}
	if stat.String() != "lol" {
		t.Error("Should be 'lol'", stat.String())
	}
}

func TestFuncStat(t *testing.T) {
	stat := &FuncStat{func() string { return "lol" }}
	if stat.String() != "lol" {
		t.Error("Should be 'lol'", stat.String())
	}
}

func TestCounterStat(t *testing.T) {
	stat := NewCounterStat()
	var i int
	for i = 0; i < 10; i++ {
		stat.Increment(1)
	}
	time.Sleep(1) // Force the intenal goroutine to catch up with the counts
	if stat.String() != "10" {
		t.Error("Should be '10'", stat.String())
	}
	for i = 0; i < 10; i++ {
		stat.Decrement(1)
	}
	time.Sleep(1)
	if stat.String() != "0" {
		t.Error("Should be '0'", stat.String())
	}
	stat.SetCount(100)
	if stat.String() != "100" {
		t.Error("Should be '100'", stat.String())
	}
}

func TestTimerStat(t *testing.T) {
	stat := NewTimerStat()
	time.Sleep(time.Duration(1) * time.Second)
	if stat.String() != "1" {
		t.Error("Should be '1'", stat.String())
	}
}
