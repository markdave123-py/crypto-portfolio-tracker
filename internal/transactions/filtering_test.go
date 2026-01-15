package transactions

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestApplyFilters_Match(t *testing.T) {
	now := time.Now()

	tx := Transaction{
		Type:      TypeSwap,
		Status:    StatusSuccess,
		Token:     "ETH",
		Timestamp: now,
	}

	f := Filters{
		Type:   ptrType(TypeSwap),
		Status: ptrStatus(StatusSuccess),
	}

	ok := applyFilters(tx, f)

	require.True(t, ok)
}

func TestApplyFilters_TypeMismatch(t *testing.T) {
	tx := Transaction{
		Type: TypeSend,
	}

	f := Filters{
		Type: ptrType(TypeSwap),
	}

	ok := applyFilters(tx, f)

	require.False(t, ok)
}

func TestApplyFilters_DateRange(t *testing.T) {
	now := time.Now()
	past := now.Add(-time.Hour)
	future := now.Add(time.Hour)

	tx := Transaction{
		Timestamp: now,
	}

	f := Filters{
		StartDate: &past,
		EndDate:   &future,
	}

	ok := applyFilters(tx, f)

	require.True(t, ok)
}

func ptrType(v TransactionType) *TransactionType {
	return &v
}

func ptrStatus(v TransactionStatus) *TransactionStatus {
	return &v
}
