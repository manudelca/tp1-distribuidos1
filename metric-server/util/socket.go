package util

import "net"

func ReadFromConnection(conn net.Conn, bytes []byte, bytesToRead int) error {
	amountRead := 0
	for amountRead < bytesToRead {
		bytesRead, err := conn.Read(bytes)
		if err != nil {
			return err
		}
		amountRead = bytesRead
	}
	return nil
}
