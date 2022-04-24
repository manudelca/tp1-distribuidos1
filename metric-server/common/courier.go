package common

import (
	"net"

	"github.com/manudelca/tp1-distribuidos1/metric-server/events"
	"github.com/manudelca/tp1-distribuidos1/metric-server/protocol"
	"github.com/sirupsen/logrus"
)

func RejectClient(clientConn net.Conn) {
	err := protocol.SendServerError("Server not available. Try again later", clientConn)
	if err != nil {
		logrus.Fatalf("[COURIER] An error ocurred while trying to reject client. Error: ", err)
	}
}

func HandleClientConnection(clientConn net.Conn, metricEventsQueue chan events.MetricEvent, queryEventsQueue chan events.QueryEvent) {
	event, err := protocol.GetEventFromMessage(clientConn)
	if err != nil {
		logrus.Fatalf("[COURIER] Error trying to getEventFromMessage. Error: %s", err)
	}
	logrus.Infof("[COURIER] Event type %d succesfully parsed", event.GetType())
	logrus.Infof("[COURIER] Event parsed: ", event)
	if metricEvent, ok := event.(events.MetricEvent); ok {
		if cap(metricEventsQueue) == len(metricEventsQueue) {
			logrus.Infof("[COURIER] MetricEventsQueue full. Rejecting client")
			RejectClient(clientConn)
			clientConn.Close()
			return
		}
		metricEventsQueue <- metricEvent
	} else if queryEvent, ok := event.(events.QueryEvent); ok {
		if cap(queryEventsQueue) == len(queryEventsQueue) {
			logrus.Infof("[COURIER] QueryEventsQueue full. Rejecting client")
			RejectClient(clientConn)
			clientConn.Close()
			return
		}
		queryEventsQueue <- queryEvent
	} else {
		logrus.Fatalf("[COURIER] Event type assertion failed")
	}
	clientConn.Close()
}

func ServeClients(clientsToServe chan net.Conn, metricEventsQueue chan events.MetricEvent, queryEventsQueue chan events.QueryEvent) {
	for clientConn := range clientsToServe {
		HandleClientConnection(clientConn, metricEventsQueue, queryEventsQueue)
	}
}
