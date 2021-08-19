package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func (c *Client) GenerateSignature(payload []byte) string {
	mac := hmac.New(sha256.New, []byte(c.apiSecret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}
