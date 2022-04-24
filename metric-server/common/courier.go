package common

import (
	"net"
)

func HandleClientConnection(clientConn net.Conn) {
	// event, err := protocol.GetEventFromMessage(clientConn)
	// if err != nil {
	// 	// What?
	// }
	clientConn.Close()
}

func ServeClients(clientsToServe chan net.Conn) {
	for clientConn := range clientsToServe {
		HandleClientConnection(clientConn)
	}
}
