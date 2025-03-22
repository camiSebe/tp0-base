package common

import (
	"bufio"
	"fmt"

	"encoding/binary"
)

// SendMessage serializes and sends a message to the server
func (c *Client) SendMessage(msgID int, betData BetData) error {
	log.Infof("action: sending_message | result: in progress | client_id: %v | msg_id: %v", c.config.ID, msgID)
	message := fmt.Sprintf("%v %v %v %v %v", betData.Nombre, betData.Apellido, betData.DNI, betData.Nacimiento, betData.Numero)
	serializedMsg, err := SerializeMessage(message)
	_, err = c.conn.Write(serializedMsg)
	if err != nil {
		log.Criticalf("action: sending_message | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return err
	}

	log.Infof("action: sending_message | result: success | client_id: %v | msg_id: %v", c.config.ID, msgID)
	return nil
}

// ReceiveMessage reads and deserializes a message from the server
func (c *Client) ReceiveMessage() (string, error) {
	log.Infof("action: receive_message_size | result: in progress | client_id: %v", c.config.ID)
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