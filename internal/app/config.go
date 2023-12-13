package app

type Config struct {
	WebhookCreds   map[string]string
	BaseURL        string
	CSRFKeyBytes   []byte
	HMACKeyBytes   []byte
	Production     bool
	RequestLogging bool
}
