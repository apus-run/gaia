package helloworld

//go:generate protoc -I . -I ../../../third_party --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. ./helloworld.proto
//go:generate protoc -I . -I ../../../third_party --go_out=paths=source_relative:. --go-gin_out=paths=source_relative:. ./helloworld.proto
