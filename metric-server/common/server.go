package common

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type ServerConfig struct {
	Port    string
	Workers int
}

type Server struct {
	config   ServerConfig
	listener net.Listener
}

func NewServer(config ServerConfig) (*Server, error) {
	// No way to set listenBacklog :(
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
	clientsToServe := make(chan net.Conn, s.config.Workers)
	for i := 0; i < s.config.Workers; i++ {
		go s.serveClients(clientsToServe)
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

func (s *Server) handleClientConnection(clientConn net.Conn) {
	bytesRecv := make([]byte, 1024)
	bytesRead, error := clientConn.Read(bytesRecv)
	if error != nil {
		logrus.Infof(fmt.Sprintf("[SERVER] Error while reading socket %s. Error: %s", clientConn.RemoteAddr(), error))
		clientConn.Close()
		return
	}
	msgRecv := string(bytesRecv[:bytesRead])
	logrus.Infof(fmt.Sprintf("[SERVER] Message received from connection %s. Msg: %s", clientConn.RemoteAddr(), msgRecv))
	msgToSend := fmt.Sprintf("Your Message has been received: %s\n", msgRecv)
	if _, error := clientConn.Write([]byte(msgToSend)); error != nil {
		logrus.Infof(fmt.Sprintf("[SERVER] Error while sending to socket %s. Error: %s", clientConn.RemoteAddr(), error))
	}
	clientConn.Close()
}

func (s *Server) serveClients(clientsToServe chan net.Conn) {
	for true {
		clientConn := <-clientsToServe
		s.handleClientConnection(clientConn)
	}
}
