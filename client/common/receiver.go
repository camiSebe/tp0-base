package common

import (
	"net"
)

const SIZE_UINT32 = 4

// ReceiveAll receives bytesToRead bytes of data through the connection
func ReceiveAll(conn net.Conn, data []byte, bytesToRead uint32) error {
	var total uint32 = 0
	for total < bytesToRead {
		n, err := conn.Read(data[total:])
		if err != nil {
			log.Criticalf("action: ReceiveAll | result: fail | error: Only %v of %v bytes were read", total, bytesToRead)
			return err
		}
		total += uint32(n)
	}
	log.Infof("action: ReceiveAll | result: success | Bytes read: %v", total)
	return nil
}