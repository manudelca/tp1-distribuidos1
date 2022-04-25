package common

import (
	"fmt"
	"net"

	"github.com/manudelca/tp1-distribuidos1/metric-server/events"
	"github.com/manudelca/tp1-distribuidos1/metric-server/file_monitor"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type ServerConfig struct {
	Port                string
	Couriers            int
	MetricEventsBacklog int
	QueryEventsBacklog  int
	MetricEventsWorkers int
	QueryEventsWorkers  int
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
	metricEventsToServe := make(chan events.MetricEvent, s.config.MetricEventsBacklog)
	queryEventsToServe := make(chan events.QueryEvent, s.config.QueryEventsBacklog)
	fileMonitor := file_monitor.FileMonitor{}
	for i := 0; i < s.config.Couriers; i++ {
		courier := NewCourier(metricEventsToServe, queryEventsToServe)
		go courier.ServeClients(clientsToServe)
	}
	for i := 0; i < s.config.MetricEventsWorkers; i++ {
		metricEventsWorker := NewMetricEventsWorker(metricEventsToServe, &fileMonitor)
		go metricEventsWorker.ServeMetricEvents()
	}
	for i := 0; i < s.config.QueryEventsWorkers; i++ {
		queryEventsWorker := NewQueryEventsWorker(queryEventsToServe, &fileMonitor)
		go queryEventsWorker.ServeQueryEvents()
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
