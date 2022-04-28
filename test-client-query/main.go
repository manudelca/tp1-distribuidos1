package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	conn, err := net.Dial("tcp", "metric-server:12345")
	if err != nil {
		log.Fatalf(
			"[CLIENT] Could not connect to server. Error: %v",
			err,
		)
	}
	sendQuery(1, conn)
	readServerResponse(conn)
	conn.Close()
}

func sendQuery(msgID int, conn net.Conn) {
	logrus.Infof("[CLIENT QUERY] About to send Query event")
	eventType := uint8(1)

	metricFieldType := uint8(0)
	metric := fmt.Sprintf("metric_%d", msgID)
	metricBytes := []byte(metric)
	metricLen := uint8(len(metric))

	aggregationFieldType := uint8(2)
	aggregation := uint8(0)

	aggregationWindowsSecsFieldType := uint8(3)
	aggregationWindowsSecs := float32(1.2)
	var aggregationWindowsSecsBytes [4]byte
	binary.BigEndian.PutUint32(aggregationWindowsSecsBytes[:], math.Float32bits(aggregationWindowsSecs))

	dateTimeIntervalFieldType := uint8(4)
	layout := "2006-01-02 15:04:05"
	fromDate, err := time.Parse(layout, "2022-04-28 00:15:00")
	if err != nil {
		logrus.Infof("[CLIENT QUERY] Error parsing date %s", err.Error())
	}
	from := fromDate.Unix()
	fromBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(fromBytes, uint64(from))
	toDate, err := time.Parse(layout, "2022-04-28 00:20:00")
	if err != nil {
		logrus.Infof("[CLIENT QUERY] Error parsing date %s", err.Error())
	}
	to := toDate.Unix()
	toBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(toBytes, uint64(to))
	logrus.Infof("fromBytes ", fromBytes)
	logrus.Infof("toBytes ", toBytes)

	var bytes []byte

	bytes = append(bytes, eventType)
	bytes = append(bytes, metricFieldType)
	bytes = append(bytes, metricLen)
	bytes = append(bytes, metricBytes[:]...)
	bytes = append(bytes, aggregationFieldType)
	bytes = append(bytes, aggregation)
	bytes = append(bytes, aggregationWindowsSecsFieldType)
	bytes = append(bytes, aggregationWindowsSecsBytes[:]...)
	bytes = append(bytes, dateTimeIntervalFieldType)
	bytes = append(bytes, fromBytes...)
	bytes = append(bytes, toBytes...)

	msgLen := uint16(len(bytes))
	var msgLenBytes [2]byte
	binary.BigEndian.PutUint16(msgLenBytes[:], msgLen)

	bytes = append(msgLenBytes[:], bytes[:]...)

	a := time.Minute / 5
	time.Sleep(a)

	logrus.Infof("[CLIENT QUERY] Bytes to be sent", bytes)
	logrus.Infof("[CLIENT QUERY] Finished")
	conn.Write(bytes)
}

func readServerResponse(conn net.Conn) {
	logrus.Infof("[CLIENT QUERY] Reading server response")

	var bytes [2]byte
	_, err := conn.Read(bytes[:])
	if err != nil {
		logrus.Infof("[CLIENT QUERY] Error reading: ", err)
	}
	len := binary.BigEndian.Uint16(bytes[:])
	msg := make([]byte, len)
	_, err = conn.Read(msg)
	if err != nil {
		logrus.Infof("[CLIENT QUERY] Error reading: ", err)
	}
	logrus.Infof("[CLIENT QUERY] Msg Len: %d", len)
	logrus.Infof("[CLIENT QUERY] Msg Type: %d", msg[0])

	if msg[0] == 1 {
		amountOfIntervalResults := binary.BigEndian.Uint32(msg[1:5])
		logrus.Infof("[CLIENT QUERY] Msg: %d", amountOfIntervalResults)

		intervals := msg[5:]
		for i := uint32(0); i < amountOfIntervalResults; i++ {
			currentInterval := intervals[i : 8+8+4]
			fromDate := int64(binary.BigEndian.Uint64(currentInterval[:8]))
			toDate := int64(binary.BigEndian.Uint64(currentInterval[8:16]))
			result := float32(binary.BigEndian.Uint32(currentInterval[16:20]))
			logrus.Infof("[CLIENT QUERY] Parsed interval (from, to, result): %d, %d, %f", fromDate, toDate, result)
		}
	} else {
		logrus.Infof("[CLIENT QUERY] Msg: %s", msg[1:])
	}
}
