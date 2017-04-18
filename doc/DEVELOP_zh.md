### LHTTP过滤模块开发
所谓过滤模块就是处理特定的LHTTP header的模块。
此处以upstream过滤模块为事例描述一个过滤模块开发的过程。
1. 在lhttpDefine.go文件中定义header名称
```go
HEADER_KEY_UPSTREAM = "upstream"
```
2. 创建一个upstreamFilter.go文件，编写过滤模块
```go
type upstreamHeadFilter struct {
    *HeadFilterBase
}
```
这里必须组合HeadFilterBase,这样你自定义的过滤模块就可以介入LHTTP请求的四个阶段
3. 实现过滤方法
```go
type HeadFilterHandler interface{
    OnOpenFilterHandler(*WsHandler)      //在打开链接时处理
	BeforeRequestFilterHandle(*WsHandler) //在业务模块处理前处理
	AfterRequestFilterHandle(*WsHandler)  //在业务模块处理后处理
	OnCloseFilterHandle(*WsHandler)       //在关闭连接时处理
}
```
比如upstream过滤模块就需要在请求处理后将信息发送给上游服务器，
所以我们为upstream模块定义AfterRequestFilterHandle方法
```go
func (*upstreamHeadFilter) AfterRequestFilterHandle(ws *WsHandler) {
	var value string
	if value = ws.GetHeader(HEADER_KEY_UPSTREAM); value == "" {
		log.Print("no upstream header found:", ws.message.message, ws.message.headers)
		return
	}
    //....
    //send message to upstream server
}
```
4. 将upstream过滤器注册进去headerFilter.go
```go
	RegistHeadFilter(&upstreamHeadFilter{})
```
