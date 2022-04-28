package common

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

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
	serverOn bool
}

func NewServer(config ServerConfig) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", config.Port))
	if err != nil {
		return nil, errors.Wrapf(err, "Could not bind to port %s. Error: %s", config.Port, err.Error())
	}
	server := &Server{
		config:   config,
		listener: listener,
		serverOn: true,
	}
	return server, nil
}

func (s *Server) Run() {
	clientsToServe := make(chan net.Conn, s.config.Couriers)
	metricEventsToServe := make(chan events.MetricEvent, s.config.MetricEventsBacklog)
	queryEventsToServePool := make([]chan events.Event, s.config.QueryEventsWorkers)
	fileMonitor := file_monitor.NewFileMonitor()

	logrus.Infof("[SERVER] Initializing Metric and Query workers")
	var waitMetricWorkers sync.WaitGroup
	for i := 0; i < s.config.MetricEventsWorkers; i++ {
		waitMetricWorkers.Add(1)
		metricEventsWorker := NewMetricEventsWorker(metricEventsToServe, fileMonitor, &waitMetricWorkers)
		go metricEventsWorker.ServeMetricEvents()
	}
	var waitQueryEventsWorkers sync.WaitGroup
	for i := 0; i < s.config.QueryEventsWorkers; i++ {
		waitQueryEventsWorkers.Add(1)
		queryEventsToServe := make(chan events.Event, s.config.QueryEventsBacklog)
		queryEventsToServePool[i] = queryEventsToServe
		queryEventsWorker := NewQueryEventsWorker(queryEventsToServe, fileMonitor, &waitQueryEventsWorkers)
		go queryEventsWorker.ServeQueryEvents()
	}

	logrus.Infof("[SERVER] Initializing Alert worker")
	alertEventsWorkerShutdown := make(chan bool)
	queueForAlerts := make(chan events.Event)
	alertEventsWorker := NewAlertEventsWorker(queueForAlerts, s.config.AlertFileName, alertEventsWorkerShutdown)
	go alertEventsWorker.ServeAlertEvents()

	logrus.Infof("[SERVER] Initializing Couriers workers")
	var waitCouriers sync.WaitGroup
	for i := 0; i < s.config.Couriers; i++ {
		waitCouriers.Add(1)
		courier := NewCourier(metricEventsToServe, queryEventsToServePool, queueForAlerts, &waitCouriers)
		go courier.ServeClients(clientsToServe)
	}

	logrus.Infof("[SERVER] Initializing Clock worker")
	alertEvent := events.AlertEvent{
		MetricId:               s.config.AlertMetricId,
		Aggregation:            s.config.AlertAggregation,
		Limit:                  s.config.AlertLimit,
		AggregationWindowsSecs: s.config.AlertAggregationWindowSecs,
	}
	clockShutdownCommand := make(chan bool)
	clockIsOff := make(chan bool)
	clockWorker := NewClockWorker(queueForAlerts, alertEvent, clockShutdownCommand, clockIsOff)
	go clockWorker.Run()

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, syscall.SIGTERM)
	go s.handleSignal(signalChannel)
	for s.serverOn {
		client_conn, err := s.acceptNewConnection()
		if err != nil {
			logrus.Infof("[SERVER] Could not accept new connection. Error: %s", err.Error())
			continue
		}
		clientsToServe <- client_conn
	}

	// Shutdown Couriers And Clock
	close(clientsToServe)
	waitCouriers.Wait()
	logrus.Infof("[SERVER] Ended all couriers")
	clockShutdownCommand <- true
	close(clockIsOff)
	logrus.Infof("[SERVER] Ended clock")

	// Shutdown
	close(metricEventsToServe)
	for _, queryEventsToServe := range queryEventsToServePool {
		close(queryEventsToServe)
	}
	close(queueForAlerts)
	waitMetricWorkers.Wait()
	logrus.Infof("[SERVER] Ended all metric event workers")
	waitQueryEventsWorkers.Wait()
	logrus.Infof("[SERVER] Ended all query event workers")
	<-alertEventsWorkerShutdown
	logrus.Infof("[SERVER] Ended alert event worker")

	logrus.Infof("[SERVER] Goodbye :)")
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

func (s *Server) handleSignal(signalChannel chan os.Signal) {
	for signal := range signalChannel {
		switch signal {
		case syscall.SIGTERM:
			logrus.Infof("[SERVER] SIGTERM received. Proceeding to shutdown")
			s.Shutdown()
		}
	}
}

func (s *Server) Shutdown() {
	s.serverOn = false
	s.listener.Close()
}
