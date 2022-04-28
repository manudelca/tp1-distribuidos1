package common

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/manudelca/tp1-distribuidos1/metric-server/events"
	"github.com/sirupsen/logrus"
)

type AlertEventsWorker struct {
	eventsQueue    chan events.Event
	orderedMetrics []events.MetricEvent
	fileName       string
}

func NewAlertEventsWorker(eventsQueue chan events.Event, fileName string) *AlertEventsWorker {
	return &AlertEventsWorker{
		eventsQueue:    eventsQueue,
		orderedMetrics: make([]events.MetricEvent, 0),
		fileName:       fileName,
	}
}

func (a *AlertEventsWorker) ServeAlertEvents() {
	for event := range a.eventsQueue {
		logrus.Infof("[ALERT EVENTS WORKER] Processing event: ", event)
		switch eventType := event.GetType(); eventType {
		case events.ALERT:
			alert, ok := event.(events.AlertEvent)
			if !ok {
				logrus.Infof("[ALERT EVENTS WORKER] Could not assert event to events.AlertEvent")
				continue
			}
			logrus.Infof("[ALERT EVENTS WORKER] AlertEvent received: ", alert)
			a.handleAlert(alert)
			logrus.Infof("[ALERT EVENTS WORKER] Finished processing AlertEvent. Cleaning metrics")
			a.orderedMetrics = make([]events.MetricEvent, 0)
		case events.METRIC:
			metric, ok := event.(events.MetricEvent)
			if !ok {
				logrus.Infof("[ALERT EVENTS WORKER] Could not assert event to events.MetricEvent")
				continue
			}
			logrus.Infof("[ALERT EVENTS WORKER] MetricEvent received: ", metric)
			a.insertIntoMetrics(metric)
		}
	}
}

func (a *AlertEventsWorker) handleAlert(alert events.AlertEvent) {
	for i, metric := range a.orderedMetrics {
		left := time.Unix(metric.Date, 0)
		right := time.Unix(metric.Date, 0).Add(time.Duration(alert.AggregationWindowsSecs * 1e9))
		a.processInterval(alert, a.orderedMetrics[i:], left, right)
	}
}

func (a *AlertEventsWorker) insertIntoMetrics(metric events.MetricEvent) {
	insertIndex := len(a.orderedMetrics)
	for i, currentMetric := range a.orderedMetrics {
		if currentMetric.Date > metric.Date {
			insertIndex = i
			break
		}
	}
	if insertIndex == len(a.orderedMetrics) {
		a.orderedMetrics = append(a.orderedMetrics, metric)
	}
	a.orderedMetrics = append(a.orderedMetrics[:insertIndex+1], a.orderedMetrics[insertIndex:]...)
	a.orderedMetrics[insertIndex] = metric
}

func (a *AlertEventsWorker) alert(aggregation string, result float32, left time.Time, right time.Time) {
	logrus.Infof("[ALERT EVENTS WORKER] %d - %d %s exceeded %f", left, right, aggregation, result)
	a.writeAlertOnFile(fmt.Sprintf("%d %d %s %f\n", left, right, aggregation, result))
}

func (a *AlertEventsWorker) writeAlertOnFile(line string) {
	file, err := os.OpenFile(a.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Infof("[ALERT EVENTS WORKER] Failed to open file %s. Error: %s", a.fileName, err.Error())
		return
	}
	for i := 0; i < len(line); {
		written, err := file.WriteString(line[i:])
		if err != nil {
			logrus.Infof("[ALERT EVENTS WORKER] Failed to write alert on file %s. Error: %s", a.fileName, err.Error())
			return
		}
		i = i + written
	}
	file.Close()
}

func (a *AlertEventsWorker) processInterval(alert events.AlertEvent, metrics []events.MetricEvent, left time.Time, right time.Time) {
	count := float32(0)
	min := float32(math.MaxFloat32)
	max := float32(-math.MaxFloat32)
	sum := float32(0)

	for _, metric := range a.orderedMetrics {
		if metric.MetricId != alert.MetricId {
			continue
		}
		value := metric.Value
		count += 1
		sum += value
		if value > max {
			max = value
		}
		if value < min {
			min = value
		}
	}
	aggregationStr := ""
	result := float32(0)
	switch aggregation := alert.Aggregation; aggregation {
	case events.AVG:
		result = sum / count
		aggregationStr = "AVG"
	case events.MIN:
		result = min
		aggregationStr = "MIN"
	case events.MAX:
		result = max
		aggregationStr = "MAX"
	case events.COUNT:
		result = count
		aggregationStr = "COUNT"
	}
	if alert.Limit < result {
		a.alert(aggregationStr, result, left, right)
	}
}
