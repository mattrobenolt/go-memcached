package memcached

import "errors"

const (
	StatusEnd = "END"
	StatusError = "ERROR"
	StatusServerError = "SERVER_ERROR %s"
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
)

var (
	ClientError = errors.New(StatusClientError)
	NotFound = errors.New(StatusNotFound)
)
