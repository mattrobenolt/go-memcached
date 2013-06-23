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
}

func (c *CounterStat) Increment(num int) {
	c.Count = c.Count + num
}

func (c *CounterStat) SetCount(num int) {
	c.Count = num
}

func (c *CounterStat) Decrement(num int) {
	c.Count = c.Count - num
}

func (c *CounterStat) String() string {
	return strconv.Itoa(c.Count)
}

func NewStats() Stats {
	s := make(Stats)
	s["pid"] = &StaticStat{strconv.Itoa(os.Getpid())}
	s["uptime"] = &TimerStat{time.Now().Unix()}
	s["time"] = &FuncStat{func() string { return strconv.Itoa(int(time.Now().Unix())) }}
	s["version"] = &StaticStat{VERSION}
	s["golang"] = &StaticStat{runtime.Version()}
	s["goroutines"] = &FuncStat{func() string { return strconv.Itoa(runtime.NumGoroutine()) }}
	s["cmd_get"] = &CounterStat{}
	s["cmd_set"] = &CounterStat{}
	s["get_hits"] = &CounterStat{}
	s["get_misses"] = &CounterStat{}
	s["curr_connections"] = &CounterStat{}
	s["total_connections"] = &CounterStat{}
	s["evictions"] = &CounterStat{}
	return s
}