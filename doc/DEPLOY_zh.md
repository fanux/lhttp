### 部署文档
### 消息中心依赖gnatsd进程
```go
go get github.com/nats-io/nats

go get github.com/nats-io/gnatsd
```
### 运行gnatsd进程

### LHTTP仅是个框架，main函数需要自己写，不过websocketServer目录下面提供了一个示例
```go
cd websocketServer
go run test.go
```
可根据需要修改test.go,或者开发自己的业务模块。
