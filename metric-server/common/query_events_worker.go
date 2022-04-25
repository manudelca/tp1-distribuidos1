package common

import (
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
}

func (q *QueryEventsWorker) ServeQueryEvents() {
	for queryEvent := range q.queryEventsQueue {
		q.handleQueryEvent(queryEvent)
	}
}
