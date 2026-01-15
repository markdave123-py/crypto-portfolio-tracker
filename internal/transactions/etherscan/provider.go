package etherscan

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/markdave123-py/crypto-portfolio-tracker/internal/transactions"
	"github.com/markdave123-py/crypto-portfolio-tracker/internal/utils"
)

type Provider struct {
	client ClientAPI
}

func NewProvider(client ClientAPI) *Provider {
	return &Provider{client: client}
}

func (p *Provider) GetTransactions(
	ctx context.Context,
	chain string,
	wallet string,
	page int,
	limit int,
) ([]transactions.Transaction, error) {

	var resp *txListResponse

	err := utils.Retry(ctx, utils.RetryConfig{
		MaxRetries: 3,
		BaseDelay:  500 * time.Millisecond,
		MaxDelay:   4 * time.Second,
	}, func() error {

		r, err := p.client.FetchTxList(ctx, chain, wallet, page, limit)
		if err != nil {
			return err
		}

		resp = r
		return nil
	})

	if err != nil {
		return nil, err
	}

	txs := make([]transactions.Transaction, 0, len(resp.Result))

	for _, item := range resp.Result {
		txType := classifyType(item)
		status := classifyStatus(item)

		ts, _ := strconv.ParseInt(item.TimeStamp, 10, 64)

		amount := weiToEther(item.Value)

		tx := transactions.Transaction{
			ID:        item.Hash,
			Chain:     chain,
			Hash:      item.Hash,
			From:      strings.ToLower(item.From),
			To:        strings.ToLower(item.To),
			Amount:    amount,
			Type:      txType,
			Status:    status,
			Timestamp: time.Unix(ts, 0),
		}

		txs = append(txs, tx)
	}

	return txs, nil
}
