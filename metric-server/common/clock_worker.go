package common

import (
	"time"

	"github.com/manudelca/tp1-distribuidos1/metric-server/events"
	"github.com/sirupsen/logrus"
)

type ClockWorker struct {
	eventsQueue chan events.Event
	alertToSend events.AlertEvent
}

func NewClockWorker(eventsQueue chan events.Event, alertToSend events.AlertEvent) *ClockWorker {
	return &ClockWorker{
		eventsQueue: eventsQueue,
		alertToSend: alertToSend,
	}
}

func (c *ClockWorker) Run() {
	for true {
		time.Sleep(time.Minute)
		logrus.Infof("[CLOCK WORKER] Woke up")
		c.eventsQueue <- c.alertToSend
		logrus.Infof("[CLOCK WORKER] Pushed AlertEvent into queue")
	}
}
