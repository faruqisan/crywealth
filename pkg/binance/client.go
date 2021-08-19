package binance

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var (
	defaultTimeout = time.Second * 10
)

type (
	Client struct {
		host       string
		httpClient *http.Client
		apiKey     string
		apiSecret  string
	}
)

func NewClient(host string, apiKey, apiSecret string) *Client {
	client := &http.Client{
		Timeout: defaultTimeout,
	}
	return &Client{
		host:       host,
		httpClient: client,
		apiKey:     apiKey,
		apiSecret:  apiSecret,
	}
}

// Do function add required header
func (c *Client) doSigned(method, endpoint string, params map[string]string) (*http.Response, error) {

	baseURL := fmt.Sprintf("%s/%s", c.host, endpoint)

	req, err := http.NewRequest(http.MethodGet, baseURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-MBX-APIKEY", c.apiKey)

	q := req.URL.Query()
	for key, val := range params {
		q.Set(key, val)
	}

	ts := strconv.FormatInt(unixMillis(time.Now()), 10)
	sig := c.GenerateSignature([]byte(q.Encode() + "&timestamp=" + ts))
	req.URL.RawQuery = fmt.Sprintf("%s&timestamp=%s&signature=%s", q.Encode(), ts, sig)

	return c.httpClient.Do(req)
}

func (c *Client) do(method, endpoint string, params map[string]string) (*http.Response, error) {
	baseURL := fmt.Sprintf("%s/%s", c.host, endpoint)
	req, err := http.NewRequest(http.MethodGet, baseURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	q := req.URL.Query()
	for key, val := range params {
		q.Set(key, val)
	}
	req.URL.RawQuery = q.Encode()
	return c.httpClient.Do(req)
}

func unixMillis(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}
