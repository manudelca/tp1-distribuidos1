package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"net"

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
	sendMetric(1, conn)
	readServerResponse(conn)
	conn.Close()
}

func sendMetric(msgID int, conn net.Conn) {
	logrus.Infof("[CLIENT] About to send Metric event")
	metric := fmt.Sprintf("metric_%d", msgID)
	metricBytes := []byte(metric)
	metricLen := uint8(len(metric))
	metricFieldType := uint8(0)
	value := float32(1.2)
	var valueBytes [4]byte
	binary.BigEndian.PutUint32(valueBytes[:], math.Float32bits(value))
	valueFieldType := uint8(1)
	eventType := uint8(0)
	msgLen := uint16(1 + 1 + 1 + len(metricBytes) + 1 + 4)
	var msgLenBytes [2]byte
	binary.BigEndian.PutUint16(msgLenBytes[:], msgLen)

	var bytes []byte
	bytes = append(bytes, msgLenBytes[:]...)
	bytes = append(bytes, eventType)
	bytes = append(bytes, metricFieldType)
	bytes = append(bytes, metricLen)
	bytes = append(bytes, metricBytes[:]...)
	bytes = append(bytes, valueFieldType)
	bytes = append(bytes, valueBytes[:]...)
	logrus.Infof("[CLIENT] Bytes to be sent", bytes)
	logrus.Infof("[CLIENT] Finished")
	conn.Write(bytes)
}

func readServerResponse(conn net.Conn) {
	logrus.Infof("[CLIENT] Reading server response")

	var bytes [2]byte
	_, err := conn.Read(bytes[:])
	if err != nil {
		logrus.Infof("[CLIENT] Error reading: ", err)
	}
	len := binary.BigEndian.Uint16(bytes[:])
	msg := make([]byte, len)
	_, err = conn.Read(msg)
	if err != nil {
		logrus.Infof("[CLIENT] Error reading: ", err)
	}
	logrus.Infof("[CLIENT] Msg Len: %d", len)
	logrus.Infof("[CLIENT] Msg Type: %d", msg[0])
	logrus.Infof("[CLIENT] Msg: %s", msg[1:])
}
