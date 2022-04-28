package common

import (
	"time"

	"github.com/manudelca/tp1-distribuidos1/metric-server/events"
	"github.com/sirupsen/logrus"
)

type ClockWorker struct {
	eventsQueue          chan events.Event
	alertToSend          events.AlertEvent
	clockShutdownCommand chan bool
	clockIsOff           chan bool
}

func NewClockWorker(eventsQueue chan events.Event, alertToSend events.AlertEvent, clockShutdownCommand chan bool, clockIsOff chan bool) *ClockWorker {
	return &ClockWorker{
		eventsQueue:          eventsQueue,
		alertToSend:          alertToSend,
		clockShutdownCommand: clockShutdownCommand,
		clockIsOff:           clockIsOff,
	}
}

func (c *ClockWorker) Run() {
	for true {
		time.Sleep(time.Minute)
		logrus.Infof("[CLOCK WORKER] Woke up")
		select {
		case <-c.clockShutdownCommand:
			logrus.Infof("[CLOCK WORKER] Proceeding to shutdown")
			break
		default:
			c.eventsQueue <- c.alertToSend
			logrus.Infof("[CLOCK WORKER] Pushed AlertEvent into queue")
		}
	}
	c.clockIsOff <- true
}
