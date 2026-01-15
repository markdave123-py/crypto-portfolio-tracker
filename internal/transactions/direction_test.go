package transactions

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDetectDirection_Outgoing(t *testing.T) {
	tx := Transaction{
		From: "0xabc",
		To:   "0xdef",
	}

	out := detectDirection(tx, "0xabc")

	require.Equal(t, DirectionOut, out.Direction)
}

func TestDetectDirection_Incoming(t *testing.T) {
	tx := Transaction{
		From: "0xabc",
		To:   "0xdef",
	}

	out := detectDirection(tx, "0xdef")

	require.Equal(t, DirectionIn, out.Direction)
}
