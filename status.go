package memcached

import "errors"

const (
	StatusEnd = "END"
	StatusError = "ERROR"
	StatusServerError = "SERVER_ERROR"
	StatusClientError = "CLIENT_ERROR %s"
	StatusStored = "STORED"
	StatusNotStored = "NOT_STORED"
	StatusExists = "EXISTS"
	StatusNotFound = "NOT_FOUND"
	StatusDeleted = "DELETED"
	StatusTouched = "TOUCHED"
	StatusOK = "OK"
	StatusVersion = "VERSION %s"
	StatusValue = "VALUE %s %d %d"
	StatusStat = "STAT %s %s"
)

var (
	// An error caused by an invalid command from the client
	ClientError = errors.New(StatusClientError)
	// The key was not found
	NotFound = errors.New(StatusNotFound)
	// An error occured servicing this request
	ServerError = errors.New(StatusServerError)
	// Generic error
	Error = errors.New(StatusError)
)
