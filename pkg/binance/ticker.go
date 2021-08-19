package binance

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	tickerURL = "api/v3/ticker/price"
)

type (
	TickerResponse struct {
		Symbol string `json:"symbol"`
		Price  string `json:"price"`
	}
)

func (c *Client) GetSymbolTicker(symbol string) (resp TickerResponse, err error) {
	res, err := c.do(http.MethodGet, tickerURL, map[string]string{
		"symbol": symbol,
	})
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return resp, err
	}

	if err = json.Unmarshal(body, &resp); err != nil {
		return resp, err
	}
	return
}

func (c *Client) GetAllTicker() (resp []TickerResponse, err error) {
	res, err := c.do(http.MethodGet, tickerURL, nil)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return resp, err
	}

	if err = json.Unmarshal(body, &resp); err != nil {
		return resp, err
	}
	return
}
