package icbt

//go:generate protoc --go_out=paths=source_relative:. --twirp_out=paths=source_relative:. service.proto
//go:generate protoc-go-inject-tag -input=service.pb.go
