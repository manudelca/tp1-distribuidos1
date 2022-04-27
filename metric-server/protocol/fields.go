package protocol

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/manudelca/tp1-distribuidos1/metric-server/events"
	"github.com/pkg/errors"
)

type FieldType uint8

const (
	METRICID FieldType = iota
	VALUE
	AGGREGATION
	AGGREGATIONWINDOWSSECS
	DATEINTERVAL
)

func parseFloat(message []byte, i int) (float32, int, error) {
	if len(message) < 4 {
		errorMsg := fmt.Sprintf("The message is shorter than 4 bytes")
		return 0, i, InvalidFloatFieldError{errorMsg: errorMsg}
	}
	value := math.Float32frombits(binary.BigEndian.Uint32(message[:4]))
	return value, i + 4, nil
}

func parseDate(message []byte, i int) (int64, int, error) {
	if len(message) < 8 {
		errorMsg := fmt.Sprintf("The message is shorter than Unix timestamp (int64)")
		return 0, i, InvalidDateFieldError{errorMsg: errorMsg}
	}
	date := int64(binary.BigEndian.Uint64(message[:8]))
	return date, i + 8, nil
}

func parseMetricId(message []byte, i int) (string, int, error) {
	metricIdLen := uint8(message[0])
	messageLen := len(message[1:])
	if int(metricIdLen) > messageLen {
		errorMsg := fmt.Sprintf("Metric len received is longer than the whole message (Len received: %d, message len: %d)", int(metricIdLen), messageLen)
		return "", i, InvalidMetricIdFieldError{errorMsg: errorMsg}
	}
	metricReceived := string(message[1 : metricIdLen+1])
	newIndex := i + int(metricIdLen+1)
	return metricReceived, newIndex, nil
}

func parseValue(message []byte, i int) (float32, int, error) {
	value, i, err := parseFloat(message, i)
	if err != nil {
		return 0, i, errors.Wrapf(err, "Could not parse Value")
	}
	return value, i, nil
}

func parseAggregation(message []byte, i int) (events.AggregationType, int, error) {
	if len(message) < 1 {
		errorMsg := fmt.Sprintf("The message is shorter than 1 bytes")
		return 0, i, InvalidAggregationFieldError{errorMsg: errorMsg}
	}
	aggregation := events.AggregationType(message[0])
	return aggregation, i + 1, nil
}

func parseAggregationWindowsSecs(message []byte, i int) (float32, int, error) {
	aggregationWindowsSecs, i, err := parseFloat(message, i)
	if err != nil {
		return 0, i, errors.Wrapf(err, "Could not parse aggregationWindowsSecs")
	}
	return aggregationWindowsSecs, i, nil
}

func parseDateInterval(message []byte, i int) (int64, int64, int, error) {
	from, newIndex, err := parseDate(message, i)
	if err != nil {
		return from, from, i, errors.Wrapf(err, "Could not parse From date")
	}
	to, newIndex, err := parseDate(message[8:], newIndex)
	if err != nil {
		return from, to, i, errors.Wrapf(err, "Could not parse To date")
	}
	return from, to, newIndex, nil
}
