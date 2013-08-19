package memcached

import (
	"io"
	"fmt"
)

type MemcachedResponse interface {
	WriteResponse(io.Writer)
}

type ItemResponse struct {
	Item *Item
}

func (r *ItemResponse) WriteResponse(writer io.Writer) {
	fmt.Fprintf(writer, StatusValue, r.Item.Key, r.Item.Flags, len(r.Item.Value))
	writer.Write(r.Item.Value)
	writer.Write(crlf)
}

type BulkResponse struct {
	Responses []MemcachedResponse
}

func (r *BulkResponse) WriteResponse(writer io.Writer) {
	for _, response := range r.Responses {
		if response != nil {
			response.WriteResponse(writer)
		}
	}
}

type ClientErrorResponse struct {
	Reason string
}

func (r *ClientErrorResponse) WriteResponse(writer io.Writer) {
	fmt.Fprintf(writer, StatusClientError, r.Reason)
}
