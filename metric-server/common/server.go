package common

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type ServerConfig struct {
	Port     string
	Couriers int
}

type Server struct {
	config   ServerConfig
	listener net.Listener
}

func NewServer(config ServerConfig) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", config.Port))
	if err != nil {
		return nil, errors.Wrapf(err, "Could not bind to port %s. Error: %s", config.Port, err)
	}
	server := &Server{
		config:   config,
		listener: listener,
	}
	return server, nil
}

func (s *Server) Run() {
	clientsToServe := make(chan net.Conn, s.config.Couriers)
	for i := 0; i < s.config.Couriers; i++ {
		go ServeClients(clientsToServe)
	}
	for true {
		client_conn, err := s.acceptNewConnection()
		if err != nil {
			logrus.Fatalf("[SERVER] Could not accept new connection. Error: %s", err)
			continue
		}
		clientsToServe <- client_conn
	}
}

func (s *Server) acceptNewConnection() (net.Conn, error) {
	logrus.Infof("[SERVER] Proceed to accept new connections")
	clientConn, err := s.listener.Accept()
	if err != nil {
		return nil, errors.Wrapf(err, "Could not accept new connection")
	}
	logrus.Infof("[SERVER] Got connection from %s", clientConn.LocalAddr())
	return clientConn, nil
}
