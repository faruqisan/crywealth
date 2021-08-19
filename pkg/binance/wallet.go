package binance

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	accountSnapshotURL = "sapi/v1/accountSnapshot"
)

type (
	AccountSnapshotResponse struct {
		Code        int                 `json:"code"`
		Msg         string              `json:"msg"`
		SnapshotVos []AccountSnapshotVo `json:"snapshotVos"`
	}

	AccountSnapshotVo struct {
		Type       string `json:"type"`
		UpdateTime int64  `json:"updateTime"`
		Data       struct {
			TotalAssetOfBtc string `json:"totalAssetOfBtc"`
			Balances        []struct {
				Asset  string `json:"asset"`
				Free   string `json:"free"`
				Locked string `json:"locked"`
			} `json:"balances"`
		} `json:"data"`
	}
)

func (c *Client) AccountSnapshot(snapshotType string) (AccountSnapshotResponse, error) {
	var (
		resp AccountSnapshotResponse
		err  error
	)

	hresp, err := c.doSigned(http.MethodGet, accountSnapshotURL, map[string]string{
		"type": snapshotType,
	})
	if err != nil {
		return resp, err
	}
	defer hresp.Body.Close()

	body, err := ioutil.ReadAll(hresp.Body)
	if err != nil {
		return resp, err
	}

	if err = json.Unmarshal(body, &resp); err != nil {
		return resp, err
	}

	return resp, err
}
