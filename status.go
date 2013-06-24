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
	ClientError = errors.New(StatusClientError)
	NotFound = errors.New(StatusNotFound)
	ServerError = errors.New(StatusServerError)
	Error = errors.New(StatusError)
)
