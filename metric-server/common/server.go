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
	Port                       string
	Couriers                   int
	MetricEventsBacklog        int
	QueryEventsBacklog         int
	MetricEventsWorkers        int
	QueryEventsWorkers         int
	AlertFileName              string
	AlertMetricId              string
	AlertAggregation           events.AggregationType
	AlertAggregationWindowSecs float32
	AlertLimit                 float32
}

type Server struct {
	config   ServerConfig
	listener net.Listener
}

func NewServer(config ServerConfig) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", config.Port))
	if err != nil {
		return nil, errors.Wrapf(err, "Could not bind to port %s. Error: %s", config.Port, err.Error())
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
	queryEventsToServePool := make([]chan events.Event, s.config.QueryEventsWorkers)
	fileMonitor := file_monitor.NewFileMonitor()

	logrus.Infof("[SERVER] Initializing Metric and Query workers")
	for i := 0; i < s.config.MetricEventsWorkers; i++ {
		metricEventsWorker := NewMetricEventsWorker(metricEventsToServe, fileMonitor)
		go metricEventsWorker.ServeMetricEvents()
	}
	for i := 0; i < s.config.QueryEventsWorkers; i++ {
		queryEventsToServe := make(chan events.Event, s.config.QueryEventsBacklog)
		queryEventsToServePool[i] = queryEventsToServe
		queryEventsWorker := NewQueryEventsWorker(queryEventsToServe, fileMonitor)
		go queryEventsWorker.ServeQueryEvents()
	}

	logrus.Infof("[SERVER] Initializing Alert worker")
	queueForAlerts := make(chan events.Event)
	alertEventsWorker := NewAlertEventsWorker(queueForAlerts, s.config.AlertFileName)
	go alertEventsWorker.ServeAlertEvents()

	logrus.Infof("[SERVER] Initializing Couriers workers")
	for i := 0; i < s.config.Couriers; i++ {
		courier := NewCourier(metricEventsToServe, queryEventsToServePool, queueForAlerts)
		go courier.ServeClients(clientsToServe)
	}

	logrus.Infof("[SERVER] Initializing Clock worker")
	alertEvent := events.AlertEvent{
		MetricId:               s.config.AlertMetricId,
		Aggregation:            s.config.AlertAggregation,
		Limit:                  s.config.AlertLimit,
		AggregationWindowsSecs: s.config.AlertAggregationWindowSecs,
	}
	clockWorker := NewClockWorker(queueForAlerts, alertEvent)
	go clockWorker.Run()
	for true {
		client_conn, err := s.acceptNewConnection()
		if err != nil {
			logrus.Infof("[SERVER] Could not accept new connection. Error: %s", err.Error())
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
	logrus.Infof("[SERVER] Got connection from %s", clientConn.LocalAddr().String())
	return clientConn, nil
}
