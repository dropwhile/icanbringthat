// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

//go:build generate

//go:generate protoc --proto_path=./protos --go_out=paths=source_relative:./icbt --twirp_out=paths=source_relative:./icbt ./protos/service.proto
//go:generate protoc-go-inject-tag -input=./icbt/service.pb.go

package rpc
