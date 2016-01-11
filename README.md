#lhttp http long live server with websocket
###Discribe
lhttp is a http like protocol using websocket to provide long live, 
bulid your IM service quickly scalable without XMPP! 

Everything is customizable.

####Protocol stack:
```go
+--------------------+
|       lhttp        |
+--------------------+
|     websocket      |
+--------------------+
|        TCP         |
+--------------------+
```

####Architecture
```go
        +---------------------------------------+
        |    message center cluster (gnatsd)    |
        +---------------------------------------+
 ........|.................|...............|..................
| +-------------+   +-------------+   +-------------+        | 
| |lhttp server |   |lhttp server |   |lhttp server |   ...  |  lhttp server cluster
| +-------------+   +-------------+   +-------------+        | 
 .....|..........._____|  |___.............|  |_________......
      |          |            |            |            |       <----using websocket link
 +--------+  +--------+   +--------+   +--------+   +--------+   
 | client |  | client |   | client |   | client |   | client |   
 +--------+  +--------+   +--------+   +--------+   +--------+  
```

###Protocol
```go
LHTTP/1.0 Command\r\n                --------start line, define command, and protocol [protocol/version] [command]\r\n
Header1:value\r\n                    --------headers
Header2:value\r\n
\r\n
body                                 --------message body
```
for example:
```go
LHTTP/1.0 chat\r\n
content-type:json\r\n
publish:channel_jack\r\n
\r\n
{
    to:jack,
    from:mike,
    message:hello jack,
    time:1990-1210 5:30:48
}
```
###Usage
1. define your processor, you need combine ```BaseProcessor```
```go
type ChatProcessor struct {
    *lhttp.BaseProcessor
}
```
if you don't like ```BaseProcessor```, define your struct witch must has ```OnOpen(*WsHandler)``` 
```OnClose(*WsHandler)``` method
like this:(don't recommand)
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
then if command is "chat" ChatProcessor will handle it 
3. define your onmessage handle
```go
func (p *ChatProcessor)OnMessage(h *WsHandler) {
    h.Send(h.GetBody())
}
```
###Start websocket server
```go
http.Handler("/echo",lhttp.Handler(lhttp.StartServer))
http.ListenAndServe(":8081")
```
### Example , echo
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
###Test
open  websocketServer and run:
```bash
cd websocketServer
go run test.go
```
as we can see, both of the new header are added and new command are set by the server. 
If we don't set a header or command ,then they will return the same result as they 
requested. 

open an other bash, and run client in websocketClient
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
client2 publish a message by channelID, client1 subscribe it, so client 1 will receive the message.
if client1 send unsubscribe channelID, he will not recevie message any more in channelID

support multiple channelID:
```go
LHTTP/1.0 chat\r\n
subscribe:channelID1 channelID2 channelID3\r\n
\r\n
```
####Using HTTP publish message! 
lhttp support publish message by standard HTTP. 
url: /publish . 
method: POST . 
body: use lhttp publish message as HTTP body.
for example I wan't send a message to who subscribe channel_test by HTTP.
```go
    resp,err := http.POST("https://www.yourserver.com/publish", "text/plain",
    "LHTTP/1.0 chat\r\npublish:channel_test\r\n\r\nhello channel_test guys!")
```
when lhttp server receive this message, will publish whole body to channel_test.


###Upstream
we can use lhttp as a proxy:
```go
LHTTP/1.0 command\r\n
upstream:POST http://www.xxx.com\r\n
\r\n
body
```
lhttp will use hole message as http body, post to http://www.xxx.com
if method is GET, lhttp  send http GET request **ignore lhttp message body**:
MESSAGE:=
```go
LHTTP/1.0 command\r\n
upstream:GET http://www.xxx.com?user=user_a&age=26\r\n
\r\n
body
```

####This case will show you about upstream proxy:
jack use lhttp chat with mike, lhttp is third part module, we can't modify lhttp server but
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
                         V
                   +------------+
                   | http server|  upstream server(http://www.xxx.com/record)
                   +------------+
                   (save chat record)
    
```
jack:
MESSAGE_UPSTREAM:=
```go
LHTTP/1.0 chat\r\n
upstream:POST http://www.xxx.com/record\r\n
publish:channel_mike\r\n
\r\n
hello mike,I am jack
```
mike:
```go
LHTTP/1.0 chat\r\n
subscribe:channel_mike\r\n
\r\n
```
when jack send publish message, not only mike will receive the message, the http server will
also receive it. witch http body is:```MESSAGE_UPSTREAM```, so http serve can do anything about
message include save the record

###Multipart form data
forexample a file upload message, the multipart header record the offset of each data part, 
each part can has it own headers
```go
LHTTP/1.0 upload\r\n
multipart:0 56\r\n
\r\n
content-type:text/json\r\n
\r\n
{filename:file.txt,fileLen:5}
content-type:text/plain\r\n
\r\n
hello
```
```go
content-type:text/json\r\n\r\n{filename:file.txt,fileLen:5}content-type:text/plain\r\n\r\nhello
^                                                          ^
|<---------------------first part------------------------->|<---------second part------------>|
0                                                          56                           
```
why not boundary but use offset? if use boundary lhttp need ergodic hole message, that behaviour 
is poor efficiency. instead we use offset to cut message 

