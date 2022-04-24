package common

import (
	"net"

	"github.com/manudelca/tp1-distribuidos1/metric-server/events"
	"github.com/manudelca/tp1-distribuidos1/metric-server/protocol"
	"github.com/sirupsen/logrus"
)

func HandleClientConnection(clientConn net.Conn, metricEventsQueue chan events.MetricEvent, queryEventsQueue chan events.QueryEvent) {
	event, err := protocol.GetEventFromMessage(clientConn)
	if err != nil {
		logrus.Fatalf("[COURIER] Error trying to getEventFromMessage. Error: %s", err)
	}
	logrus.Infof("[COURIER] Event type %d succesfully processed", event.GetType())
	logrus.Infof("[COURIER] Event processed: ", event)
	if metricEvent, ok := event.(events.MetricEvent); ok {
		metricEventsQueue <- metricEvent
	} else if queryEvent, ok := event.(events.QueryEvent); ok {
		queryEventsQueue <- queryEvent
	} else {
		logrus.Fatalf("[COURIER] Event type assertion failed")
	}
	// switch event.GetType() {
	// case events.METRIC:
	// 	metricEventsQueue <- events.MetricEvent.(event)
	// case events.QUERY:
	// 	queryEventsQueue <- events.QueryEvent(event)
	// }
	clientConn.Close()
}

func ServeClients(clientsToServe chan net.Conn, metricEventsQueue chan events.MetricEvent, queryEventsQueue chan events.QueryEvent) {
	for clientConn := range clientsToServe {
		HandleClientConnection(clientConn, metricEventsQueue, queryEventsQueue)
	}
}
