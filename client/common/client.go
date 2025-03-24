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
	bets    []Bet
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

// makeBet Sends a bet to the server and receives confirmation
func (c *Client) makeBet(msgID int, bet Bet) {
    if err := c.SendBet(msgID, bet); err != nil {
		log.Criticalf("action: send_bet | result: fail | client_id: %v | msg_id: %v | error: %v", c.config.ID, msgID, err)
		c.conn.Close()
		return
	}
}

func setBatchOfBets(bets []Bet, start int, end int) []Bet {
	if end > len(bets) {
		end = len(bets)
	}
	log.Debugf("action: set_batch_of_bets | result: success | start: %v | end: %v", start, end)
	return bets[start:end]
}


func (c *Client) SendBatchOfBets(batchOfBets []Bet) {
	for i := 0; i < len(batchOfBets); i++ {
		c.makeBet(i, batchOfBets[i])
	}
}

func (c *Client) ProcessBatch(i int) {	
	batch := setBatchOfBets(c.bets, i*c.config.BatchMaxAmount, (i+1)*c.config.BatchMaxAmount)

	c.SendBatchSize(len(batch))
	
	c.SendBatchOfBets(batch)
}

func (c *Client) ProcessBets() {
	batches := len(c.bets) / c.config.BatchMaxAmount
	if len(c.bets) % c.config.BatchMaxAmount != 0 {
		batches++
	}

	for i := 0; i < batches; i++ {
		log.Infof("action: sending_batch | result: in_progress | client_id: %v | batch: %v / %v", c.config.ID, i, batches)
		c.createClientSocket()
		c.ProcessBatch(i)
		confirmation, err := c.ReceiveConfirmation()
		if err != nil {
			log.Criticalf("action: receive_confirmation | result: fail | client_id: %v | error: %v", c.config.ID, err)
			return
		}

		if confirmation == BET_SUCCESS {
			log.Debugf("action: batch_recibido | result: success | client_id: %v", c.config.ID)
		} else {
			log.Criticalf("action: batch_recibido | result: fail")
		}
		c.conn.Close()
		// time.Sleep(c.config.LoopPeriod)
	}
	log.Infof("action: sending_batch | result: success | client_id: %v", c.config.ID)
}


// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	var err error
	c.bets, err = LoadBetsFromCSV(c.config.ID, fmt.Sprintf(".data/dataset/agency-%s.csv", c.config.ID))
	if err != nil {
		log.Criticalf("action: load_bets | result: fail | client_id: %v | error: %v", c.config.ID, err)
		return
	}
	c.ProcessBets()
}

func (c *Client) StopClient() {
	if c.conn != nil {
		c.conn.Close()
		log.Infof("action: stop_client | result: success | client_id: %v", c.config.ID)
	}
}