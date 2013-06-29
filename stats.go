package memcached

import (
	"fmt"
	"runtime"
	"os"
	"strconv"
	"time"
)

type Stats map[string] fmt.Stringer


type StaticStat struct {
	Value string
}

func (s *StaticStat) String() string {
	return s.Value
}

type TimerStat struct {
	Value int64
}

func (t *TimerStat) String() string {
	return strconv.Itoa(int(time.Now().Unix() - t.Value))
}

type FuncStat struct {
	Callable func() string
}

func (f *FuncStat) String() string {
	return f.Callable()
}

type CounterStat struct {
	Count int
	calculations chan int
}

func (c *CounterStat) Increment(num int) {
	c.calculations <- num
}

func (c *CounterStat) SetCount(num int) {
	c.Count = num
}

func (c *CounterStat) Decrement(num int) {
	c.calculations <- -num
}

func (c *CounterStat) String() string {
	return strconv.Itoa(c.Count)
}

func (c *CounterStat) work() {
	for num := range c.calculations {
		c.Count = c.Count + num
	}
}

func NewCounterStat() *CounterStat {
	c := &CounterStat{}
	go c.work()
	return c
}

func NewStats() Stats {
	s := make(Stats)
	s["pid"] = &StaticStat{strconv.Itoa(os.Getpid())}
	s["uptime"] = &TimerStat{time.Now().Unix()}
	s["time"] = &FuncStat{func() string { return strconv.Itoa(int(time.Now().Unix())) }}
	s["version"] = &StaticStat{VERSION}
	s["golang"] = &StaticStat{runtime.Version()}
	s["goroutines"] = &FuncStat{func() string { return strconv.Itoa(runtime.NumGoroutine()) }}
	s["cmd_get"] = NewCounterStat()
	s["cmd_set"] = NewCounterStat()
	s["get_hits"] = NewCounterStat()
	s["get_misses"] = NewCounterStat()
	s["curr_connections"] = NewCounterStat()
	s["total_connections"] = NewCounterStat()
	s["evictions"] = NewCounterStat()
	return s
}
