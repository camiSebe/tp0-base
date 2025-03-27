package common

import (
	"encoding/csv"
	"os"
	"strconv"
)

const MAX_BATCH_AMOUNT_FOR_PROTOCOL = 100

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
		log.Debugf("action: sending_batch | result: in_progress | batch: %v / %v", i, batches)
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

		log.Debugf("action: sending_batch | result: success | batch: %v / %v", i, batches)
	}
	
	log.Infof("action: sending_all_batches | result: success")
}

// CalulateBatchAmount calculates the number of batches needed to send all bets
func CalulateBatchAmount(betsAmount int, batchMaxAmount int) int {
	if batchMaxAmount > MAX_BATCH_AMOUNT_FOR_PROTOCOL {
		batchMaxAmount = MAX_BATCH_AMOUNT_FOR_PROTOCOL
	}
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

// GetWinners asks the server for the winners
func GetWinners(protocol *ProtocolClient, agency string) {
	log.Infof("action: consulta_ganadores | result: in_progress")

	err := protocol.SendGetWinnersCode()
	if err != nil {
		log.Criticalf("action: consulta_ganadores | result: fail | error: %v", err)
		return
	}

	agencyNumber, err := strconv.Atoi(agency)
	if err != nil {
		log.Criticalf("action: consulta_ganadores | result: fail | error: invalid agency number | error: %v", err)
		return
	}

	err = protocol.SendAgencyNumber(agencyNumber)
	if err != nil {
		log.Criticalf("action: consulta_ganadores | result: fail | error: %v", err)
		return
	}

	amountOfWinners, err := protocol.ReceiveListOfWinnersSize()
	if err != nil {
		log.Criticalf("action: consulta_ganadores | result: fail | error: %v", err)
		return
	}

	for i := 0; i < amountOfWinners; i++ {
		document, err := protocol.ReceiveWinnerDocument()
		if err != nil {
			log.Criticalf("action: consulta_ganadores | result: fail | error: %v", err)
			return
		}
		log.Infof("action: ganador_obtenido | result: success | DNI_ganador: %v", document)
	}

	log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %v", amountOfWinners)
}