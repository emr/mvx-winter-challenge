package transactions

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const transactionsQuery = "%s/transactions?from=%d&size=%d&sender=%s&withScResults=false"

type Fetcher struct {
	apiEndpoint string
	maxRetries  int
	retryDelay  time.Duration
	pageSize    int
}

func NewFetcher(apiEndpoint string, maxRetries int, retryDelay time.Duration, pageSize int) *Fetcher {
	return &Fetcher{
		apiEndpoint: apiEndpoint,
		maxRetries:  maxRetries,
		retryDelay:  retryDelay,
		pageSize:    pageSize,
	}
}

func (f *Fetcher) FetchTransactions(address string, from int) ([]Transaction, error) {
	url := fmt.Sprintf(transactionsQuery, f.apiEndpoint, from, f.pageSize, address)

	var response []Transaction
	for attempt := 0; attempt < f.maxRetries; attempt++ {
		resp, err := http.Get(url)
		if err != nil {
			time.Sleep(f.retryDelay)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			time.Sleep(f.retryDelay)
			continue
		}

		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, fmt.Errorf("failed to decode response: %w", err)
		}
		return response, nil
	}
	return nil, fmt.Errorf("max retries reached")
}
