package protocol

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/manudelca/tp1-distribuidos1/metric-server/common"
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
	valueBytes := message[:4]
	if _, err := strconv.ParseUint(string(valueBytes), 10, 16); err != nil {
		return 0, i, errors.Wrapf(err, "Could not be parsed as a float32")
	}
	value := math.Float32frombits(binary.BigEndian.Uint32(message[:4]))
	return value, i + 4, nil
}

func parseDate(message []byte, i int) (time.Time, int, error) {
	layout := "2006-01-02 03:04:05"
	if len(message) < len(layout) {
		errorMsg := fmt.Sprintf("The message is shorter than the layout proposed (YYYY-MM-DD HH:mm:ss)")
		return time.Time{}, i, InvalidDateFieldError{errorMsg: errorMsg}
	}
	date, err := time.Parse(layout, string(message[:len(layout)]))
	if err != nil {
		return time.Time{}, i, errors.Wrapf(err, "Could not parse date using layout YYYY-MM-DD HH:mm:ss")
	}
	return date, i + len(layout), nil
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

func parseAggregation(message []byte, i int) (common.AggregationType, int, error) {
	if len(message) < 1 {
		errorMsg := fmt.Sprintf("The message is shorter than 1 bytes")
		return 0, i, InvalidAggregationFieldError{errorMsg: errorMsg}
	}
	aggregation := common.AggregationType(message[0])
	return aggregation, i + 1, nil
}

func parseAggregationWindowsSecs(message []byte, i int) (float32, int, error) {
	aggregationWindowsSecs, i, err := parseFloat(message, i)
	if err != nil {
		return 0, i, errors.Wrapf(err, "Could not parse aggregationWindowsSecs")
	}
	return aggregationWindowsSecs, i, nil
}

func parseFrom(message []byte, i int) (time.Time, int, error) {
	date, i, err := parseDate(message, i)
	if err != nil {
		return date, i, errors.Wrapf(err, "Could not parse From date")
	}
	return date, i, nil
}

func parseTo(message []byte, i int) (time.Time, int, error) {
	date, i, err := parseDate(message, i)
	if err != nil {
		return date, i, errors.Wrapf(err, "Could not parse To date")
	}
	return date, i, nil
}

func parseDateInterval(message []byte, i int) (time.Time, time.Time, int, error) {
	from, newIndex, err := parseDate(message, i)
	if err != nil {
		return from, from, i, errors.Wrapf(err, "Could not parse From date")
	}
	to, newIndex, err := parseDate(message, newIndex)
	if err != nil {
		return from, to, i, errors.Wrapf(err, "Could not parse To date")
	}
	return from, to, newIndex, nil
}
