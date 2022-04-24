package protocol

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/manudelca/tp1-distribuidos1/metric-server/events"
	"github.com/manudelca/tp1-distribuidos1/metric-server/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func getLen(clientConn net.Conn) (uint16, error) {
	uint16Size := 2
	bytes, err := util.ReadFromConnection(clientConn, uint16Size)
	if err != nil {
		logrus.Fatalf("Failed to get message len. Error: %s", err)
		return 0, errors.Wrapf(err, "Could not read message len from connection")
	}
	len := binary.BigEndian.Uint16(bytes)
	logrus.Infof("Message len readed: %d", len)
	return len, nil
}

func GetEventFromMessage(clientConn net.Conn) (events.Event, error) {
	len, err := getLen(clientConn)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to get length of message")
	}
	bytes, err := util.ReadFromConnection(clientConn, int(len))
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read message")
	}
	logrus.Infof("Bytes received: ", bytes)
	logrus.Infof("EventType received: %d", bytes[0])
	switch eventType := bytes[0]; eventType {
	case byte(events.METRIC):
		return buildMetricMessage(bytes[1:])
	case byte(events.QUERY):
		return buildQueryMessage(bytes[1:])
	default:
		return nil, InvalidEventTypeError{eventType: uint8(eventType)}
	}
}

func buildMetricMessage(message []byte) (events.MetricEvent, error) {
	logrus.Infof("About to build MetricEvent")
	metricEvent := events.MetricEvent{}
	fieldsReceived := make(map[FieldType]bool)
	for i := 0; i < len(message); {
		logrus.Infof("FieldType received: %d", message[i])
		switch fieldType := message[i]; fieldType {
		case byte(METRICID):
			logrus.Infof("About to parse MetricId field")
			metricId, newIndex, err := parseMetricId(message[i+1:], i+1)
			logrus.Infof("MetricId parsed: %s", metricId)
			if err != nil {
				return events.MetricEvent{}, errors.Wrapf(err, "Could not parse MetricId when trying to build MetricEvent from message received")
			}
			metricEvent.MetricId = metricId
			i = newIndex
			fieldsReceived[FieldType(fieldType)] = true
		case byte(VALUE):
			logrus.Infof("About to parse value field")
			value, newIndex, err := parseValue(message[i+1:], i+1)
			logrus.Infof("Value parsed: %d", value)
			if err != nil {
				return events.MetricEvent{}, errors.Wrapf(err, "Could not parse Value when trying to build MetricEvent from message received")
			}
			metricEvent.Value = value
			i = newIndex
			fieldsReceived[FieldType(fieldType)] = true
		default:
			err := fmt.Sprintf("Invalid message format for MetricEvent. Unrecognized Field type %d", fieldType)
			return events.MetricEvent{}, InvalidMessageFormatError{errorMsg: err}
		}
	}
	areAllFieldsPresent := fieldsReceived[METRICID] && fieldsReceived[VALUE]
	if !areAllFieldsPresent {
		return events.MetricEvent{}, InvalidMessageFormatError{errorMsg: "Missing required fields. MetricID and Value are needed"}
	}
	return metricEvent, nil
}

func buildQueryMessage(message []byte) (events.QueryEvent, error) {
	logrus.Infof("About to build QueryMessage")
	query := events.QueryEvent{}
	fieldsReceived := make(map[FieldType]bool)
	for i := 0; i < len(message); {
		switch fieldType := message[i]; fieldType {
		case byte(METRICID):
			logrus.Infof("About to parse MetricId field")
			metricId, newIndex, err := parseMetricId(message[i+1:], i+1)
			logrus.Infof("MetricId parsed: %s", metricId)
			if err != nil {
				return events.QueryEvent{}, errors.Wrapf(err, "Could not parse MetricId when trying to build QueryMessage from message received")
			}
			query.MetricId = metricId
			i = newIndex
		case byte(AGGREGATION):
			logrus.Infof("About to parse Aggregation field")
			aggregation, newIndex, err := parseAggregation(message[i+1:], i+1)
			logrus.Infof("Aggregation parsed: %d", aggregation)
			if err != nil {
				return events.QueryEvent{}, errors.Wrapf(err, "Could not parse Aggregation when trying to build QueryMessage from message received")
			}
			query.Aggregation = aggregation
			i = newIndex
		case byte(AGGREGATIONWINDOWSSECS):
			logrus.Infof("About to parse AggregationWindowsSecs field")
			aggregationWindowsSecs, newIndex, err := parseAggregationWindowsSecs(message[i+1:], i+1)
			logrus.Infof("Aggregation parsed: %d", aggregationWindowsSecs)
			if err != nil {
				return events.QueryEvent{}, errors.Wrapf(err, "Could not parse AggregationWindowsSecs when trying to build QueryMessage from message received")
			}
			query.AggregationWindowsSecs = aggregationWindowsSecs
			i = newIndex
		case byte(DATEINTERVAL):
			logrus.Infof("About to parse DateInterval field")
			from, to, newIndex, err := parseDateInterval(message[i+1:], i+1)
			logrus.Infof("From parsed: %s. To parsed: %s", from, to)
			if err != nil {
				return events.QueryEvent{}, errors.Wrapf(err, "Could not parse Date interval when trying to build QueryMessage from message received")
			}
			query.From = from
			query.To = to
			i = newIndex
		default:
			err := fmt.Sprintf("Invalid message format for MetricEvent. Unrecognized Field type %d", fieldType)
			return events.QueryEvent{}, InvalidMessageFormatError{errorMsg: err}
		}
	}
	areAllFieldsPresent := fieldsReceived[METRICID] && fieldsReceived[AGGREGATION] && fieldsReceived[AGGREGATIONWINDOWSSECS]
	if !areAllFieldsPresent {
		return events.QueryEvent{}, InvalidMessageFormatError{errorMsg: "Missing required fields. MetricID, Aggregation, AggregationWindowsSecs are needed"}
	}
	return query, nil
}
