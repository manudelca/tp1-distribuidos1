package common

import "time"

type Event interface {
	GetType() EventType
}

type EventType uint8

const (
	METRIC EventType = iota
	QUERY
)

type MetricEvent struct {
	MetricId string
	Value    float32
}

func (m MetricEvent) GetType() EventType {
	return METRIC
}

func (m MetricEvent) Validate() (MetricEvent, error) {
	if m.MetricId == "" {
		return MetricEvent{}, InvalidMetricEventError{errorMsg: "Missing metricId"}
	}
	if m.Value == 0 { //TODO: Find another way of checking if the Value was set...
		return MetricEvent{}, InvalidMetricEventError{errorMsg: "Missing value"}
	}
	return m, nil
}

type QueryEvent struct {
	MetricId               string
	Aggregation            AggregationType
	AggregationWindowsSecs float32
	From                   time.Time
	To                     time.Time
}

func (m QueryEvent) GetType() EventType {
	return QUERY
}

type AggregationType uint8

const (
	AVG AggregationType = iota
	MIN
	MAX
	COUNT
)

type InvalidMetricEventError struct {
	errorMsg string
}

func (e InvalidMetricEventError) Error() string {
	return e.errorMsg
}
