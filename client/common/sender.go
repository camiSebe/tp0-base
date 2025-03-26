package common

import (
	"net"
)

// SendAll sends all the data through the connection
func SendAll(conn net.Conn, data []byte) error {
	var total = 0
	for total < len(data) {
		n, err := conn.Write(data[total:])
		if err != nil {
			log.Criticalf("action: SendAll | result: fail | error: Only %v of %v bytes were sent", total, len(data))
			return err
		}
		total += n
	}
	log.Debugf("action: SendAll | result: success | Bytes sent: %v | Data: %v", total, data)
	return nil
}
