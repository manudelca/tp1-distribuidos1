package common

import (
	"sync"

	"github.com/manudelca/tp1-distribuidos1/metric-server/events"
	"github.com/manudelca/tp1-distribuidos1/metric-server/file_monitor"
	"github.com/manudelca/tp1-distribuidos1/metric-server/storage_protocol"
	"github.com/sirupsen/logrus"
)

type MetricEventsWorker struct {
	metricEventsQueue chan events.MetricEvent
	fileMonitor       *file_monitor.FileMonitor
	wait              *sync.WaitGroup
}

func NewMetricEventsWorker(metricEventsQueue chan events.MetricEvent, fileMonitor *file_monitor.FileMonitor, wait *sync.WaitGroup) *MetricEventsWorker {
	return &MetricEventsWorker{
		metricEventsQueue: metricEventsQueue,
		fileMonitor:       fileMonitor,
		wait:              wait,
	}
}

func (m *MetricEventsWorker) ServeMetricEvents() {
	for metricEvent := range m.metricEventsQueue {
		logrus.Infof("[METRIC EVENTS WORKER] Processing metric event: ", metricEvent)
		metricToWrite := storage_protocol.ParseMetrictToBytes(metricEvent)
		fileToWrite := storage_protocol.GetFileName(metricEvent.MetricId, metricEvent.Date)
		err := m.fileMonitor.WriteMetricOnFile(metricToWrite, fileToWrite)
		if err != nil {
			logrus.Infof("[METRIC EVENTS WORKER] Failed to write metric: \"%s\" On file: \"%s\". Error %s", metricToWrite, fileToWrite, err.Error())
		}
		logrus.Infof("[METRIC EVENTS WORKER] Metric successfully written: \"%s\" On file: \"%s\"", metricToWrite, fileToWrite)
	}
	logrus.Infof("[METRIC EVENTS WORKER] metricEventsQueue closed. Proceeding to shutdown")
	m.wait.Done()
}
