package common

import (
	"bufio"
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
		c.sendMetric(msgID)
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
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

func (c *Client) sendMetric(msgID int) {
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
	logrus.Infof("MsgLen: %s", msgLen)
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
	logrus.Infof("Bytes to be sent", bytes)
	logrus.Infof("Finished")
	c.conn.Write(bytes)
}
