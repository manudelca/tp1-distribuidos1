package storage_protocol

import (
	"strconv"
	"strings"
	"time"

	"github.com/manudelca/tp1-distribuidos1/metric-server/events"
	"github.com/pkg/errors"
)

func ParseLine(line string) (events.MetricEvent, error) {
	lineSplit := strings.Split(line[:len(line)-1], " ")
	if len(lineSplit) < 3 {
		return events.MetricEvent{}, errors.New("Could not parse line. Length of splitted line is less than 3")
	}
	date, err := time.Parse("2006-01-02 03:04:05", lineSplit[0])
	if err != nil {
		return events.MetricEvent{}, errors.Wrapf(err, "Could not parse Date from line")
	}
	metricId := lineSplit[1]
	value, err := strconv.ParseFloat(lineSplit[2], 32)
	if err != nil {
		return events.MetricEvent{}, errors.Wrapf(err, "Could not parse Value from line")
	}
	return events.MetricEvent{MetricId: metricId, Date: date, Value: float32(value)}, nil
}
