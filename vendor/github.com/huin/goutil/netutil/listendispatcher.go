package netutil

import (
	"net"
)

type AcceptErr struct {
	L   net.Listener
	Err error
}

var _ error = AcceptErr{}

func (ae AcceptErr) Error() string {
	return ae.Err.Error()
}

func AcceptToChan(l net.Listener, connChan chan<- net.Conn, errChan chan<- AcceptErr) {
	for {
		conn, err := l.Accept()
		if err != nil {
			errChan <- AcceptErr{L: l, Err: err}
			return
		}
		connChan <- conn
	}
}

type ListenDispatcher interface {
	NewListenAddr(string)
}

type listenDispatcher struct {
	addrChan     chan string
	connDispatch func(net.Conn)
	listenErr    func(error)
	acceptErr    func(error)
}

// NewListenDispatcher creates a new ListenDispatcher, which will call
// connDispatch for every connection it accepts. listenErr and acceptErr may be
// nil, but if provided are called for errors when attempting to listen and
// accept respectively. A call to either indicates that the listener will stop
// until its NewListenAddr method is called.
func NewListenDispatcher(connDispatch func(net.Conn), listenErr, acceptErr func(error)) ListenDispatcher {
	if connDispatch == nil {
		panic("connDispatch must not be nil")
	}
	ld := &listenDispatcher{
		addrChan:     make(chan string),
		connDispatch: connDispatch,
		listenErr:    listenErr,
		acceptErr:    acceptErr,
	}
	go ld.run()
	return ld
}

func (ld *listenDispatcher) NewListenAddr(addr string) {
	ld.addrChan <- addr
}

func (ld *listenDispatcher) run() {
	var l net.Listener

	acceptErr := make(chan AcceptErr)
	newConn := make(chan net.Conn)

	for {
		select {
		case addr := <-ld.addrChan:
			if l != nil {
				l.Close()
			}
			if addr == "" {
				l = nil
			} else {
				var err error
				l, err = net.Listen("tcp", addr)
				if err != nil && ld.listenErr != nil {
					ld.listenErr(err)
				}
			}

		case conn := <-newConn:
			ld.connDispatch(conn)

		case err := <-acceptErr:
			if err.L == l {
				l.Close()
				l = nil
				if ld.acceptErr != nil {
					ld.acceptErr(err.Err)
				}
			}
		}
	}
}
