package protocol

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"strconv"

	"github.com/manudelca/tp1-distribuidos1/metric-server/common"
	"github.com/manudelca/tp1-distribuidos1/metric-server/util"
	"github.com/pkg/errors"
)

func getLen(clientConn net.Conn) (uint16, error) {
	uint16Size := 2
	bytes, err := util.ReadFromConnection(clientConn, uint16Size)
	if err != nil {
		return 0, errors.Wrapf(err, "Could not read message len from connection")
	}
	if _, err := strconv.ParseUint(string(bytes), 10, 16); err != nil {
		return 0, errors.Wrapf(err, "Len field could not be parsed as 16 bits uint")
	}
	len := binary.BigEndian.Uint16(bytes)
	return len, nil
}

func GetMessage(clientConn net.Conn) (common.Event, error) {
	len, err := getLen(clientConn)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get length of message")
	}
	bytes, err := util.ReadFromConnection(clientConn, int(len))
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read message")
	}
	switch eventType := bytes[0]; eventType {
	case byte(common.METRIC):
		return buildMetricMessage(bytes[1:])
	case byte(common.QUERY):
		return buildQueryMessage(bytes[1:])
	default:
		return nil, InvalidEventTypeError{eventType: uint8(eventType)}
	}
}

func buildMetricMessage(message []byte) (common.MetricEvent, error) {
	metricEvent := common.MetricEvent{}
	for i := 0; i < len(bytes); {
		switch fieldType := bytes[i]; fieldType {
		case byte(METRICID):
			metricId, newIndex , err := parseMetricId(message[i:], i)
			if err != nil {
				return nil, errors.Wrapf(err, "Could not parse MetricId when trying to build MetricEvent from message received")
			}
			metricEvent.MetricId = metricId
			i = newIdex
		case byte(VALUE):
			value, newIndex , err := parseValue(message[i:], i)
			if err != nil {
				return nil, errors.Wrapf(err, "Could not parse Value when trying to build MetricEvent from message received")
			}
			metricEvent.Value = value
			i = newIdex
		default:
			err := fmt.SPrintf("Invalid message format for MetricEvent. Unrecognized Field type %d", fieldType) 
			return nil, InvalidMessageFormatError{errorMsg: err}
		}
	}
	return metricEvent.Validate()
}

func buildQueryMessage(bytes []byte) (common.QueryEvent, error) {
	query := common.QueryEvent{}
	for i := 0; i < len(bytes); {
		switch (bytes[i]) {
		case byte(METRICID):
		case byte(VALUE):
		case byte(AGGREGATION):
			// value, i, err := getAggregationField()
			// if err {
			// 
			// }
			// metric.Field = value
		case byte(AGGREGATIONWINDOWSSECS):
		case byte(FROM):
		case byte(TO):
		default:
			return nil, //Error custom
		}
	}
	return query.Validate()
}

// func buildQueryMessage(bytes []byte) (QueryMessage, error) {
// 	// Esto es mas dificil tbh
// 	aggregation := AggregationType(bytes[0])
// 	if _, err := strconv.ParseFloat(string(bytes[1:4]), 32); err != nil {
// 		return QueryMessage{}, errors.Wrapf(err, "\"Value\" field could not be parsed as 32 bits Float")
// 	}
// 	aggregationWindowsSecs := math.Float32frombits(binary.BigEndian.Uint32(bytes[1:4]))

// 	layout := "2006-01-02 03:04:05"
// 	if bytes[5] > 0 {
// 		from, err := time.Parse(layout, string(bytes[6:len(layout)+6]))
// 		if err != nil {

// 		}
// 		to, err := time.Parse(layout, string(bytes[len(layout)+6:len(layout)+6+len(layout)+6]))
// 		if err != nil {

// 		}
// 	}
// 	hasTimeInterval := (bytes[5])

// 	return QueryMessage{
// 		Aggregation:            aggregation,
// 		AggregationWindowsSecs: aggregationWindowsSecs,
// 		To:                     to,
// 		From:                   from,
// 	}, nil
// }
