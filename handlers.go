package memcached

type RequestHandler interface{}

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
