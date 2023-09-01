// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package icbt

//go:generate protoc --go_out=paths=source_relative:. --twirp_out=paths=source_relative:. service.proto
//go:generate protoc-go-inject-tag -input=service.pb.go
