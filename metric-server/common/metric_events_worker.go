package common

import (
	"github.com/manudelca/tp1-distribuidos1/metric-server/events"
	"github.com/sirupsen/logrus"
)

func handleMetricEvent(metricEvent events.MetricEvent) {
	logrus.Infof("[METRIC EVENTS WORKER] Processing metric event: ", metricEvent)
}

func ServeMetricEvents(metricEventsToServe chan events.MetricEvent) {
	for metricEvent := range metricEventsToServe {
		handleMetricEvent(metricEvent)
	}
}
