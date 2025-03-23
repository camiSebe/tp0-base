package common

import (
	"net"
	"encoding/binary"
)

// SendSerializedBet sends:
// - len(nombre) in 4 bytes
// - nombre
// - len(apellido) in 4 bytes
// - apellido
// - len(dni) in 4 bytes
// - dni
// - len(nacimiento) in 4 bytes
// - nacimiento
// - len(numero) in 4 bytes
// - numero
func SendSerializedBet(conn net.Conn, bet SerializedBet) error {
	for _, data := range bet.data {
		err := SendLen(conn, uint32(len(data)))
		if err != nil {
			return err
		}

		err = SendData(conn, data, uint32(len(data)))
		if err != nil {
			return err
		}
	}
	return nil
}

// SendLen sends a uint32 through the connection
func SendLen(conn net.Conn, data uint32) error {
	err := binary.Write(conn, binary.BigEndian, data)
	if err != nil {
		log.Criticalf("action: SendLen | result: fail | error: %v", err)
		return err
	}
	log.Infof("action: SendLen | result: success | Len: %v", data)
	return nil
}

// SendData sends bytesToSend bytes of data through the connection
func SendData(conn net.Conn, data []byte, bytesToSend uint32) error {
	var total uint32 = 0
	for total < bytesToSend {
		n, err := conn.Write(data[total:])
		if err != nil {
			log.Criticalf("action: SendData | result: fail | error: Only %v of %v bytes were sent", total, bytesToSend)
			return err
		}
		total += uint32(n)
	}
	log.Infof("action: SendData | result: success | Bytes sent: %v | Data: %v", total, data)
	return nil
}

