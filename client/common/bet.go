package common

type Bet struct {
	Nombre        string
	Apellido      string
	DNI           string
	Nacimiento    string
	Numero        string
	Agencia		  string
}

const BET_SUCCESS = 0
const BET_FAIL = 1

const BET_RECORD_LENGTH = 5