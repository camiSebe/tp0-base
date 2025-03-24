package common

import (
	"encoding/csv"
	"os"
)

// LoadBetsFromCSV loads bets from a CSV file
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
