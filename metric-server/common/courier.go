package common

import (
	"net"

	"github.com/manudelca/tp1-distribuidos1/metric-server/events"
	"github.com/manudelca/tp1-distribuidos1/metric-server/protocol"
	"github.com/sirupsen/logrus"
)

func ServeClients(clientsToServe chan net.Conn, metricEventsQueue chan events.MetricEvent, queryEventsQueue chan events.QueryEvent) {
	for clientConn := range clientsToServe {
		handleClientConnection(clientConn, metricEventsQueue, queryEventsQueue)
	}
}

func handleClientConnection(clientConn net.Conn, metricEventsQueue chan events.MetricEvent, queryEventsQueue chan events.QueryEvent) {
	event, err := protocol.GetEventFromMessage(clientConn)
	if err != nil {
		logrus.Fatalf("[COURIER] Error trying to getEventFromMessage. Error: %s", err)
	}
	logrus.Infof("[COURIER] Event type %d succesfully parsed", event.GetType())
	logrus.Infof("[COURIER] Event parsed: ", event)
	if metricEvent, ok := event.(events.MetricEvent); ok {
		answerMetricEvent(metricEvent, clientConn, metricEventsQueue)
	} else if queryEvent, ok := event.(events.QueryEvent); ok {
		answerQueryEvent(queryEvent, clientConn, queryEventsQueue)
	} else {
		logrus.Fatalf("[COURIER] Event type assertion failed")
	}
	logrus.Infof("[COURIER] Closing client connection")
	clientConn.Close()
}

func answerMetricEvent(metricEvent events.MetricEvent, clientConn net.Conn, metricEventsQueue chan events.MetricEvent) {
	if cap(metricEventsQueue) == len(metricEventsQueue) {
		logrus.Infof("[COURIER] MetricEventsQueue full. Rejecting client")
		RejectClient(clientConn)
		return
	}
	logrus.Infof("[COURIER] Storing metric event in queue: ", metricEvent)
	metricEventsQueue <- metricEvent
	err := protocol.SendSuccess("Metric succesfully received", clientConn)
	if err != nil {
		logrus.Fatalf("[COURIER] An error ocurred while trying to answer client metric. Error: ", err)
	}
}

func answerQueryEvent(queryEvent events.QueryEvent, clientConn net.Conn, queryEventsQueue chan events.QueryEvent) {
	if cap(queryEventsQueue) == len(queryEventsQueue) {
		logrus.Infof("[COURIER] QueryEventsQueue full. Rejecting client")
		RejectClient(clientConn)
		return
	}
	logrus.Infof("[COURIER] Storing query event in queue: ", queryEvent)
	queryEventsQueue <- queryEvent
}

func RejectClient(clientConn net.Conn) {
	err := protocol.SendServerError("Server not available. Try again later", clientConn)
	if err != nil {
		logrus.Fatalf("[COURIER] An error ocurred while trying to reject client. Error: ", err)
	}
}
