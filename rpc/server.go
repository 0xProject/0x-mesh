// +build !js

package rpc

import (
	"context"
	"net"
	"net/http"
	"strings"
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
// stop listening. Listen blocks until there is an error or the given context is
// canceled.
func (s *Server) Listen(ctx context.Context) error {
	s.mut.Lock()

	rpcService := &rpcService{
		rpcHandler: s.rpcHandler,
	}
	s.rpcServer = rpc.NewServer()
	if err := s.rpcServer.RegisterName("mesh", rpcService); err != nil {
		log.WithField("error", err.Error()).Fatal("could not register RPC service")
	}
	listener, err := net.Listen("tcp4", s.addr)
	if err != nil {
		s.mut.Unlock()
		log.WithField("error", err.Error()).Fatal("could not start listener")
	}
	s.listener = listener
	s.mut.Unlock()

	// Close the server when the context is canceled.
	go func() {
		<-ctx.Done()
		s.rpcServer.Stop()
		_ = s.listener.Close()
	}()

	if err := http.Serve(s.listener, s.rpcServer.WebsocketHandler([]string{"*"})); err != nil {
		// HACK(albrow): http.Serve doesn't accept a context. This means that
		// everytime we close the context for our rpc.Server, we see a "use of
		// closed network connection" error.
		if isClosedNetworkConnectionErr(err) {
			// Check whether the context is canceled in order to determine whether we
			// are in the process of tearing down the server.
			select {
			case <-ctx.Done():
				// If we are tearing down the server, this is okay and we don't need to
				// return the error.
				return nil
			default:
				// If we are not tearing down the server, the error is not expected, and
				// we should return it.
				return err
			}
		}
		return err
	}
	return nil
}

func isClosedNetworkConnectionErr(err error) bool {
	if opErr, ok := err.(*net.OpError); ok {
		if strings.Contains(opErr.Error(), "use of closed network connection") {
			return true
		}
	}
	return false
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
