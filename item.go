package memcached

import "time"

const MAX_EXPTIME = 60*60*24*30  // 30 days

type Item struct {
	Key, Value []byte
	Length, Flags, Ttl int
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

func (i *Item) SetValue(value []byte) {
	i.Value = value
	i.Length = len(value)
}
