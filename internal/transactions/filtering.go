package transactions

import "time"

type Filters struct {
	Type   *TransactionType
	Status *TransactionStatus
	Token  *string

	StartDate *time.Time
	EndDate   *time.Time
}

// This filters the transaction by Type, Status, Time
func applyFilters(tx Transaction, f Filters) bool {
	if f.Type != nil && tx.Type != *f.Type {
		return false
	}
	if f.Status != nil && tx.Status != *f.Status {
		return false
	}
	if f.Token != nil && tx.Token != *f.Token {
		return false
	}
	if f.StartDate != nil && tx.Timestamp.Before(*f.StartDate) {
		return false
	}
	if f.EndDate != nil && tx.Timestamp.After(*f.EndDate) {
		return false
	}
	return true
}
