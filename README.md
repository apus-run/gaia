# Gaia
Gaia[ˈɡaɪə] 一个轻量级gRPC业务框架

##  安装 protoc
```
$ go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
$ go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## 编译
```
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    helloworld/helloworld.proto
```

## 感谢
- [go-kiss/sniper](https://github.com/go-kiss/sniper) 一个轻量级 go 业务框架.
- [go-kratos/kratos](https://github.com/go-kratos/kratos) 一套轻量级 go 微服务框架，包含大量微服务相关框架及工具.
- [tal-tech/go-zero](https://github.com/tal-tech/go-zero) 是一个集成了各种工程实践的 web 和 rpc 框架.
- [go-kit/kit](https://github.com/go-kit/kit) is a programming toolkit for building microservices in go.
- [asim/go-micro](https://github.com/asim/go-micro) a distributed systems development framework.
- [google/go-cloud](https://github.com/google/go-cloud) is go cloud development kit.
