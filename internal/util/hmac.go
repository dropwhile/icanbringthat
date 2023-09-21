package util

import (
	"crypto/hmac"
	"crypto/sha256"
)

type Hmac struct {
	key []byte
}

func (h *Hmac) Validate(message, messageMAC []byte) bool {
	mac := hmac.New(sha256.New, h.key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

func (h *Hmac) Generate(message []byte) []byte {
	mac := hmac.New(sha256.New, h.key)
	mac.Write(message)
	return mac.Sum(nil)
}

func NewHmac(key []byte) *Hmac {
	return &Hmac{key}
}
