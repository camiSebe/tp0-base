package common

type SerializedBet struct {
	data [][]byte
}

// SerializeBet serializes a bet into a SerializedBet struct that has all their fields as byte slices
func SerializeBet(betData BetConfig) SerializedBet {
	return SerializedBet{
		data: [][]byte{
			[]byte(betData.Nombre),
			[]byte(betData.Apellido),
			[]byte(betData.DNI),
			[]byte(betData.Nacimiento),
			[]byte(betData.Numero),
			[]byte(betData.Agencia),
		},
	}
}