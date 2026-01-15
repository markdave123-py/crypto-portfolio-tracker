package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

type Client struct {
	httpClient *http.Client
	limiter    *rate.Limiter
	apiKey     string
	baseURL    string
}

func NewClient(apiKey, baseURL string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		// rate limiting requests forwarded to coingecko
		limiter: rate.NewLimiter(rate.Every(2*time.Second), 1),
		apiKey:  apiKey,
		baseURL: baseURL,
	}
}

func (c *Client) FetchTokenPrices(
	ctx context.Context,
	chain string,
	contracts []string,
) (TokenPriceResponse, error) {

	if c.limiter != nil {
		if err := c.limiter.Wait(ctx); err != nil {
			return nil, err
		}
	}

	url := fmt.Sprintf(
		"%s/simple/token_price/%s?contract_addresses=%s&vs_currencies=usd",
		c.baseURL,
		chain,
		strings.Join(contracts, ","),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-cg-demo-api-key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("coingecko error %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var decoded TokenPriceResponse
	if err := json.Unmarshal(body, &decoded); err != nil {
		return nil, err
	}

	return decoded, nil
}
