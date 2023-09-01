// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package mail

type Config struct {
	Hostname    string
	Host        string
	User        string
	Pass        string
	DefaultFrom string
	Port        int
}
