### LHTTP header filter develop
We know lhttp has some standard headers like subscribe/publish/unsubscribe/upstream ... But if you want to define your own header for example auth, then you need to add your own filter, this document shows you upstream header develop

* Add the header name in lhttpDefine.go
```go
HEADER_KEY_UPSTREAM = "upstream"
```
* create a file named upstreamFilter.go, add filter module
```go
type upstreamHeadFilter struct {
    *HeadFilterBase
}
```
Here we need combine `HeadFilterBase`, so we need't implements all the interface witch we don't need.

* implements filter method
```go
type HeadFilterHandler interface{
  	OnOpenFilterHandler(*WsHandler)      //when open the link, this method been called.
	BeforeRequestFilterHandle(*WsHandler) //before the service module handle, this method been called
	AfterRequestFilterHandle(*WsHandler)  //after the service module handle,this method been called
	OnCloseFilterHandle(*WsHandler)       //close the link, this method been called.
}
```
for example, after the service handle, upstream module send the message to upstream server,
so we implement `AfterRequestFilterHandle` interface.
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
* Regist upstream filter in headerFilter.go
```go
	RegistHeadFilter(&upstreamHeadFilter{})
```
