package transactions

type Transaction struct {
	TxHash        string `json:"txHash"`
	Nonce         uint64 `json:"nonce"`
	Timestamp     int64  `json:"timestamp"`
	Sender        string `json:"sender"`
	Receiver      string `json:"receiver"`
	Function      string `json:"function"`
	Value         string `json:"value"`
	MiniBlockHash string `json:"miniBlockHash"`
	Status        string `json:"status"`
}
