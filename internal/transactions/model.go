package transactions

import "time"

type Transaction struct {
	ID    string // internal ID (can be tx hash)
	Chain string // ethereum, polygon, etc
	Hash  string // on-chain tx hash

	From string
	To   string

	Token     string // ETH, USDC, WBTC, etc
	TokenAddr string // empty for native ETH

	Amount float64

	Type   TransactionType
	Status TransactionStatus

	Direction Direction

	Timestamp time.Time
}
