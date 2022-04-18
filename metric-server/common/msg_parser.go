package common

import (
	"encoding/binary"
	"net"

	"github.com/manudelca/tp1-distribuidos1/metric-server/util"
)

// struct Message {
// 	Msg

// }

func getLen(clientConn net.Conn) (uint16, error) {
	uint16Size := 16
	bytes := make([]byte, uint16Size)
	err := util.ReadFromConnection(clientConn, bytes, uint16Size)
	if err != nil {
		return 0, err
	}
	len := binary.BigEndian.Uint16(bytes)
	return len, nil
}

// func GetMessage(clientConn net.Conn) (, error) {
// 	len, err := getLen(clientConn)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// msgRecv := string(bytesRecv[:bytesRead])
// }
