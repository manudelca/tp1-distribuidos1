package common

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"strconv"

	"github.com/manudelca/tp1-distribuidos1/metric-server/util"
	"github.com/pkg/errors"
)

type InvalidMessageTypeError struct {
	msgType MsgType
}

func (e InvalidMessageTypeError) Error() string {
	return fmt.Sprintf("Invalid message type: %d", e.msgType)
}

type InvalidMetricMessageFormatError struct {
	errorMsg string
}

func (e InvalidMetricMessageFormatError) Error() string {
	return e.errorMsg
}

func getLen(clientConn net.Conn) (uint16, error) {
	uint16BytesSize := 2
	bytes, err := util.ReadFromConnection(clientConn, uint16BytesSize)
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
	default:
		return nil, InvalidMessageTypeError{msgType: MsgType(msgType)}
	}
}

func buildMetricMessage(bytes []byte) (MetricMessage, error) {
	if len(bytes) < (4 + 1) {
		return MetricMessage{}, InvalidMetricMessageFormatError{errorMsg: "There should be a 32 bits float and at least 1 character for MetricId"}
	}
	if _, err := strconv.ParseFloat(string(bytes[:4]), 32); err != nil {
		return MetricMessage{}, errors.Wrapf(err, "Value field could not be parsed as 32 bits Float")
	}
	value := math.Float32frombits(binary.BigEndian.Uint32(bytes[:4]))
	metricId := string(bytes[4:])
	return MetricMessage{
		Value:    value,
		MetricId: metricId,
	}, nil
}

func buildQueryMessage(bytes []byte) (QueryMessage, error) {
	// Esto es mas dificil tbh
	return QueryMessage{}, nil
}
