package common

import (
	"net"
)

const BET_MESSAGE_CODE = 1
const START_OF_BATCH_MESSAGE_CODE = 2
const END_OF_BATCH_MESSAGE_CODE = 3
const GET_WINNERS_MESSAGE_CODE = 4

const CONFIRMATION_MESSAGE_SIZE = 1
const READING_ERROR = 1

type ProtocolClient struct {
	conn net.Conn
}

// NewProtocolClient initializes a new ProtocolClient
func NewProtocolClient(conn net.Conn) *ProtocolClient {
	return &ProtocolClient{conn}
}

// SendBatch sends a batch of bets to the server
func (p *ProtocolClient) SendBatch(batch []Bet) error {
	err := p.SendStartOfBatchCode()
	if err != nil {
		return err
	}

	err = p.SendBatchSize(len(batch))
	if err != nil {
		return err
	}

	for _, bet := range batch {
		err = p.SendBet(bet)
		if err != nil {
			return err
		}
	}

	return nil
}

// SendStartOfBatchCode sends the start of batch code to the server
func (p *ProtocolClient) SendStartOfBatchCode() error {
	log.Debugf("action: sending_start_of_batch_code | result: in_progress")
	code := serializeUint8(START_OF_BATCH_MESSAGE_CODE)
	if code == nil {
		log.Criticalf("action: sending_start_of_batch_code | result: fail | error: Not able to serialize as uint8")
		return nil
	}

	err := SendAll(p.conn, code)
	if err != nil {
		log.Criticalf("action: sending_start_of_batch_code | result: fail | error: %v", err)
		return err
	}
	log.Debugf("action: sending_start_of_batch_code | result: success")
	return nil
}


// SendBatchSize sends the batch size to the server
func (p *ProtocolClient) SendBatchSize(batchSize int) error {
	log.Debugf("action: sending_batch_size | result: in_progress")

	batchSizeBytes := serializeUint32(batchSize)

	err := SendAll(p.conn, batchSizeBytes)
	if err != nil {
		log.Criticalf("action: sending_batch_size | result: fail | error: %v", err)
		return err
	}

	log.Debugf("action: sending_batch_size | result: success")
	return nil
}

// SendBet sends a bet to the server
func (p *ProtocolClient) SendBet(bet Bet) error {
	log.Debugf("action: sending_bet | result: in_progress")

	// Because i'm sending a batch of bets, server already knows that the message is a bet -> no need to send the message code
	/*
	code := serializeUint8(BET_MESSAGE_CODE)
	if code == nil {
		log.Criticalf("action: sending_bet | result: fail | error: Not able to serialize as uint8")
		return nil
	}
	err := SendAll(p.conn, code)
	if err != nil {
		return err
	}
	*/

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

	log.Debugf("action: sending_bet | result: success")
	return nil
}

// ReceiveConfirmation receives a confirmation message from the server
func (p *ProtocolClient) ReceiveConfirmation() (int, error) {
	log.Debugf("action: receive_confirmation | result: in_progress")

	data := make([]byte, CONFIRMATION_MESSAGE_SIZE)
	err := ReceiveAll(p.conn, data)
	if err != nil {
		log.Criticalf("action: receive_confirmation | result: fail | error: %v", err)
		return READING_ERROR, err
	}

	confirmation := int(data[0])
	log.Debugf("action: receive_confirmation | result: success | confirmation: %v", confirmation)
	return confirmation, nil
}