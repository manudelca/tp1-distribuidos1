package protocol

import (
	"encoding/binary"
	"fmt"
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
	fieldsReceived := make(map[FieldType]bool)
	for i := 0; i < len(message); {
		switch fieldType := message[i]; fieldType {
		case byte(METRICID):
			metricId, newIndex, err := parseMetricId(message[i:], i)
			if err != nil {
				return common.MetricEvent{}, errors.Wrapf(err, "Could not parse MetricId when trying to build MetricEvent from message received")
			}
			metricEvent.MetricId = metricId
			i = newIndex
			fieldsReceived[FieldType(fieldType)] = true
		case byte(VALUE):
			value, newIndex, err := parseValue(message[i:], i)
			if err != nil {
				return common.MetricEvent{}, errors.Wrapf(err, "Could not parse Value when trying to build MetricEvent from message received")
			}
			metricEvent.Value = value
			i = newIndex
			fieldsReceived[FieldType(fieldType)] = true
		default:
			err := fmt.Sprintf("Invalid message format for MetricEvent. Unrecognized Field type %d", fieldType)
			return common.MetricEvent{}, InvalidMessageFormatError{errorMsg: err}
		}
	}
	areAllFieldsPresent := fieldsReceived[METRICID] && fieldsReceived[VALUE]
	if !areAllFieldsPresent {
		return common.MetricEvent{}, InvalidMessageFormatError{errorMsg: "Missing required fields. MetricID and Value are needed"}
	}
	return metricEvent, nil
}

func buildQueryMessage(message []byte) (common.QueryEvent, error) {
	query := common.QueryEvent{}
	fieldsReceived := make(map[FieldType]bool)
	for i := 0; i < len(message); {
		switch fieldType := message[i]; fieldType {
		case byte(METRICID):
			metricId, newIndex, err := parseMetricId(message[i:], i)
			if err != nil {
				return common.QueryEvent{}, errors.Wrapf(err, "Could not parse MetricId when trying to build QueryMessage from message received")
			}
			query.MetricId = metricId
			i = newIndex
		case byte(AGGREGATION):
			aggregation, newIndex, err := parseAggregation(message[i:], i)
			if err != nil {
				return common.QueryEvent{}, errors.Wrapf(err, "Could not parse Aggregation when trying to build QueryMessage from message received")
			}
			query.Aggregation = aggregation
			i = newIndex
		case byte(AGGREGATIONWINDOWSSECS):
			aggregationWindowsSecs, newIndex, err := parseAggregationWindowsSecs(message[i:], i)
			if err != nil {
				return common.QueryEvent{}, errors.Wrapf(err, "Could not parse AggregationWindowsSecs when trying to build QueryMessage from message received")
			}
			query.AggregationWindowsSecs = aggregationWindowsSecs
			i = newIndex
		case byte(DATEINTERVAL):
			from, to, newIndex, err := parseDateInterval(message[i:], i)
			if err != nil {
				return common.QueryEvent{}, errors.Wrapf(err, "Could not parse Date interval when trying to build QueryMessage from message received")
			}
			query.From = from
			query.To = to
			i = newIndex
		default:
			err := fmt.Sprintf("Invalid message format for MetricEvent. Unrecognized Field type %d", fieldType)
			return common.QueryEvent{}, InvalidMessageFormatError{errorMsg: err}
		}
	}
	areAllFieldsPresent := fieldsReceived[METRICID] && fieldsReceived[AGGREGATION] && fieldsReceived[AGGREGATIONWINDOWSSECS]
	if !areAllFieldsPresent {
		return common.QueryEvent{}, InvalidMessageFormatError{errorMsg: "Missing required fields. MetricID, Aggregation, AggregationWindowsSecs are needed"}
	}
	return query, nil
}
