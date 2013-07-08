package memcached

import (
	"time"
	"fmt"
)

// The maximum time to send from a client before
// the timestamp is considered an absolute unix
// timestamp.
const MAX_EXPTIME = 60*60*24*30  // 30 days

type Item struct {
	Key string
	Value []byte
	Flags, Ttl int
	Expires time.Time
}

// Check if an Item is expired based on it's Ttl.
// If an item has no Ttl set, it is considered to never
// be expired.
func (i *Item) IsExpired() bool {
	return !i.Expires.IsZero() && i.Expires.Before(time.Now())
}

// Set the Ttl and Expires based on the exptime send from
// a client. This follows standard memcached rules, and an
// exptime greater than 30 days is treated as an absolute
// unix timestamp.
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

// Initialize a new Item
func NewItem() *Item {
	return &Item{}
}
