// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package app

type Config struct {
	WebhookCreds   map[string]string
	BaseURL        string
	CSRFKeyBytes   []byte
	HMACKeyBytes   []byte
	Production     bool
	RequestLogging bool
	RpcApi         bool
}
