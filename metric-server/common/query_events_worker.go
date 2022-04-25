package common

import (
	"time"

	"github.com/manudelca/tp1-distribuidos1/metric-server/events"
	"github.com/manudelca/tp1-distribuidos1/metric-server/file_monitor"
	"github.com/sirupsen/logrus"
)

type QueryEventsWorker struct {
	queryEventsQueue chan events.QueryEvent
	fileMonitor      *file_monitor.FileMonitor
}

func NewQueryEventsWorker(queryEventsQueue chan events.QueryEvent, fileMonitor *file_monitor.FileMonitor) *QueryEventsWorker {
	return &QueryEventsWorker{
		queryEventsQueue: queryEventsQueue,
		fileMonitor:      fileMonitor,
	}
}

func (q *QueryEventsWorker) handleQueryEvent(queryEvent events.QueryEvent) {
	logrus.Infof("[QUERY EVENTS WORKER] Processing query event: ", queryEvent)
	fromDate := queryEvent.From
	toDate := queryEvent.To
	timeWindow := time.Duration(queryEvent.AggregationWindowsSecs * 1e9)
	logrus.Infof("[QUERY EVENTS WORKER] While parse using From date: %s, To Date: %s, Window secs: %s", fromDate, toDate, timeWindow)
	for leftWindowLimit := fromDate; leftWindowLimit.Before(toDate); {
		rightWindowLimit := leftWindowLimit.Add(timeWindow)
		if rightWindowLimit.After(toDate) {
			rightWindowLimit = toDate
		}
		leftWindowLimit = rightWindowLimit
	}
	logrus.Infof("[QUERY EVENTS WORKER] Finished processing every time window")
}

func (q *QueryEventsWorker) ServeQueryEvents() {
	for queryEvent := range q.queryEventsQueue {
		q.handleQueryEvent(queryEvent)
	}
}
