package common

import (
	"net"
	"reflect"

	"github.com/manudelca/tp1-distribuidos1/metric-server/events"
	"github.com/manudelca/tp1-distribuidos1/metric-server/protocol"
	"github.com/sirupsen/logrus"
)

type Courier struct {
	metricEventsQueue chan events.MetricEvent
	queryEventsPool   []chan events.QueryEvent
}

func NewCourier(metricEventsQueue chan events.MetricEvent, queryEventsPool []chan events.QueryEvent) *Courier {
	courier := Courier{
		metricEventsQueue: metricEventsQueue,
		queryEventsPool:   queryEventsPool,
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
		logrus.Infof("[COURIER] Error trying to getEventFromMessage. Error: %s", err.Error())
		return
	}
	logrus.Infof("[COURIER] Event type %d succesfully parsed", event.GetType())
	logrus.Infof("[COURIER] Event parsed: ", event)
	if metricEvent, ok := event.(events.MetricEvent); ok {
		c.answerMetricEvent(metricEvent, clientConn)
	} else if queryEvent, ok := event.(events.QueryEvent); ok {
		c.answerQueryEvent(queryEvent, clientConn)
	} else {
		logrus.Infof("[COURIER] Event type assertion failed")
	}
}

func (c *Courier) answerMetricEvent(metricEvent events.MetricEvent, clientConn net.Conn) {
	select {
	case c.metricEventsQueue <- metricEvent:
	default:
		logrus.Infof("[COURIER] MetricEventsQueue full. Rejecting client")
		c.rejectClient(clientConn)
		return
	}
	logrus.Infof("[COURIER] Stored metric event in queue: ", metricEvent)
	c.metricEventsQueue <- metricEvent
	err := protocol.SendSuccess("Metric succesfully received", clientConn)
	if err != nil {
		logrus.Infof("[COURIER] An error ocurred while trying to answer client metric. Error: %s", err.Error())
	}
}

func (c *Courier) answerQueryEvent(queryEvent events.QueryEvent, clientConn net.Conn) {
	cases := make([]reflect.SelectCase, len(c.queryEventsPool)+1)
	for i, query := range c.queryEventsPool {
		cases[i] = reflect.SelectCase{Dir: reflect.SelectSend, Chan: reflect.ValueOf(query), Send: reflect.ValueOf(queryEvent)}
	}
	defualtCaseIndex := len(c.queryEventsPool)
	cases[defualtCaseIndex] = reflect.SelectCase{Dir: reflect.SelectDefault, Chan: reflect.Value{}, Send: reflect.Value{}}
	chosen, _, _ := reflect.Select(cases)
	if chosen == defualtCaseIndex {
		logrus.Infof("[COURIER] QueryEventsQueue full. Rejecting client")
		c.rejectClient(clientConn)
		return
	}
	logrus.Infof("[COURIER] Stored query event in queue %d: ", chosen, queryEvent)
}

func (c *Courier) rejectClient(clientConn net.Conn) {
	err := protocol.SendServerError("Server not available. Try again later", clientConn)
	if err != nil {
		logrus.Infof("[COURIER] An error ocurred while trying to reject client. Error: %s", err.Error())
	}
}
