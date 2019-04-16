package netutil

import (
	"io"
	"net"
	"time"
)

// ConnWrapper forwards Read and Write method calls to Reader and Writer
// respectively. All other methods calls go to UnderlyingConn directly. If
// ReadCloser or WriteCloser are not nil, then a call to Close on them will be
// made prior to calling Close on UnderlyingConn. Fields Reader, Writer and
// UnderlyingConn must *not* be nil.
type ConnWrapper struct {
	io.Reader
	io.Writer
	UnderlyingConn net.Conn
	ReadCloser     io.Closer
	WriteCloser    io.Closer
}

var _ net.Conn = new(ConnWrapper)

func (c *ConnWrapper) Close() error {
	var multiErr MultiError
	if c.ReadCloser != nil {
		multiErr.RecordError(c.ReadCloser.Close())
	}
	if c.WriteCloser != nil {
		multiErr.RecordError(c.WriteCloser.Close())
	}
	multiErr.RecordError(c.UnderlyingConn.Close())
	return multiErr.ToError()
}

func (c *ConnWrapper) LocalAddr() net.Addr                { return c.UnderlyingConn.LocalAddr() }
func (c *ConnWrapper) RemoteAddr() net.Addr               { return c.UnderlyingConn.RemoteAddr() }
func (c *ConnWrapper) SetDeadline(t time.Time) error      { return c.UnderlyingConn.SetDeadline(t) }
func (c *ConnWrapper) SetReadDeadline(t time.Time) error  { return c.UnderlyingConn.SetReadDeadline(t) }
func (c *ConnWrapper) SetWriteDeadline(t time.Time) error { return c.UnderlyingConn.SetWriteDeadline(t) }
