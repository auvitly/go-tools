package labhttp

import (
	"crypto/tls"
	"fmt"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"
)

type Server struct {
	mu       sync.Mutex
	server   *http.Server
	router   *mux.Router
	handlers map[string]map[*testing.T]func(w http.ResponseWriter, r *http.Request)
}

// NewHTTPServer a method for creating a test server for HTTP clients.
func NewHTTPServer() *Server {
	addr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:0")
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}

	defer func() {
		err = conn.Close()
		if err != nil {
			panic(err)
		}
	}()

	var (
		router = mux.NewRouter()
		server = &http.Server{
			Addr:    fmt.Sprintf("localhost:%d", addr.Port),
			Handler: router,
		}
	)

	return &Server{
		router: router,
		server: server,
	}
}

// Host return server host.
func (s *Server) Host() string {
	return s.server.Addr
}

// SetTLSConfig set the configuration for TLS connection.
func (s *Server) SetTLSConfig(config *tls.Config) *Server {
	s.server.TLSConfig = config

	return s
}

// SetReadTimeout set the timeout for http server.
func (s *Server) SetReadTimeout(timeout time.Duration) *Server {
	s.server.ReadTimeout = timeout

	return s
}

// Serve - call ListenAndServe on http.Server. Listens on the TCP network address srv.Addr and then calls Serve to
// handle requests on incoming connections. Accepted connections are configured to enable TCP keep-alives.
func (s *Server) Serve() {
	go func() {
		err := s.server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()
}

// Close - immediately closes all active net.Listeners and any connections in state StateNew, StateActive,
// or StateIdle. For a graceful shutdown, use Shutdown.
func (s *Server) Close() {
	err := s.server.Close()
	if err != nil {
		panic(err)
	}
}

// HandlerFunc adds a handler for the router. Takes the test handler as the first argument.
func (s *Server) HandlerFunc(t *testing.T, path string, fn func(http.ResponseWriter, *http.Request)) *Server {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.handlers == nil {
		s.handlers = make(map[string]map[*testing.T]func(w http.ResponseWriter, r *http.Request))
	}

	_, ok := s.handlers[path]
	if !ok {
		s.router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) { s.handlers[path][t](w, r) })
		s.handlers[path] = make(map[*testing.T]func(w http.ResponseWriter, r *http.Request))
	}

	_, ok = s.handlers[path][t]
	if !ok {
		s.handlers[path][t] = fn
	}

	return s
}