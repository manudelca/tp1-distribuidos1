package util

import "net"

func ReadFromConnection(conn net.Conn, bytesToRead int) ([]byte, error) {
	bytes := make([]byte, bytesToRead)
	amountRead := 0
	for amountRead < bytesToRead {
		bytesRead, err := conn.Read(bytes)
		if err != nil {
			return nil, err
		}
		amountRead = bytesRead
	}
	return bytes, nil
}
