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

type QueryEvent struct {
	MetricId               string
	Aggregation            AggregationType
	AggregationWindowsSecs float32
	From                   time.Time
	To                     time.Time
}

func (q QueryEvent) GetType() EventType {
	return QUERY
}

func (q QueryEvent) Validate() (QueryEvent, error) {
	return q, nil
}

type AggregationType uint8

const (
	AVG AggregationType = iota
	MIN
	MAX
	COUNT
)
