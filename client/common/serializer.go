package common

import (
	"encoding/binary"
)

const SIZE_OF_UINT32 = 4
const MIN_VALUE_FOR_UINT8 = 0
const MAX_VALUE_FOR_UINT8 = 255

// SerializeUint32 serializes an int into a byte array
func serializeUint32(value int) []byte {
	data := make([]byte, SIZE_OF_UINT32)
	binary.BigEndian.PutUint32(data, uint32(value))
	return data
}

// SerializeUint8 serializes an int into a byte array
func serializeUint8(value int) []byte {
	if value < MIN_VALUE_FOR_UINT8 || value > MAX_VALUE_FOR_UINT8 {
		log.Criticalf("action: serialize_uint8 | result: fail | error: value %v is out of range", value)
		return nil
	}
	return []byte{uint8(value)}
}