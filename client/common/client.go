package common

import (
	"net"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	bet    BetConfig
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig, bet BetConfig) *Client {
	client := &Client{
		config: config,
		bet: bet,
	}
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Criticalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return err
	}
	c.conn = conn
	return nil
}


// makeBet Sends a bet to the server and receives confirmation
func (c *Client) makeBet(msgID int) {
	if err := c.SendBet(msgID, c.bet); err != nil {
		log.Criticalf("action: send_bet | result: fail | client_id: %v | msg_id: %v | error: %v", c.config.ID, msgID, err)
		c.conn.Close()
		return
	}

	confirmation, err := c.ReceiveConfirmation()
	if err != nil {
		log.Criticalf("action: receive_confirmation | result: fail | client_id: %v | msg_id: %v | error: %v", c.config.ID, msgID, err)
		c.conn.Close()
		return
	}

	if confirmation == BET_SUCCESS {
		log.Infof("action: apuesta_enviada | result: success | dni: %s | numero: %s", c.bet.DNI, c.bet.Numero)
	} else {
		log.Criticalf("action: apuesta_enviada | result: failed_confirmation | dni: %s | numero: %s", c.bet.DNI, c.bet.Numero)
	}
}


// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	// There is an autoincremental msgID to identify every message sent
	// Messages if the message amount threshold has not been surpassed
	for msgID := 1; msgID <= c.config.LoopAmount; msgID++ {
		// Create the connection the server in every loop iteration. Send an
		c.createClientSocket()

		c.makeBet(msgID)

		c.conn.Close()

		// Wait a time between sending one message and the next one
		time.Sleep(c.config.LoopPeriod)

	}
	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}

func (c *Client) StopClient() {
	if c.conn != nil {
		c.conn.Close()
		log.Infof("action: stop_client | result: success | client_id: %v", c.config.ID)
	}
}