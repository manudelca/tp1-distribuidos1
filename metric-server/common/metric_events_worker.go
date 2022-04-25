package common

import (
	"fmt"

	"github.com/manudelca/tp1-distribuidos1/metric-server/events"
	"github.com/manudelca/tp1-distribuidos1/metric-server/file_monitor"
	"github.com/sirupsen/logrus"
)

type MetricEventsWorker struct {
	metricEventsQueue chan events.MetricEvent
	fileMonitor       *file_monitor.FileMonitor
}

func NewMetricEventsWorker(metricEventsQueue chan events.MetricEvent, fileMonitor *file_monitor.FileMonitor) *MetricEventsWorker {
	return &MetricEventsWorker{
		metricEventsQueue: metricEventsQueue,
		fileMonitor:       fileMonitor,
	}
}

func (m *MetricEventsWorker) ServeMetricEvents() {
	for metricEvent := range m.metricEventsQueue {
		logrus.Infof("[METRIC EVENTS WORKER] Processing metric event: ", metricEvent)
		metricToWrite := fmt.Sprintf("%s %s %f", metricEvent.Date, metricEvent.MetricId, metricEvent.Value)
		year, month, day := metricEvent.Date.Date()
		hours, minutes, _ := metricEvent.Date.Clock()
		fileToWrite := fmt.Sprintf("%s_%d%02d%02d_%02d%02d", metricEvent.MetricId, year, month, day, hours, minutes)
		logrus.Infof("[METRIC EVENTS WORKER] Metric to write: \"%s\" On file: \"%s\"", metricToWrite, fileToWrite)
		m.fileMonitor.WriteLineOnFile(metricToWrite, fileToWrite)
	}
}
