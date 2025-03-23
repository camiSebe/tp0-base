package common

import (
	"encoding/binary"
	"fmt"
)

type SerializedBet struct {
	data [][]byte
}

// SerializeBet serializes a bet into a SerializedBet struct
func SerializeBet(betData BetConfig) SerializedBet {
	return SerializedBet{
		data: [][]byte{
			[]byte(betData.Nombre),
			[]byte(betData.Apellido),
			[]byte(betData.DNI),
			[]byte(betData.Nacimiento),
			[]byte(betData.Numero),
		},
	}
}

// DeserializeMessage extracts the message size and content from a byte slice
func DeserializeMessage(data []byte) (string, error) {
	if len(data) < SIZE_UINT32 {
		err := fmt.Errorf("invalid data length")
		log.Criticalf("action: deserialize_message | result: fail | error: %v", err)
		return "", err
	}

	msgSize := binary.BigEndian.Uint32(data[:SIZE_UINT32])
	if int(msgSize) != len(data[SIZE_UINT32:]) {
		err := fmt.Errorf("message size mismatch")
		log.Criticalf("action: deserialize_message | result: fail | error: %v", err)
		return "", err
	}

	log.Infof("action: deserialize_message | result: success | size: %v | msg: %v", msgSize, string(data[SIZE_UINT32:]))
	return string(data[SIZE_UINT32:]), nil
}
