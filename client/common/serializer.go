package common

import (
	"encoding/binary"
	"bytes"
	"fmt"
)

// SerializeMessage serializes the message by prefixing it with its size
func SerializeMessage(msg string) ([]byte, error) {
	msgBytes := []byte(msg)
	msgSize := uint32(len(msgBytes))
	log.Infof("action: serialize_message | result: success | size: %v | msg: %v", msgSize, msg)

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, msgSize); err != nil {
		log.Criticalf("action: serialize_message | result: fail | error: %v", err)
		return nil, err
	}
	buf.Write(msgBytes)

	return buf.Bytes(), nil
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
