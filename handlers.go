package memcached

type RequestHandler interface{}

// A Getter is an object who responds to a simple
// "get" command.
type Getter interface {
	RequestHandler
	Get(string) (*Item, error)
}

// A Setter is an object who response to a simple
// "set" command.
type Setter interface {
	RequestHandler
	Set(*Item) error
}

// A Delter is an object who responds to a simple
// "delete" command.
type Deleter interface {
	RequestHandler
	Delete(string) error
}
