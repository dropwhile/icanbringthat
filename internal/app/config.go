package app

type Config struct {
	HMACKey      string
	WebhookCreds map[string]string
	CSRFKeyBytes []byte
	HMACKeyBytes []byte
	Production   bool
	BaseURL      string
}
