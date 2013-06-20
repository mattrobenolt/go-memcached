package memcached

import (
	"net"
	"io"
	"bufio"
	"bytes"
	"strconv"
	"fmt"
)

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
	defer c.Close()
	for {
		err := c.handleRequest()
		if err != nil {
			if err == io.EOF {
				return
			}
			c.end(fmt.Sprintf(err.Error(), "Unsupported"))
		}
	}
}

func (c *conn) end(s string) {
	io.WriteString(c.rwc, s)
	c.rwc.Write(crlf)
	c.rwc.Flush()
}

func (c *conn) handleRequest() error {
	line, err := c.ReadLine()
	if err != nil || len(line) == 0 {
		return err
	}
	switch line[0] {
	case 'g':
		key := line[4:]
		getter, ok := c.server.Handler.(Getter)
		if !ok {
			return ClientError
		}
		item, err := getter.Get(key)
		if err != nil {
			c.end(StatusEnd)
		} else {
			fmt.Fprintf(c.rwc, StatusValue, item.Key, item.Flags, item.Length)
			c.rwc.Write(crlf)
			c.rwc.Write(item.Value)
			c.rwc.Write(crlf)
			c.end(StatusEnd)
		}
	case 's':
		setter, ok := c.server.Handler.(Setter)
		if !ok {
			return ClientError
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

		err = setter.Set(item)
		if err != nil {
			c.end(StatusNotStored)
		} else {
			c.end(StatusStored)
		}
	default:
		return ClientError
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
	item.Key = make([]byte, len(pieces[0]))
	copy(item.Key, pieces[0])

	// lol, no error handling here
	item.Flags, _ = strconv.Atoi(string(pieces[1]))
	exptime, _ := strconv.ParseInt(string(pieces[2]), 10, 64)
	item.SetExpires(exptime)
	item.Length, _ = strconv.Atoi(string(pieces[3]))
}
