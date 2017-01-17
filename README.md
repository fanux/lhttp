#Your star is my power!! :rocket: :star: :star: :star: :star: :star:

[![License MIT](https://img.shields.io/npm/l/express.svg)](http://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/fanux/lhttp)](https://goreportcard.com/report/github.com/fanux/lhttp) [![GoDoc](https://godoc.org/github.com/fanux/lhttp?status.svg)](http://godoc.org/github.com/fanux/lhttp) 
[![Awesome](https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg)](https://github.com/avelino/awesome-go/blob/master/README.md#networking) 

###Discribe
lhttp is a http like protocol using websocket to provide long live, 
build your IM service quickly scalable without XMPP! 

Everything is customizable.

###[简体中文](https://github.com/fanux/lhttp/blob/master/doc/README_zh.md)

###Features
*   simple easy but powerful!
*   fast, publish 10000 messages using 0.04s(single-core CPU,1G memory).
*   support cluster.
*   easy to customize and expansion.
*   work well with HTTP. So LHTTP can work with others language like PHP java python etc,.

### A simple [chat room demo](https://github.com/fanux/lhttp-web-demo)
![chat-demo](https://github.com/fanux/lhttp-web-demo/blob/master/web-demo.gif)
with [lhttp javascript sdk](https://github.com/fanux/lhttp-javascript-sdk) we complete a simple chat room within 40 lines code!!

###SDKs 
- [x] [javascript SDK](https://github.com/fanux/lhttp-javascript-sdk) webapp or website.
- [ ] [c SDK](https://github.com/fanux/lhttp-c-sdk) ARM application or some c/c++ application.

###[Header filter development](https://github.com/fanux/lhttp/blob/master/doc/DEVELOP.md)

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

####Quick start
```bash
go get github.com/nats-io/nats
go get github.com/fanux/lhttp
```
We need run gnatsd first:
```bash
cd bin
./gnatsd &
./lhttpd 
```
Open anohter bash run lhttpClient, then input your command:
```bash
cd bin
./lhttpClient
```

###Ship on docker
```
$ docker build -t lhttp:latest .
$ docker run -p 9090:9090 -p 8081:8081 lhttp:latest
```
Open two windows in your browser, enter `http://localhost:9090`.

Lhttp server port is 8081, your own websocket client can connect to `ws://localhost:8081`

Enjoy the chat...

Alternative, pull image from docker hub.
```
$ docker run -p 9090:9090 -p 8081:8081 fanux/lhttp:latest
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
 > define your processor, you need combine ```BaseProcessor```
 
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

> regist your processor

```go
lhttp.Regist("chat",&ChatProcessor{&lhttp.BaseProcessor{}})
```
**then if command is "chat" ChatProcessor will handle it** 

> define your onmessage handle

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
***

###Test
open  websocketServer and run:
```bash
cd websocketServer
go run test.go
```
as we can see, both of the new headers are added and new command is set by the server. 
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
if client1 send unsubscribe channelID, he will not receive message any more in channelID

support multiple channelID:
```go
LHTTP/1.0 chat\r\n
subscribe:channelID1 channelID2 channelID3\r\n
\r\n
```
####Using HTTP publish message! 
lhttp support publish message by standard HTTP. 
URL: /publish . 
method: POST . 
body: use lhttp publishes message as HTTP body.
for example I want send a message to who subscribe channel_test by HTTP.
```go
    resp,err := http.POST("https://www.yourserver.com/publish", "text/plain",
    "LHTTP/1.0 chat\r\npublish:channel_test\r\n\r\nhello channel_test guys!")
```
when lhttp server receive this message, will publish whole body to channel_test.

your can use ```Publish``` function in tools.go
```go
//func Publish(channelID []string, command string, header map[string]string, body string) (err error) {
//}
//send message to who subscribe mike.

Publish("mike", "yourCommand", nil, "hello mike!")
```


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

```go
LHTTP/1.0 command\r\n
upstream:GET http://www.xxx.com?user=user_a&age=26\r\n
\r\n
body
```

####This case will show you about upstream proxy:
jack use lhttp chat with mike, lhttp is third part module, we can't modify lhttp server but
we want to save the chat record, how can we do?

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
`MESSAGE_UPSTREAM`
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
also receive it. witch http body is:```MESSAGE_UPSTREAM```, so http server can do anything about
message include save the record

###Multipart data
for example a file upload message, the multipart header record the offset of each data part, 
each part can have it own headers
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

####How to get multipart data
for example this is client message.
```go
LHTTP/1.0 upload\r\nmultipart:0 14\r\n\r\nk1:v1\r\n\r\nbody1k2:v2\r\n\r\nbody2
```
server code:
```go
type UploadProcessor struct {
	*lhttp.BaseProcessor
}

func (*UploadProcessor) OnMessage(ws *lhttp.WsHandler) {
	for m := ws.GetMultipart(); m != nil; m = m.GetNext() {
		log.Print("multibody:", m.GetBody(), " headers:", m.GetHeaders())
	}
}

//don't forget to tegist your command processor

lhttp.Regist("upload", &UploadProcessor{&lhttp.BaseProcessor{}})
```

## Partners
[![](https://yunbi.com/logos/logo.svg)](https://yunbi.com)
