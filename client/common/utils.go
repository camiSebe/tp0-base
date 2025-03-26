package common

import (
	"encoding/csv"
	"os"
)

const BET_RECORD_LENGTH = 5
const BATCH_SUCCESS = 0

func ProcessFileOfBets(protocol *ProtocolClient, agency string, filePath string, batchMaxAmount int) {
	var err error
	var bets []Bet
	bets, err = LoadBetsFromCSV(agency, filePath)
	if err != nil {
		log.Criticalf("action: load_bets | result: fail | error: %v", err)
		return
	}

	ProcessBets(bets, protocol, batchMaxAmount)
}

// LoadBetsFromCSV loads bets from a CSV file and returns a slice of Bet
func LoadBetsFromCSV(agency string, filePath string) ([]Bet, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Criticalf("action: Opening_Bets_from_csv | result: fail | error: %v", err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	var bets []Bet

	records, err := reader.ReadAll()
	if err != nil {
		log.Criticalf("action: Loading_Bets_from_csv | result: fail | error: %v", err)
		return nil, err
	}

	for _, record := range records {
		if len(record) < BET_RECORD_LENGTH {
			log.Criticalf("action: Loading_Bets_from_csv | result: fail | error: invalid record")
			continue
		}

		bet := Bet{
			Nombre:     record[0],
			Apellido:   record[1],
			DNI:        record[2],
			Nacimiento: record[3],
			Numero:     record[4],
			Agencia:    agency,
		}
		bets = append(bets, bet)
	}

	return bets, nil
}

// ProcessBets separates bets into batches and sends them to the server
func ProcessBets(bets []Bet, protocol *ProtocolClient, batchMaxAmount int) {
	batches := CalulateBatchAmount(len(bets), batchMaxAmount)
	log.Infof("action: sending_all_batches | result: in_progress | batches: %v", batches)
	

	for i := 0; i < batches; i++ {
		log.Infof("action: sending_batch | result: in_progress | batch: %v / %v", i, batches)
		betsOfThisBatch := GetBetsOfThisBatch(bets, i, batchMaxAmount)

		err := protocol.SendBatch(betsOfThisBatch)
		if err != nil {
			log.Criticalf("action: send_batch | result: fail | error: %v", err)
			return
		}

		confirmation, err := protocol.ReceiveConfirmation()
		if err != nil {
			log.Criticalf("action: receive_confirmation | result: fail | error: %v", err)
			return
		}

		if confirmation == BATCH_SUCCESS {
			log.Debugf("action: batch_recibido | result: success")
		} else {
			log.Criticalf("action: batch_recibido | result: fail")
		}

		log.Infof("action: sending_batch | result: success | batch: %v / %v", i, batches)
	}
	
	log.Infof("action: sending_all_batches | result: success")

	// protocol.SendEndOfBatchesCode()
}

// CalulateBatchAmount calculates the number of batches needed to send all bets
func CalulateBatchAmount(betsAmount int, batchMaxAmount int) int {
	batches := betsAmount / batchMaxAmount
	if betsAmount%batchMaxAmount != 0 {
		batches++
	}
	return batches
}

// GetBetsOfThisBatch returns a slice of bets that belong to a batch
func GetBetsOfThisBatch(bets []Bet, i int, batchMaxAmount int) []Bet {
	start := i * batchMaxAmount
	end := (i + 1) * batchMaxAmount
	if end > len(bets) {
		end = len(bets)
	}
	return bets[start:end]
}
