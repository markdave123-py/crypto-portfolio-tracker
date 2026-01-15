package transactions

type TransactionType string

const (
	TypeSend    TransactionType = "send"
	TypeReceive TransactionType = "receive"
	TypeSwap    TransactionType = "swap"
	TypeStake   TransactionType = "stake"
)

type TransactionStatus string

const (
	StatusSuccess TransactionStatus = "success"
	StatusFailed  TransactionStatus = "failed"
)

type Direction string

const (
	DirectionIn  Direction = "in"
	DirectionOut Direction = "out"
)
