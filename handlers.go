package memcached

type RequestHandler interface{}

type Getter interface {
	RequestHandler
	Get([]byte) (*Item, error)
}

type Setter interface {
	RequestHandler
	Set(*Item) error
}
