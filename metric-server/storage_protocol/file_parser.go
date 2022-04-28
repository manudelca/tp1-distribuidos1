package storage_protocol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/manudelca/tp1-distribuidos1/metric-server/events"
)

func ParseBytesToMetric(bytes []byte, metricId string) (events.MetricEvent, error) {
	if len(bytes) < (8 + 4) {
		return events.MetricEvent{}, errors.New("Bytes array is to short. 12 bytes needed (int64 and float32)")
	}
	date := int64(binary.BigEndian.Uint64(bytes[:8]))
	value := math.Float32frombits(binary.BigEndian.Uint32(bytes[8:]))
	return events.MetricEvent{MetricId: metricId, Date: date, Value: value}, nil
}

func ParseMetrictToBytes(metricEvent events.MetricEvent) []byte {
	bytes := make([]byte, 8+4)
	binary.BigEndian.PutUint64(bytes[:8], uint64(metricEvent.Date))
	binary.BigEndian.PutUint32(bytes[8:], math.Float32bits(metricEvent.Value))
	return bytes

}

func GetFileName(metricId string, date int64) string {
	timeParsed := time.Unix(date, 0)
	year, month, day := timeParsed.Date()
	hours, minutes, _ := timeParsed.Clock()
	return fmt.Sprintf("%s_%d%02d%02d_%02d%02d", metricId, year, month, day, hours, minutes)
}
