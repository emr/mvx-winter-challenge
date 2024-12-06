package api

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
)

type Client struct {
	baseURL string
}

type TokenData struct {
	Identifier string `json:"identifier"`
	Name       string `json:"name"`
	Balance    *big.Int
	Owner      string `json:"owner"`
	Decimals   int    `json:"decimals"`
}

func NewClient(baseURL string) *Client {
	return &Client{baseURL: baseURL}
}

func (c *Client) GetAccountTokens(address, tokenName string) ([]TokenData, error) {
	url := fmt.Sprintf("%s/accounts/%s/tokens?type=FungibleESDT&name=%s",
		c.baseURL, address, tokenName)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tokens: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var tokens []TokenData
	if err := json.Unmarshal(body, &tokens); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return tokens, nil
}

// UnmarshalJSON implements json.Unmarshaler
func (t *TokenData) UnmarshalJSON(data []byte) error {
	type Alias TokenData
	aux := struct {
		Balance string `json:"balance"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Balance != "" {
		balance := new(big.Int)
		balance, ok := balance.SetString(aux.Balance, 10)
		if !ok {
			return fmt.Errorf("failed to parse balance: %s", aux.Balance)
		}
		t.Balance = balance
	}

	return nil
}

// MarshalJSON implements json.Marshaler
func (t TokenData) MarshalJSON() ([]byte, error) {
	type Alias TokenData
	return json.Marshal(&struct {
		Balance string `json:"balance"`
		*Alias
	}{
		Balance: t.Balance.String(),
		Alias:   (*Alias)(&t),
	})
}
