package common

import (
	"net"
	"encoding/binary"
)

// SendSerializedBet sends:
// - len(nombre) in 1 byte
// - nombre
// - len(apellido) in 1 byte
// - apellido
// - len(dni) in 1 byte
// - dni
// - len(nacimiento) in 1 byte
// - nacimiento
// - len(numero) in 1 byte
// - numero
// - len(agencia) in 1 byte
// - agencia
func SendSerializedBet(conn net.Conn, bet SerializedBet) error {
	for _, data := range bet.data {
		err := SendLen(conn, uint8(len(data)))
		if err != nil {
			return err
		}

		err = SendData(conn, data, uint8(len(data)))
		if err != nil {
			return err
		}
	}
	return nil
}

// SendLen sends a uint32 through the connection
func SendLen(conn net.Conn, data uint8) error {
	err := binary.Write(conn, binary.BigEndian, data)
	if err != nil {
		log.Criticalf("action: SendLen | result: fail | error: %v", err)
		return err
	}
	log.Debugf("action: SendLen | result: success | Len: %v", data)
	return nil
}

// SendData sends bytesToSend bytes of data through the connection
func SendData(conn net.Conn, data []byte, bytesToSend uint8) error {
	var total uint8 = 0
	for total < bytesToSend {
		n, err := conn.Write(data[total:])
		if err != nil {
			log.Criticalf("action: SendData | result: fail | error: Only %v of %v bytes were sent", total, bytesToSend)
			return err
		}
		total += uint8(n)
	}
	log.Debugf("action: SendData | result: success | Bytes sent: %v | Data: %v", total, data)
	return nil
}

