#lhttp http long live server with websocket
###discribe
lhttp is a http like protocol but using websocket provide long live
###protocol
```go
LHTTP/1.0 Command\r\n                --------start line,define command,and protocol [protocol/version] [command]\r\n
Header1:value\r\n                    --------headers
Header2:value\r\n
\r\n
body                                 --------message body
```
for example:
```go
LHTTP/1.0 chat\r\n
content-type:json\r\n
to:jack\r\n
from:mike\r\n
\r\n
{
message:hello jack,
time:1990-1210 5:30:48
}
```
###usage
1. define your processor,you need combine ```BaseProcessor```
```go
type ChatProcessor struct {
    *lhttp.BaseProcessor
}
```
if you dont like ```BaseProcessor```,define your struct witch must has ```OnOpen(*WsHandler)``` 
```OnClose(*WsHandler)``` method
like this:(dont recommand)
```go
type ChatProcessor struct {
}
func (p ChatProcessor)OnOpen(h *WsHandler) {
    //your logic
}
func (p ChatProcessor)OnClose(h *WsHandler) {
    //your logic
}
func (p ChatProcessor)OnMessage(h *WsHandler) {
    //your logic
}
```
2. regist your processor
```go
lhttp.Regist("chat",&ChatProcessor{&lhttp.BaseProcessor{}})
```
then if command is 'chat' ChatProcessor will handle it 
3. define your onmessage handle
```go
func (p *ChatProcessor)OnMessage(h *WsHandler) {
    h.Send(h.GetBody())
}
```
###start websocket server
```go
http.Handler("/echo",lhttp.Handler(lhttp.StartServer))
http.ListenAndServe(":8081")
```
### example ,echo
```go
type ChatProcessor struct {
    *lhttp.BaseProcessor
}

func (p *ChatProcessor) OnMessage (h *lhttp.WsHandler) {
    log.Print("on message :", h.GetBody())
    h.Send(h.GetBody())
}

func main(){
    lhttp.Regist("chat", &ChatProcessor{&lhttp.BaseProcessor{}})

    http.Handle("/echo",lhttp.Handler(lhttp.StartServer))
    http.ListenAndServe(":8081",nil)
}
```
###test
open  websocketServer and run:
```bash
cd websocketServer
go run test.go
```
as you can see ,server add new header and set new command, if server not change headers or command,
response headers and command will same as request
open an other bash ,and run client in websocketClient
```bash
cd websocketClient
go run test.go
```
###Subscribe/Publish
client1:
```go
LHTTP/1.0 command\r\n
subscribe:channelID\r\n
\r\n
body optional
```
client2:
```go
LHTTP/1.0 command\r\n
publish:channelID\r\n
\r\n
body require
```
client1:
```go
LHTTP/1.0 command\r\n
unsubscribe:channelID\r\n
\r\n
body optional
```
client2 publish a message by channelID, client1 subscribe it,so client 1 will receive the message.
if client1 send unsubscribe channelID,he will not recevie message any more in channelID

support multiple channelID:
```go
LHTTP1.0 chat\r\n
subscribe:channelID1 channelID2 channelID3\r\n
\r\n
```

###Upstream
we can use lhttp as a proxy:
```go
LHTTP/1.0 command\r\n
upstream:post http://www.xxx.com\r\n
\r\n
body
```
lhttp will use hole message as http body,post to http://www.xxx.com
if method is get,lhttp act message as an argument and send http get request:
MESSAGE:=
```go
LHTTP/1.0 command\r\n
upstream:get http://www.xxx.com\r\n
\r\n
body
```
will send http://www.xxx.com?lhttp=MESSAGE

####this case will show you about upstream proxy:
jack use lhttp chat with mike, lhttp is third part module,we cant modify lhttp server but
we want save the chat record, how can we do?

```
        +----+                  +----+
        |jack|                  |mike|
        +----+                  +----+
         |_____________    _______|
                       |  |
                   +------------+
                   |lhttp server|
                   +------------+
                         |(http request with chat record)
                         |
                         V
                   +------------+
                   | http server|  upstream server(http://www.xxx.com/record)
                   +------------+
                   (save chat record)
    
```
jack:
MESSAGE_UPSTREAM:=
```go
LHTTP1.0 chat\r\n
upstream:post http://www.xxx.com/record\r\n
publish:channel_mike\r\n
\r\n
hello mike,I am jack
```
mike:
```go
LHTTP1.0 chat\r\n
subscribe:channel_mike\r\n
\r\n
```
when jack send publish message,not only mike will receive the message,the http server will
also receive it. witch http body is:```MESSAGE_UPSTREAM```,so http serve can do anything about
message include save the record

###Multipart form data
forexample a file upload message,the multipart header record the offset of each data part
```go
LHTTP1.0 upload\r\n
multipart:0 54\r\n
\r\n
content-type:text/json\r\n
{filename:file.txt,fileLen:5}
content-type:text/plain\r\n
hello
```
```go
content-type:text/json\r\n{filename:file.txt,fileLen:5}content-type:text/plain\r\nhello
^                                                      ^
|<-------------------first part----------------------->|<---------second part---------|
0                                                      54                           
```
why not boundary but use offset? if use boundary lhttp need ergodic hole message,that behaviour 
is poor efficiency. instead we use offset to cut message 

