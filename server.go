package memcached

import (
	"net"
	"io"
	"bufio"
	"bytes"
	"strconv"
	"fmt"
)

const VERSION = "0.0.0"

var (
	crlf = []byte("\r\n")
)

type conn struct {
	server *Server
	conn net.Conn
	rwc *bufio.ReadWriter
}

type Server struct {
	Addr string
	Handler RequestHandler
	Stats Stats
}

func (s *Server) newConn(rwc net.Conn) (c *conn, err error) {
	c = new(conn)
	c.server = s
	c.conn = rwc
	c.rwc = bufio.NewReadWriter(bufio.NewReaderSize(rwc, 1048576), bufio.NewWriter(rwc))
	return c, nil
}

func (s *Server) ListenAndServe() error {
	addr := s.Addr
	if addr == "" {
		addr = ":11211"
	}
	l, e := net.Listen("tcp", addr)
	if e != nil {
		return e
	}
	return s.Serve(l)
}

func (s *Server) Serve(l net.Listener) error {
	defer l.Close()
	for {
		rw, e := l.Accept()
		if e != nil {
			return e
		}
		c, err := s.newConn(rw)
		if err != nil {
			continue
		}
		go c.serve()
	}
}

func (c *conn) serve() {
	defer func() {
		c.server.Stats["curr_connections"].(*CounterStat).Decrement(1)
		c.Close()
	}()
	c.server.Stats["total_connections"].(*CounterStat).Increment(1)
	c.server.Stats["curr_connections"].(*CounterStat).Increment(1)
	for {
		err := c.handleRequest()
		if err != nil {
			if err == io.EOF {
				return
			}
			c.end(err.Error())
		}
	}
}

func (c *conn) end(s string) {
	c.rwc.WriteString(s)
	c.rwc.Write(crlf)
	c.rwc.Flush()
}

func (c *conn) handleRequest() error {
	line, err := c.ReadLine()
	if err != nil || len(line) == 0 {
		return io.EOF
	}
	switch line[0] {
	case 'g':
		key := string(line[4:]) // get
		getter, ok := c.server.Handler.(Getter)
		if !ok {
			return Error
		}
		c.server.Stats["cmd_get"].(*CounterStat).Increment(1)
		item, err := getter.Get(key)
		if err != nil {
			c.server.Stats["get_misses"].(*CounterStat).Increment(1)
			c.end(StatusEnd)
		} else {
			c.server.Stats["get_hits"].(*CounterStat).Increment(1)
			fmt.Fprintf(c.rwc, StatusValue, item.Key, item.Flags, len(item.Value))
			c.rwc.Write(crlf)
			c.rwc.Write(item.Value)
			c.rwc.Write(crlf)
			c.end(StatusEnd)
		}
	case 's':
		switch line[1] {
		case 'e':
			setter, ok := c.server.Handler.(Setter)
			if !ok {
				return Error
			}
			item := &Item{}
			parseStorageLine(line, item)
			value, err := c.ReadLine()
			if err != nil {
				return ClientError
			}

			// Copy the value into the *Item
			item.Value = make([]byte, len(value))
			copy(item.Value, value)

			c.server.Stats["cmd_set"].(*CounterStat).Increment(1)
			err = setter.Set(item)
			if err != nil {
				c.end(err.Error())
			} else {
				c.end(StatusStored)
			}
		case 't':
			for key, value := range c.server.Stats {
				fmt.Fprintf(c.rwc, StatusStat, key, value)
				c.rwc.Write(crlf)
			}
			c.end(StatusEnd)
		}
	case 'd':
		key := string(line[7:]) // delete
		deleter, ok := c.server.Handler.(Deleter)
		if !ok {
			return Error
		}
		err := deleter.Delete(key)
		if err != nil {
			c.end(StatusNotFound)
		} else {
			c.end(StatusDeleted)
		}
	default:
		return Error
	}
	return nil
}

func (c *conn) Close() {
	c.conn.Close()
}

func (c *conn) ReadLine() (line []byte, err error) {
	line, _, err = c.rwc.ReadLine()
	return
}

func ListenAndServe(addr string) error {
	s := &Server{
		Addr: addr,
	}
	return s.ListenAndServe()
}

func parseStorageLine(line []byte, item *Item) {
	pieces := bytes.Fields(line[4:])  // Skip the actual "set "
	item.Key = string(pieces[0])

	// lol, no error handling here
	item.Flags, _ = strconv.Atoi(string(pieces[1]))
	exptime, _ := strconv.ParseInt(string(pieces[2]), 10, 64)
	item.SetExpires(exptime)
}

func NewServer(listen string, handler RequestHandler) *Server {
	return &Server{listen, handler, NewStats()}
}
