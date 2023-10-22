package util

import (
	"crypto/hmac"

	"github.com/zeebo/blake3"
)

// keyed MAC using blake3
// similar to HMAC, but simpler and faster (while offering similar security)
type MAC struct {
	key []byte
}

func (h *MAC) Validate(message, messageMAC []byte) bool {
	// only panics on invalid keysize, which shouldn't happen
	mac, _ := blake3.NewKeyed(h.key)
	mac.Write(message) // #nosec G104 -- doesn't actually return errors
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

func (h *MAC) Generate(message []byte) []byte {
	// only panics on invalid keysize, which shouldn't happen
	mac, _ := blake3.NewKeyed(h.key)
	mac.Write(message) // #nosec G104 -- doesn't actually return errors
	return mac.Sum(nil)
}

func NewMAC(key []byte) *MAC {
	derivedKey := make([]byte, 32)
	blake3.DeriveKey(
		"icanbringthat 2023-10-22T04:03:24.602Z keyed mac", // context
		key,        // material
		derivedKey, // output
	)
	return &MAC{derivedKey}
}
