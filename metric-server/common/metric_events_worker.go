package common

import (
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

func (m *MetricEventsWorker) handleMetricEvent(metricEvent events.MetricEvent) {
	logrus.Infof("[METRIC EVENTS WORKER] Processing metric event: ", metricEvent)
}

func (m *MetricEventsWorker) ServeMetricEvents() {
	for metricEvent := range m.metricEventsQueue {
		m.handleMetricEvent(metricEvent)
	}
}
