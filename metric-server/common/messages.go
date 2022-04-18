package common

import "time"

type Message interface {
	GetType() MsgType
}

type MsgType uint8

const (
	METRIC MsgType = iota
	QUERY
)

type MetricMessage struct {
	MetricId string
	Value    float32
}

func (m MetricMessage) GetType() MsgType {
	return METRIC
}

type QueryMessage struct {
	MetricId               string
	Aggregation            AggregationType
	AggregationWindowsSecs float32
	Limit                  float32
	From                   time.Time
	To                     time.Time
}

func (m QueryMessage) GetType() MsgType {
	return QUERY
}

type AggregationType uint8

const (
	AVG AggregationType = iota
	MIN
	MAX
	COUNT
)
