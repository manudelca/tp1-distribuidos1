package protocol

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"

	"github.com/pkg/errors"
)

type FieldType uint8

const (
	METRICID FieldType = iota
	VALUE
	AGGREGATION
	AGGREGATIONWINDOWSSECS
	FROM
	TO
)

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
	if len(message) < 4 {
		errorMsg := fmt.Sprintf("The message is shorter than 4 bytes")
		return 0, i, InvalidValueFieldError{errorMsg: errorMsg}
	}
	valueBytes := message[:4]
	if _, err := strconv.ParseUint(string(valueBytes), 10, 16); err != nil {
		return 0, i, errors.Wrapf(err, "Value could not be parsed as a float32")
	}
	value := math.Float32frombits(binary.BigEndian.Uint32(message[:4]))
	return value, i + 4, nil
}
