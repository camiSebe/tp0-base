package common

import (
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
func (p *ProtocolClient) SendBet(client_id string, msgID int, bet BetConfig) error {
	log.Infof("action: sending_bet | result: in_progress | client_id: %v | msg_id: %v", client_id, msgID)

	fields := []string{bet.Nombre, bet.Apellido, bet.DNI, bet.Nacimiento, bet.Numero, bet.Agencia}

	for _, field := range fields {
		lengthData := serializeUint8(len(field))
		if lengthData == nil {
			log.Criticalf("action: sending_bet | result: fail | error: Not able to serialize as uint8")
			return nil
		}

		if err := SendAll(p.conn, lengthData); err != nil {
			return err
		}
		if err := SendAll(p.conn, []byte(field)); err != nil {
			return err
		}
	}

	log.Infof("action: sending_bet | result: success | client_id: %v | msg_id: %v", client_id, msgID)
	return nil
}

// ReceiveConfirmation receives a confirmation message from the server
func (p *ProtocolClient) ReceiveConfirmation(client_id string) (int, error) {
	log.Infof("action: receive_confirmation | result: in_progress | client_id: %v", client_id)

	data := make([]byte, CONFIRMATION_MESSAGE_SIZE)
	err := ReceiveAll(p.conn, data)
	if err != nil {
		log.Criticalf("action: receive_confirmation | result: fail | client_id: %v | error: %v", client_id, err)
		return READING_ERROR, err
	}

	confirmation := int(data[0])

	log.Infof("action: receive_confirmation | result: success | client_id: %v | confirmation: %v", client_id, confirmation)
	return int(confirmation), nil
}