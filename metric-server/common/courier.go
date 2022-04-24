package common

import (
	"net"

	"github.com/manudelca/tp1-distribuidos1/metric-server/protocol"
	"github.com/sirupsen/logrus"
)

func HandleClientConnection(clientConn net.Conn) {
	_, err := protocol.GetEventFromMessage(clientConn)
	if err != nil {
		logrus.Fatalf("[SERVER] Error trying to getEventFromMessage. Error: %s", err)
	}
	// logrus.Infof("Event type %s succesfully processed", event.GetType())
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
