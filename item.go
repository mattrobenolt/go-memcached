package memcached

import (
	"time"
	"fmt"
)

const MAX_EXPTIME = 60*60*24*30  // 30 days

type Item struct {
	Key string
	Value []byte
	Flags, Ttl int
	Expires time.Time
}

func (i *Item) IsExpired() bool {
	return !i.Expires.IsZero() && i.Expires.Before(time.Now())
}

func (i *Item) SetExpires(exptime int64) {
	if exptime > MAX_EXPTIME {
		i.Expires = time.Unix(exptime, 0)
		i.Ttl = int(i.Expires.Sub(time.Now()).Seconds())
	} else if exptime > 0 {
		i.Ttl = int(exptime)
		i.Expires = time.Now().Add(time.Duration(exptime)*time.Second)
	}
}

func (i *Item) String() string {
	return fmt.Sprintf("<Item %s Flags:%d Length:%d Ttl:%d Expires:%s>", i.Key, i.Flags, len(i.Value), i.Ttl, i.Expires)
}

func NewItem() *Item {
	return &Item{}
}
