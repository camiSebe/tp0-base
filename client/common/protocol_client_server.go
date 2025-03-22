package common

import (
	"bufio"
	"fmt"

	"encoding/binary"
)

// SendMessage serializes and sends a message to the server
func (c *Client) SendMessage(msgID int) error {
	message := fmt.Sprintf("[CLIENT %v] Message N°%v", c.config.ID, msgID)
	serializedMsg, err := SerializeMessage(message)
	if err != nil {
		log.Criticalf("action: serialize_message | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return err
	} else {
		log.Infof("action: serialize_message | result: success | client_id: %v | msg: %v", c.config.ID, message)
	}

	_, err = c.conn.Write(serializedMsg)
	if err != nil {
		log.Criticalf("action: send_message | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return err
	}

	log.Infof("action: send_message | result: success | client_id: %v | msg_id: %v", c.config.ID, msgID)
	return nil
}

// ReceiveMessage reads and deserializes a message from the server
func (c *Client) ReceiveMessage() (string, error) {
	reader := bufio.NewReader(c.conn)
	sizeBuf := make([]byte, SIZE_UINT32)
	_, err := reader.Read(sizeBuf)
	if err != nil {
		log.Criticalf("action: receive_message_size | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return "", err
	}

	size := binary.BigEndian.Uint32(sizeBuf)
	msgBuf := make([]byte, size)
	_, err = reader.Read(msgBuf)
	if err != nil {
		log.Criticalf("action: receive_message | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return "", err
	}

	message, err := DeserializeMessage(append(sizeBuf, msgBuf...))
	if err != nil {
		log.Criticalf("action: deserialize_message | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return "", err
	}

	log.Infof("action: receive_message | result: success | client_id: %v | msg: %v", c.config.ID, message)
	return message, nil
}