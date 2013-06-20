package memcached

import "time"

const MAX_EXPTIME = 60*60*24*30  // 30 days

type Item struct {
	Key string
	Flags int64
	Length int64
	Expires time.Time
	Value string
}

func (i *Item) IsExpired() bool {
	return !i.Expires.IsZero() && i.Expires.Before(time.Now())
}

func (i *Item) SetExpires(exptime int64) {
	if exptime > MAX_EXPTIME {
		i.Expires = time.Unix(int64(exptime), 0)
	} else if exptime > 0 {
		i.Expires = time.Now().Add(time.Duration(exptime)*time.Second)
	}
}
