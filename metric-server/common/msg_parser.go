package common

import (
	"encoding/binary"
	"math"
	"net"

	"github.com/manudelca/tp1-distribuidos1/metric-server/util"
)

func getLen(clientConn net.Conn) (uint16, error) {
	uint16Size := 16
	bytes, err := util.ReadFromConnection(clientConn, uint16Size)
	if err != nil {
		return 0, err
	}
	len := binary.BigEndian.Uint16(bytes)
	return len, nil
}

func GetMessage(clientConn net.Conn) (Message, error) {
	len, err := getLen(clientConn)
	if err != nil {
		return nil, err
	}
	bytes, err := util.ReadFromConnection(clientConn, int(len))
	if err != nil {
		return nil, err
	}
	switch msgType := bytes[0]; msgType {
	case byte(METRIC):
		return buildMetricMessage(bytes)
	case byte(QUERY):
		return buildQueryMessage(bytes)
		// case default:
		// error!!!
	}
}

func buildMetricMessage(bytes []byte) (MetricMessage, error) {
	value := math.Float32frombits(binary.BigEndian.Uint32(bytes[:32]))
	metricId := string(bytes[32:])
	// Como hago el parseo de los errores?
	return MetricMessage{
		Value:    value,
		MetricId: metricId,
	}, nil
}

func buildQueryMessage(bytes []byte) (QueryMessage, error) {
	// Esto es mas dificil tbh
	return QueryMessage{}, nil
}
