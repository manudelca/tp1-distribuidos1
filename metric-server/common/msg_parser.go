package common

import (
	"encoding/binary"
	"net"

	"github.com/manudelca/tp1-distribuidos1/metric-server/util"
)

type MsgType uint8

const (
	Metric MsgType = iota
	Query
)

// struct Message {
// 	Msg

// }

func getLen(clientConn net.Conn) (uint16, error) {
	uint16Size := 16
	bytes, err := util.ReadFromConnection(clientConn, uint16Size)
	if err != nil {
		return 0, err
	}
	len := binary.BigEndian.Uint16(bytes)
	return len, nil
}

func GetMessage(clientConn net.Conn) (, error) {
 	len, err := getLen(clientConn)
 	if err != nil {
 		return nil, err
 	}
	bytes, err := util.ReadFromConnection(clientConn, len)
	if err != nil {
		return nil, err
	}
	switch msgType := bytes[0]; msgType {
	case byte(Metric):
		
	case byte(Query):
	
	case default:
		// error!!!
	}
	// Y ahora? Como creo el struct y etc? Igual estoy un toque mataduki tbh...
	// Sigo con el otro tp
}
