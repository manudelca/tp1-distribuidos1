package common

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
)

func HandleClientConnection(clientConn net.Conn) {
	bytesRecv := make([]byte, 1024)
	bytesRead, error := clientConn.Read(bytesRecv)
	if error != nil {
		logrus.Infof(fmt.Sprintf("[SERVER] Error while reading socket %s. Error: %s", clientConn.RemoteAddr(), error))
		clientConn.Close()
		return
	}
	msgRecv := string(bytesRecv[:bytesRead])
	logrus.Infof(fmt.Sprintf("[SERVER] Message received from connection %s. Msg: %s", clientConn.RemoteAddr(), msgRecv))
	msgToSend := fmt.Sprintf("Your Message has been received: %s\n", msgRecv)
	if _, error := clientConn.Write([]byte(msgToSend)); error != nil {
		logrus.Infof(fmt.Sprintf("[SERVER] Error while sending to socket %s. Error: %s", clientConn.RemoteAddr(), error))
	}
	clientConn.Close()
}

func ServeClients(clientsToServe chan net.Conn) {
	for clientConn := range clientsToServe {
		HandleClientConnection(clientConn)
	}
}
