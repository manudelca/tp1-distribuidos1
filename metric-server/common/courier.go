package common

import (
	"net"

	"github.com/manudelca/tp1-distribuidos1/metric-server/protocol"
	"github.com/sirupsen/logrus"
)

func HandleClientConnection(clientConn net.Conn) {
	event, err := protocol.GetEventFromMessage(clientConn)
	if err != nil {
		logrus.Fatalf("[SERVER] Error trying to getEventFromMessage. Error: %s", err)
	}
	logrus.Infof("Event type %d succesfully processed", event.GetType())
	logrus.Infof("Event processed: ", event)
	// switch event.GetType() {
	// case events.METRIC:
	// case events.QUERY:
	// }
	clientConn.Close()
}

func ServeClients(clientsToServe chan net.Conn) {
	for clientConn := range clientsToServe {
		HandleClientConnection(clientConn)
	}
}
