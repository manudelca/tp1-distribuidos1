package common

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/manudelca/tp1-distribuidos1/metric-server/events"
	"github.com/manudelca/tp1-distribuidos1/metric-server/file_monitor"
	"github.com/manudelca/tp1-distribuidos1/metric-server/storage_protocol"
	"github.com/sirupsen/logrus"
)

type QueryEventsWorker struct {
	queryEventsQueue chan events.Event
	fileMonitor      *file_monitor.FileMonitor
}

func NewQueryEventsWorker(queryEventsQueue chan events.Event, fileMonitor *file_monitor.FileMonitor) *QueryEventsWorker {
	return &QueryEventsWorker{
		queryEventsQueue: queryEventsQueue,
		fileMonitor:      fileMonitor,
	}
}

func (q *QueryEventsWorker) processFile(fileName string, metricId string, left time.Time, right time.Time) (float32, float32, float32, float32, error) {
	logrus.Infof("[QUERY EVENTS WORKER] About to process file: %s", fileName)
	count := float32(0)
	min := float32(math.MaxFloat32)
	max := float32(-math.MaxFloat32)
	sum := float32(0)
	for i := 0; true; i++ {
		bytes, err := q.fileMonitor.ReadMetric(fileName, i)
		if err != nil {
			return 0, 0, 0, 0, err
		}
		if len(bytes) == 0 {
			break
		}
		metricEvent, err := storage_protocol.ParseBytesToMetric(bytes, metricId)
		if err != nil {
			logrus.Infof("[QUERY EVENTS WORKER] Failed to parse bytes: %s from file: %s. Error: %s", bytes, fileName, err.Error())
			continue
		}
		if metricEvent.Date < left.Unix() || metricEvent.Date > right.Unix() {
			continue
		}
		count += 1
		sum += metricEvent.Value
		if metricEvent.Value > max {
			max = metricEvent.Value
		}
		if metricEvent.Value < min {
			min = metricEvent.Value
		}
	}
	return count, min, max, sum, nil
}

func (q *QueryEventsWorker) processTimeInterval(metricId string, aggregationType events.AggregationType, left time.Time, right time.Time) (float32, float32, float32, float32, error) {
	count := float32(0)
	min := float32(math.MaxFloat32)
	max := float32(-math.MaxFloat32)
	sum := float32(0)
	for currentLeftTime := left.Truncate(time.Minute); currentLeftTime.Before(right); currentLeftTime = currentLeftTime.Add(time.Minute) {
		fileName := storage_protocol.GetFileName(metricId, currentLeftTime.Unix())
		if q.fileMonitor.FileExists(fileName) {
			countFile, minFile, maxFile, sumFile, err := q.processFile(fileName, metricId, currentLeftTime, right)
			if err != nil {
				logrus.Infof("[QUERY EVENTS WORKER] Failed to read line from file: %s. Error: %s", fileName, err.Error())
				return 0, 0, 0, 0, err
			}
			count += countFile
			sum += sumFile
			if maxFile > max {
				max = maxFile
			}
			if minFile < min {
				min = minFile
			}
		} else {
			// Log too noisy
			// logrus.Infof("[QUERY EVENTS WORKER] File not found: %s", fileName)
		}
	}
	if count == 0 {
		// Border case
		return 0, 0, 0, 0, errors.New(fmt.Sprintf("No file was processed in the time interval %s, %s", left, right))
	}
	return count, min, max, sum / count, nil
}

func (q *QueryEventsWorker) handleQueryEvent(queryEvent events.QueryEvent) events.QueryResultEvent {
	logrus.Infof("[QUERY EVENTS WORKER] Processing query event: ", queryEvent)
	result := make([]events.QueryIntervalResult, 0)
	fromDate := time.Unix(queryEvent.FromDate, 0)
	toDate := time.Unix(queryEvent.ToDate, 0)
	timeWindow := time.Duration(queryEvent.AggregationWindowsSecs * 1e9)
	logrus.Infof("[QUERY EVENTS WORKER] While parse using From date: %s, To Date: %s, Window secs: %s", fromDate, toDate, timeWindow)
	for leftWindowLimit := fromDate; leftWindowLimit.Before(toDate); {
		rightWindowLimit := leftWindowLimit.Add(timeWindow)
		if rightWindowLimit.After(toDate) {
			rightWindowLimit = toDate
		}
		resultInterval := float32(0)
		count, min, max, avg, err := q.processTimeInterval(queryEvent.MetricId, queryEvent.Aggregation, leftWindowLimit, rightWindowLimit)
		if err != nil {
			// Log too noisy
			// logrus.Infof("[QUERY EVENTS WORKER] Error processing time interval From Date: %s, To Date: %s", leftWindowLimit, rightWindowLimit)
		} else {
			switch queryEvent.Aggregation {
			case events.COUNT:
				resultInterval = count
			case events.MIN:
				resultInterval = min
			case events.MAX:
				resultInterval = max
			case events.AVG:
				resultInterval = avg
			}
			queryResult := events.QueryIntervalResult{
				FromDate: leftWindowLimit.Unix(),
				ToDate:   rightWindowLimit.Unix(),
				Result:   resultInterval,
			}
			result = append(result, queryResult)
		}
		leftWindowLimit = rightWindowLimit
	}
	logrus.Infof("[QUERY EVENTS WORKER] Finished processing every time window")
	logrus.Infof("[QUERY EVENTS WORKER] Result len: %d", len(result))
	return events.QueryResultEvent{Results: result}
}

func (q *QueryEventsWorker) ServeQueryEvents() {
	for queryEvent := range q.queryEventsQueue {
		switch eventType := queryEvent.GetType(); eventType {
		case events.QUERY:
			query, ok := queryEvent.(events.QueryEvent)
			if !ok {
				logrus.Infof("[QUERY EVENTS WORKER] Could not assert event to events.QueryEvent")
				continue
			}
			result := q.handleQueryEvent(query)
			q.queryEventsQueue <- result
		}
	}
}
