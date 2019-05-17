// +build !js

package ws

import (
	"net"
	"net/http"
	"sync"

	"github.com/ethereum/go-ethereum/rpc"
	log "github.com/sirupsen/logrus"
)

// Server is a JSON RPC 2.0 server implementation over WebSockets. It accepts
// requests from a client for adding orders to the 0x Mesh network.
type Server struct {
	mut          sync.Mutex
	addr         string
	listenerAddr net.Addr
	rpcHandler   RPCHandler
	listener     net.Listener
	rpcServer    *rpc.Server
}

// NewServer creates and returns a new server which will listen for new
// connections on the given addr and use the rpcHandler to handle incoming
// requests.
func NewServer(addr string, rpcHandler RPCHandler) (*Server, error) {
	return &Server{
		addr:       addr,
		rpcHandler: rpcHandler,
	}, nil
}

// Listen causes the server to listen for new connections. You can call Close to
// stop listening. Listen blocks until there is an error.
func (s *Server) Listen() error {
	s.mut.Lock()

	rpcService := &rpcService{
		rpcHandler: s.rpcHandler,
	}
	s.rpcServer = rpc.NewServer()
	if err := s.rpcServer.RegisterName("mesh", rpcService); err != nil {
		log.WithField("err", err.Error()).Fatal("could not register RPC service")
	}
	listener, err := net.Listen("tcp4", s.addr)
	if err != nil {
		s.mut.Unlock()
		log.WithField("err", err.Error()).Fatal("could not start listener")
	}
	s.listener = listener
	s.mut.Unlock()

	return http.Serve(s.listener, s.rpcServer.WebsocketHandler([]string{"*"}))
}

// Addr returns the address the server is listening on or nil if it has not yet
// started listening.
func (s *Server) Addr() net.Addr {
	s.mut.Lock()
	defer s.mut.Unlock()
	if s.listener == nil {
		return nil
	}
	return s.listener.Addr()
}

// Close closes the listener and stops it from accepting new connections or
// responding to any new requests.
func (s *Server) Close() error {
	s.rpcServer.Stop()
	return s.listener.Close()
}
