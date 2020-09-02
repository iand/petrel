package petrel

import (
	"io"
	"log"
	"net"
)

func ListenAndServe(addr string) error {
	server := &Server{Addr: addr}
	return server.ListenAndServe()
}

// A Server holds the configuration for a petrel server. The zero value for Server is a valid configuration.
type Server struct {
	Addr string
}

// ListenAndServe listens on the TCP network address s.Addr and then calls Serve to handle requests on incoming connections.
func (s *Server) ListenAndServe() error {
	addr := s.Addr
	if addr == "" {
		addr = ":3001"
	}
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return s.Serve(l)
}

// Serve accepts incoming connections on the Listener l, creating a new
// service goroutine for each. The service goroutines read requests and then
// call s.Handler to reply to them.
func (s *Server) Serve(l net.Listener) error {
	for {
		conn, err := l.Accept()
		if err != nil {
			// TODO: log failure to accept connection
			continue
		}
		go s.ServeConn(conn)
	}
}

// ServeConn runs the server on a single connection in blocking mode.
func (s *Server) ServeConn(conn io.ReadWriteCloser) {
	codec := TextCodec{}
	var msg Message

	defer conn.Close()
	err := codec.Decode(conn, &msg)
	if err != nil {
		log.Printf("error: %v\n", err)
		return
	}

	log.Printf("data: %+v\n", msg)
}
