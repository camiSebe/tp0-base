package common

import (
	"bufio"
	"net"
)

const CONFIRMATION_MESSAGE_SIZE = 1
const READING_ERROR = 1

type ProtocolClient struct {
	conn net.Conn
}

// NewProtocolClient initializes a new ProtocolClient
func NewProtocolClient(conn net.Conn) *ProtocolClient {
	return &ProtocolClient{conn}
}

// SendMessage serializes and sends a message to the server
func (p *ProtocolClient) SendBet(client_id string, msgID int, betData BetConfig) error {
	log.Infof("action: sending_bet | result: in_progress | client_id: %v | msg_id: %v", client_id, msgID)

	serializedBet := SerializeBet(betData)
	err := SendSerializedBet(p.conn, serializedBet)
	if err != nil {
		log.Criticalf("action: sending_bet | result: fail | client_id: %v | msg_id: %v | error: %v", client_id, msgID, err)
		return err
	}

	log.Infof("action: sending_bet | result: success | client_id: %v | msg_id: %v", client_id, msgID)
	return nil
}

// ReceiveConfirmation receives a confirmation message from the server
func (p *ProtocolClient) ReceiveConfirmation(client_id string) (int, error) {
	log.Infof("action: receive_confirmation | result: in_progress | client_id: %v", client_id)

	reader := bufio.NewReader(p.conn)
	confirmationBuf := make([]byte, CONFIRMATION_MESSAGE_SIZE)

	_, err := reader.Read(confirmationBuf)
	if err != nil {
		log.Criticalf("action: receive_confirmation | result: fail | client_id: %v | error: %v", client_id, err)
		return READING_ERROR, err
	}

	confirmation := confirmationBuf[0]

	log.Infof("action: receive_confirmation | result: success | client_id: %v | confirmation: %v", client_id, confirmation)
	return int(confirmation), nil
}