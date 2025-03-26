package common

import (
	"net"
)

// ReceiveAll receives all the data through the connection
func ReceiveAll(conn net.Conn, data []byte) error {
	var total = 0
	for total < len(data) {
		n, err := conn.Read(data[total:])
		if err != nil {
			log.Criticalf("action: ReceiveAll | result: fail | error: Only %v of %v bytes were received", total, len(data))
			return err
		}
		total += n
	}
	log.Debugf("action: ReceiveAll | result: success | Bytes received: %v | Data: %v", total, data)
	return nil
}