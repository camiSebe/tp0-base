package common

import (
	"bufio"
	"fmt"
)

// SendMessage receives a message and sends it to the server
func (c *Client) SendMessage(msgID int) error {
	message := fmt.Sprintf("[CLIENT %v] Message N°%v\n", c.config.ID, msgID)
	_, err := fmt.Fprintf(c.conn, "%s", message)
	if err != nil {
		log.Criticalf("action: send_message | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return err
	}
	log.Infof("action: send_message | result: success | client_id: %v | msg_id: %v", c.config.ID, msgID)
	return nil
}

// ReceiveMessage receives a message from the server
func (c *Client) ReceiveMessage() (string, error) {
	msg, err := bufio.NewReader(c.conn).ReadString('\n')
	if err != nil {
		log.Criticalf("action: receive_message | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return "", err
	}
	log.Infof("action: receive_message | result: success | client_id: %v | msg: %v", c.config.ID, msg)
	return msg, nil
}