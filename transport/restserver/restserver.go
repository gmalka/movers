package restserver

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	defaultReadTimeout  = 2 * time.Minute
	defaultWriteTimeout = 20 * time.Second
	defaultAddr         = ""
	defaultPort         = "8080"
)

type ServerBuilder struct {
	handler      http.Handler
	addr         string
	port         string
	readTimeout  time.Duration
	writeTimeout time.Duration
}

func NewServerBuilder(handler http.Handler) *ServerBuilder {
	return &ServerBuilder{
		handler:      handler,
		port:         fmt.Sprintf("%v:%v", defaultAddr, defaultPort),
		readTimeout:  defaultReadTimeout,
		writeTimeout: defaultWriteTimeout,
	}
}

func (b *ServerBuilder) Addr(port string) *ServerBuilder {
	b.addr = port
	return b
}

func (b *ServerBuilder) Port(port string) *ServerBuilder {
	b.port = port
	return b
}

func (b *ServerBuilder) ReadTimeout(timeout time.Duration) *ServerBuilder {
	b.readTimeout = timeout
	return b
}

func (b *ServerBuilder) WriteTimeout(timeout time.Duration) *ServerBuilder {
	b.writeTimeout = timeout
	return b
}

func (b *ServerBuilder) Build() *Server {
	httpServer := &http.Server{
		Handler:      b.handler,
		ReadTimeout:  b.readTimeout,
		WriteTimeout: b.writeTimeout,
		Addr:         fmt.Sprintf("%v:%v", b.addr, b.port),
	}

	s := &Server{
		server: httpServer,
		errs:   make(chan error, 1),
	}

	s.start()

	return s
}

type Server struct {
	server *http.Server
	errs   chan error
}

func (s *Server) start() {
	go func() {
		defer close(s.errs)
		s.errs <- s.server.ListenAndServe()
	}()
}

func (s *Server) Errs() <-chan error {
	return s.errs
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}
