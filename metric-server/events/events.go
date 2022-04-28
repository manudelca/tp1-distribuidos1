package events

type Event interface {
	GetType() EventType
}

type EventType uint8

const (
	METRIC EventType = iota
	QUERY
	ALERT
	QUERYRESULT
)

type MetricEvent struct {
	MetricId string
	Value    float32
	Date     int64
}

func (m MetricEvent) GetType() EventType {
	return METRIC
}

type QueryEvent struct {
	MetricId               string
	Aggregation            AggregationType
	AggregationWindowsSecs float32
	FromDate               int64
	ToDate                 int64
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

type AlertEvent struct {
	MetricId               string
	Aggregation            AggregationType
	Limit                  float32
	AggregationWindowsSecs float32
}

func (a AlertEvent) GetType() EventType {
	return ALERT
}

type QueryResultEvent struct {
	Results []QueryIntervalResult
}

func (a QueryResultEvent) GetType() EventType {
	return QUERYRESULT
}

type QueryIntervalResult struct {
	FromDate int64
	ToDate   int64
	Result   float32
}
