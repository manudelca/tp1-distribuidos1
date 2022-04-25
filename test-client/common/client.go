package common

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopLapse     time.Duration
	LoopPeriod    time.Duration
	ClientType    string
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Fatalf(
			"[CLIENT %v] Could not connect to server. Error: %v",
			c.config.ID,
			err,
		)
	}
	c.conn = conn
	return nil
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	// Create the connection the server in every loop iteration. Send an
	// autoincremental msgID to identify every message sent
	c.createClientSocket()
	msgID := 1

loop:
	// Send messages if the loopLapse threshold has been not surpassed
	for timeout := time.After(c.config.LoopLapse); ; {
		select {
		case <-timeout:
			break loop
		default:
		}

		// Send
		if c.config.ClientType == "METRIC" {
			c.sendMetric(msgID)
		} else if c.config.ClientType == "QUERY" {
			c.sendQuery(msgID)
		} else {
			logrus.Infof("Unrecognized clientType: %s", c.config.ClientType)
			c.conn.Close()
			return
		}
		msg, err := c.readServerResponse()
		msgID++

		if err != nil {
			log.Errorf(
				"[CLIENT %v] Error reading from socket. %v.",
				c.config.ID,
				err,
			)
			c.conn.Close()
			return
		}
		log.Infof("[CLIENT %v] Message from server: %v", c.config.ID, msg)

		// Wait a time between sending one message and the next one
		time.Sleep(c.config.LoopPeriod)

		// Recreate connection to the server
		c.conn.Close()
		c.createClientSocket()
	}

	log.Infof("[CLIENT %v] Closing connection", c.config.ID)
	c.conn.Close()
}

func (c *Client) readServerResponse() (string, error) {
	logrus.Infof("[CLIENT %v] About to send Metric event", c.config.ID)

	var bytes [2]byte
	_, err := c.conn.Read(bytes[:])
	if err != nil {
		return "", err
	}
	len := binary.BigEndian.Uint16(bytes[:])
	msg := make([]byte, len)
	_, err = c.conn.Read(msg)
	if err != nil {
		return "", err
	}
	logrus.Infof("[CLIENT %v] Msg Len: %d", c.config.ID, len)
	logrus.Infof("[CLIENT %v] Msg Type: %d", c.config.ID, msg[0])
	return string(msg[1:]), nil
}

func (c *Client) sendMetric(msgID int) {
	logrus.Infof("[CLIENT %v] About to send Metric event", c.config.ID)
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
	logrus.Infof("[CLIENT %v] Bytes to be sent", c.config.ID, bytes)
	logrus.Infof("[CLIENT %v] Finished", c.config.ID)
	c.conn.Write(bytes)
}

func (c *Client) sendQuery(msgID int) {
	logrus.Infof("[CLIENT %v] About to send Query event", c.config.ID)
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
	from := "2020-01-01 00:00:00"
	to := "2020-01-01 00:05:00"

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
	bytes = append(bytes, []byte(from)[:]...)
	bytes = append(bytes, []byte(to)[:]...)

	msgLen := uint16(len(bytes))
	var msgLenBytes [2]byte
	binary.BigEndian.PutUint16(msgLenBytes[:], msgLen)

	bytes = append(msgLenBytes[:], bytes[:]...)

	logrus.Infof("[CLIENT %v] Bytes to be sent", c.config.ID, bytes)
	logrus.Infof("[CLIENT %v] Finished", c.config.ID)
	c.conn.Write(bytes)
}
