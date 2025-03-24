package common

import (
	"bufio"
)

const CONFIRMATION_MESSAGE_SIZE = 1
const READING_ERROR = 1

// SendBatchSize sends the batch size to the server
func (c *Client) SendBatchSize(batchSize int) error {
	log.Debugf("action: sending_batch_size | result: in_progress | client_id: %v", c.config.ID)

	err := SenderSendBatchSize(c.conn, uint32(batchSize))
	if err != nil {
		log.Criticalf("action: sending_batch_size | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return err
	}

	log.Debugf("action: sending_batch_size | result: success | client_id: %v", c.config.ID)
	return nil
}

// SendMessage serializes and sends a message to the server
func (c *Client) SendBet(msgID int, betData Bet) error {
	log.Debugf("action: sending_bet | result: in_progress | client_id: %v | msg_id: %v", c.config.ID, msgID)

	serializedBet := SerializeBet(betData)
	err := SendSerializedBet(c.conn, serializedBet)
	if err != nil {
		log.Criticalf("action: sending_bet | result: fail | client_id: %v | msg_id: %v | error: %v", c.config.ID, msgID, err)
		return err
	}

	log.Debugf("action: sending_bet | result: success | client_id: %v | msg_id: %v", c.config.ID, msgID)
	return nil
}

// ReceiveConfirmation receives a confirmation message from the server
func (c *Client) ReceiveConfirmation() (int, error) {
	log.Debugf("action: receive_confirmation | result: in_progress | client_id: %v", c.config.ID)

	reader := bufio.NewReader(c.conn)
	confirmationBuf := make([]byte, CONFIRMATION_MESSAGE_SIZE)

	_, err := reader.Read(confirmationBuf)
	if err != nil {
		log.Criticalf("action: receive_confirmation | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return READING_ERROR, err
	}

	confirmation := confirmationBuf[0]

	log.Debugf("action: receive_confirmation | result: success | client_id: %v | confirmation: %v", c.config.ID, confirmation)
	return int(confirmation), nil
}