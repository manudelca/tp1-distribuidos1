package common

import (
	"net"

	"github.com/manudelca/tp1-distribuidos1/metric-server/events"
	"github.com/manudelca/tp1-distribuidos1/metric-server/protocol"
	"github.com/sirupsen/logrus"
)

type Courier struct {
	metricEventsQueue chan events.MetricEvent
	queryEventsQueue  chan events.QueryEvent
}

func NewCourier(metricEventsQueue chan events.MetricEvent, queryEventsQueue chan events.QueryEvent) *Courier {
	courier := Courier{
		metricEventsQueue: metricEventsQueue,
		queryEventsQueue:  queryEventsQueue,
	}
	return &courier
}

func (c *Courier) ServeClients(clientsToServe chan net.Conn) {
	for clientConn := range clientsToServe {
		c.handleClientConnection(clientConn)
		logrus.Infof("[COURIER] Closing client connection")
		clientConn.Close()
	}
}

func (c *Courier) handleClientConnection(clientConn net.Conn) {
	event, err := protocol.GetEventFromMessage(clientConn)
	if err != nil {
		logrus.Fatalf("[COURIER] Error trying to getEventFromMessage. Error: %s", err)
		return
	}
	logrus.Infof("[COURIER] Event type %d succesfully parsed", event.GetType())
	logrus.Infof("[COURIER] Event parsed: ", event)
	if metricEvent, ok := event.(events.MetricEvent); ok {
		c.answerMetricEvent(metricEvent, clientConn)
	} else if queryEvent, ok := event.(events.QueryEvent); ok {
		c.answerQueryEvent(queryEvent, clientConn)
	} else {
		logrus.Fatalf("[COURIER] Event type assertion failed")
	}
}

func (c *Courier) answerMetricEvent(metricEvent events.MetricEvent, clientConn net.Conn) {
	if cap(c.metricEventsQueue) == len(c.metricEventsQueue) {
		logrus.Infof("[COURIER] MetricEventsQueue full. Rejecting client")
		c.rejectClient(clientConn)
		return
	}
	logrus.Infof("[COURIER] Storing metric event in queue: ", metricEvent)
	c.metricEventsQueue <- metricEvent
	err := protocol.SendSuccess("Metric succesfully received", clientConn)
	if err != nil {
		logrus.Fatalf("[COURIER] An error ocurred while trying to answer client metric. Error: ", err)
	}
}

func (c *Courier) answerQueryEvent(queryEvent events.QueryEvent, clientConn net.Conn) {
	if cap(c.queryEventsQueue) == len(c.queryEventsQueue) {
		logrus.Infof("[COURIER] QueryEventsQueue full. Rejecting client")
		c.rejectClient(clientConn)
		return
	}
	logrus.Infof("[COURIER] Storing query event in queue: ", queryEvent)
	c.queryEventsQueue <- queryEvent
}

func (c *Courier) rejectClient(clientConn net.Conn) {
	err := protocol.SendServerError("Server not available. Try again later", clientConn)
	if err != nil {
		logrus.Fatalf("[COURIER] An error ocurred while trying to reject client. Error: ", err)
	}
}
