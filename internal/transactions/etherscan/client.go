package etherscan

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type ClientAPI interface {
	FetchTxList(
		ctx context.Context,
		chainID string,
		wallet string,
		page int,
		offset int,
	) (*txListResponse, error)
}

type Client struct {
	baseURL string
	apiKey  string
	client  *http.Client
	limiter *rate.Limiter
}

func NewClient(apiKey string, baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  &http.Client{},
		// rate limiting requests forwarded to etherscan
		limiter: rate.NewLimiter(rate.Every(2*time.Second), 1),
	}
}

func (c *Client) FetchTxList(
	ctx context.Context,
	chainID string,
	wallet string,
	page int,
	offset int,
) (*txListResponse, error) {

	if c.limiter != nil {
		if err := c.limiter.Wait(ctx); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf(
			"%s?module=account&chainid=%s&action=txlist&address=%s&startblock=0&endblock=99999999&page=%d&offset=%d&sort=desc&apikey=%s",
			c.baseURL,
			chainID,
			wallet,
			page,
			offset,
			c.apiKey,
		),
		nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var decoded txListResponse
	if err := json.Unmarshal(body, &decoded); err != nil {
		return nil, err
	}

	if decoded.Status != "1" {
		return nil, fmt.Errorf("etherscan error: %s", decoded.Message)
	}

	return &decoded, nil
}
