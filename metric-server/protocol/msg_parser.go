package protocol

import (
	"encoding/binary"
	"net"

	"github.com/manudelca/tp1-distribuidos1/metric-server/util"
	"github.com/pkg/errors"
)

func sendMessage(message Message, clientConn net.Conn) error {
	var bytesToBeSent []byte
	bytesToBeSent = append(bytesToBeSent, byte(message.MsgType))
	bytesToBeSent = append(bytesToBeSent, []byte(message.Message)[:]...)
	bytesLen := uint16(len(bytesToBeSent))
	var bytesLenArray [2]byte
	binary.BigEndian.PutUint16(bytesLenArray[:], bytesLen)
	bytesToBeSent = append(bytesLenArray[:], bytesToBeSent...)

	return util.SendToConnection(clientConn, bytesToBeSent)
}

func SendServerError(errorMsg string, clientConn net.Conn) error {
	message := Message{MsgType: SERVERERROR, Message: errorMsg}
	err := sendMessage(message, clientConn)
	if err != nil {
		return errors.Wrapf(err, "Could not send server error")
	}
	return nil
}
