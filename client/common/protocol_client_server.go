package common

import (
	"bufio"
	"encoding/binary"
)

// SendMessage serializes and sends a message to the server
func (c *Client) SendBet(msgID int, betData BetConfig) error {
	log.Infof("action: sending_bet | result: in progress | client_id: %v | msg_id: %v", c.config.ID, msgID)

	serializedBet := SerializeBet(betData)
	err := SendSerializedBet(c.conn, serializedBet)
	if err != nil {
		log.Criticalf("action: sending_bet | result: fail | client_id: %v | msg_id: %v | error: %v", c.config.ID, msgID, err)
		return err
	}

	log.Infof("action: sending_bet | result: success | client_id: %v | msg_id: %v", c.config.ID, msgID)
	return nil
}

// ReceiveMessage reads and deserializes a message from the server
func (c *Client) ReceiveConfirmation() (string, error) {
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