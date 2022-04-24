package common

import (
	"github.com/manudelca/tp1-distribuidos1/metric-server/events"
	"github.com/sirupsen/logrus"
)

func handleQueryEvent(queryEvent events.QueryEvent) {
	logrus.Infof("[QUERY EVENTS WORKER] Processing query event: ", queryEvent)
}

func ServeQueryEvents(queryEventsToServe chan events.QueryEvent) {
	for queryEvent := range queryEventsToServe {
		handleQueryEvent(queryEvent)
	}
}
