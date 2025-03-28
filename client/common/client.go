package common

import (
	"net"
	"time"
	"fmt"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("log")

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopAmount    int
	LoopPeriod    time.Duration
	BatchMaxAmount int
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
}

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	client := &Client{
		config: config,
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

// StartClientLoop Sends batches of bets to the server and asks for the winners
func (c *Client) StartClientLoop() {
	c.createClientSocket()

	protocol := NewProtocolClient(c.conn)

	err := ProcessFileOfBets(protocol, c.config.ID, fmt.Sprintf(".data/agency-%s.csv", c.config.ID), c.config.BatchMaxAmount)
	if err != nil {
		log.Criticalf("action: process_file_of_bets | result: fail | error: %v", err)
	}

	GetWinners(protocol, c.config.ID)

	// Sleep for the tests
	time.Sleep(50 * time.Second)

	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) StopClient() {
	if c.conn != nil {
		c.conn.Close()
		log.Infof("action: stop_client | result: success | client_id: %v", c.config.ID)
	}
}