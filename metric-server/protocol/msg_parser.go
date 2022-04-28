package protocol

import (
	"encoding/binary"
	"math"
	"net"

	"github.com/manudelca/tp1-distribuidos1/metric-server/events"
	"github.com/manudelca/tp1-distribuidos1/metric-server/util"
	"github.com/pkg/errors"
)

func sendMessage(message Message, clientConn net.Conn) error {
	var bytesToBeSent []byte
	bytesToBeSent = append(bytesToBeSent, byte(message.MsgType))
	bytesToBeSent = append(bytesToBeSent, []byte(message.Message)[:]...)
	bytesLen := uint16(len(bytesToBeSent))
	var bytesLenArray [2]byte
	binary.BigEndian.PutUint16(bytesLenArray[:], bytesLen)
	bytesToBeSent = append(bytesLenArray[:], bytesToBeSent...)

	return util.SendToConnection(clientConn, bytesToBeSent)
}

func SendServerError(errorMsg string, clientConn net.Conn) error {
	message := Message{MsgType: ServerError, Message: errorMsg}
	err := sendMessage(message, clientConn)
	if err != nil {
		return errors.Wrapf(err, "Could not send server error")
	}
	return nil
}

func SendSuccess(message string, clientConn net.Conn) error {
	messageToSend := Message{MsgType: Success, Message: message}
	err := sendMessage(messageToSend, clientConn)
	if err != nil {
		return errors.Wrapf(err, "Could not send success message")
	}
	return nil
}

func SendQueryResult(queryResult events.QueryResultEvent, clientConn net.Conn) error {
	if len(queryResult.Results) == 0 {
		return SendServerError("MetricID not found", clientConn)
	}
	var bytesToSend []byte
	msgType := SuccessQueryResult
	bytesToSend = append(bytesToSend, byte(msgType))

	amountOfIntervalResults := uint32(len(queryResult.Results))
	amountOfIntervalResultsBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(amountOfIntervalResultsBytes, amountOfIntervalResults)
	bytesToSend = append(bytesToSend, amountOfIntervalResultsBytes...)

	for _, queryIntervalResult := range queryResult.Results {

		fromDateBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(fromDateBytes, uint64(queryIntervalResult.FromDate))
		bytesToSend = append(bytesToSend, fromDateBytes...)

		toDateBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(toDateBytes, uint64(queryIntervalResult.ToDate))
		bytesToSend = append(bytesToSend, toDateBytes...)

		resultBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(resultBytes, math.Float32bits(queryIntervalResult.Result))
		bytesToSend = append(bytesToSend, resultBytes...)
	}

	return util.SendToConnection(clientConn, bytesToSend)
}
