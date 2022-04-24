package util

import "net"

func ReadFromConnection(conn net.Conn, bytesToRead int) ([]byte, error) {
	bytes := make([]byte, bytesToRead)
	for amountRead := 0; amountRead < bytesToRead; {
		bytesRead, err := conn.Read(bytes[amountRead:])
		if err != nil {
			return nil, err
		}
		amountRead = amountRead + bytesRead
	}
	return bytes, nil
}

func SendToConnection(conn net.Conn, bytesToSend []byte) error {
	for i := 0; i < len(bytesToSend); {
		bytesSent, err := conn.Write(bytesToSend[i:])
		if err != nil {
			return err
		}
		i = i + bytesSent
	}
	return nil
}
