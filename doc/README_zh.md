### 描述

Lhttp是一个基于websocket服务端框架，提供一个类似http的协议去帮助开发者开发长连接的应用。

使用Lhttp可以大量减少服务端开发的工作量，实现非常好的模块化和业务功能的解耦合。

可以定制任何你想要的功能。

### 特点
*   使用简单，功能强大
*   性能高，使用gnatsd消息队列 publish 10000 条消息耗时0.04s(single-core CPU,1G memory).
*   支持集群，横向扩展，通过增加服务器来获取更高的服务能力
*   非常容器进行定制与扩展
*   可以非常好的与http服务协同工作，如利用http发送消息，将消息转发给上游http服务器等。所以即便你不会go语言也可以开发一些应用。

###  [聊天室demo](https://github.com/fanux/lhttp-web-demo)
![chat-demo](https://github.com/fanux/lhttp-web-demo/blob/master/web-demo.gif)
[前端sdk](https://github.com/fanux/lhttp-javascript-sdk)


#### 协议栈:
```go
+--------------------+
|       lhttp        |
+--------------------+
|     websocket      |
+--------------------+
|        TCP         |
+--------------------+
```

#### 系统架构
```go
        +---------------------------------------+
        |    message center cluster (gnatsd)    |
        +---------------------------------------+
 ........|.................|...............|..................
| +-------------+   +-------------+   +-------------+        | 
| |lhttp server |   |lhttp server |   |lhttp server |   ...  |  lhttp 服务集群
| +-------------+   +-------------+   +-------------+        | 
 .....|..........._____|  |___.............|  |_________......
      |          |            |            |            |       <----使用websocket链接
 +--------+  +--------+   +--------+   +--------+   +--------+   
 | client |  | client |   | client |   | client |   | client |   
 +--------+  +--------+   +--------+   +--------+   +--------+  
```

#### 快速入门
```bash
go get github.com/nats-io/nats
go get github.com/fanux/lhttp
```
先启动gnatsd:
```bash
cd bin
./gnatsd &
./lhttpd 
```

打开另一个终端，执行客户端程序，输入命令码：
```bash
cd bin
./lhttpClient
```

### 使用docker快速体验
```
$ docker build -t lhttp:latest .
$ docker run -p 9090:9090 -p 8081:8081 lhttp:latest
```
打开浏览器，访问： `http://localhost:9090`.

打开两个窗口就可以聊起来了。

websocket 端口是 8081, 可以使用自己的websocket客户端去连 `ws://localhost:8081`

也可以从dockerhub上下载镜像:
```
$ docker run -p 9090:9090 -p 8081:8081 fanux/lhttp:latest
```

### 协议介绍
```go
LHTTP/1.0 Command\r\n                --------起始行，协议名和版本，Command:非常重要，标识这条消息的命令码是什么，服务端也是根据命令码注册对应的处理器的。
Header1:value\r\n                    --------首部
Header2:value\r\n
\r\n
body                                 --------消息体
```

事例：
```go
LHTTP/1.0 chat\r\n                  命令码叫`chat`
content-type:json\r\n               消息体使用json编码
publish:channel_jack\r\n            服务端请把这条消息publish给jack (jack订阅了channel_jack)
\r\n
{
    to:jack,
    from:mike,
    message:hello jack,
    time:1990-1210 5:30:48
}
```
### 使用教程,只需三步
 > 定义你的处理器，需要聚合 ```BaseProcessor```
 
```go
type ChatProcessor struct {
    *lhttp.BaseProcessor
}
```

> 实现三个接口，连接打开时干嘛，关闭时干嘛，消息到来时干嘛。
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

> 注册你的处理器，这里的`chat` 与消息体中的`chat对应`,也就是这个处理器仅会处理`LHTTP/1.0 chat\r\n....`这类消息.
```go
lhttp.Regist("chat",&ChatProcessor{&lhttp.BaseProcessor{}})
```
**then if command is "chat" ChatProcessor will handle it** 

> 这里比如收到消息就直接将消息返回：
```go
func (p *ChatProcessor)OnMessage(h *WsHandler) {
    h.Send(h.GetBody())
}
```
### 启动服务器
```go
http.Handler("/echo",lhttp.Handler(lhttp.StartServer))
http.ListenAndServe(":8081")
```

### 一个完整的回射例子：
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

### 订阅/发布 
下面来看用Lhttp开发及时通信应用有多简单

假设有两个客户端，这里的客户端比如浏览器应用。

client1:
```go
LHTTP/1.0 command\r\n
subscribe:channelID\r\n
\r\n
body optional
```
client1通过websocket向Lhttp发送如上字符串，就订阅了`channelId`

client2:
```go
LHTTP/1.0 command\r\n
publish:channelID\r\n
\r\n
body require
```
client2通过websocket向Lhttp发送如上字符串，就向`channelID`发布了一条消息。  因为client1订阅了channelID,所以client1会收到这条消息。


client1不想再收消息，那么发如下字符串给服务端即可:
```go
LHTTP/1.0 command\r\n
unsubscribe:channelID\r\n
\r\n
body optional
```
订阅/发布 是lhttp内置功能，服务端一行代码不用写即可获取这种服务，只需要使用特定首部`subscribe`,`publish` 和`unsubscribe`

同时订阅多个，如同时订阅多个聊天室。
```go
LHTTP/1.0 chat\r\n
subscribe:channelID1 channelID2 channelID3\r\n
\r\n
```
#### 使用http发布消息
URL: /publish . 
方法: POST . 
http body: 整个lhttp消息
for example I want send a message to who subscribe channel_test by HTTP.
如我想发送一条消息给订阅了channel_test的人。
```go
    resp,err := http.POST("https://www.yourserver.com/publish", "text/plain",
    "LHTTP/1.0 chat\r\npublish:channel_test\r\n\r\nhello channel_test guys!")
```

这里封装好了一个更好用的工具 ```Publish```  tools.go
```go
//func Publish(channelID []string, command string, header map[string]string, body string) (err error) {
//}
//send message to who subscribe mike.

Publish("mike", "yourCommand", nil, "hello mike!")
```

### 上游服务器
upstream首部可以让lhttp向上游的http服务器发送一条消息。
```go
LHTTP/1.0 command\r\n
upstream:POST http://www.xxx.com\r\n
\r\n
body
```
如果是POST方法，lhttp会把整个消息体当作http的body发送给 http://www.xxx.com
如果是GET，lhttp会忽略消息体

```go
LHTTP/1.0 command\r\n
upstream:GET http://www.xxx.com?user=user_a&age=26\r\n
\r\n
body
```

#### upstream有什么用：
如我们不想改动lhttp的代码，但是想存储聊天记录。

通过upstream可以实现很好的解耦：

并且http server可以用其它语言实现.

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
这样jack publish消息时不仅mike可以收到，后端的upstream server也可以收到，我们可以在后端服务器中处理消息存储的逻辑，如将消息

存储到redis的有序集合中。


### 分块消息
试想一下，一条消息中既有图片也有文字还有语音怎么办？ lhttp的multipart首部解决这个问题

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
http中是使用boundry实现的，lhttp使用偏移量标识分块，这样效率更高，不需要遍历整个消息体。

#### 如何获取分块消息
如客户端消息如下：
```go
LHTTP/1.0 upload\r\nmultipart:0 14\r\n\r\nk1:v1\r\n\r\nbody1k2:v2\r\n\r\nbody2
```
服务端代码，消息存在链表中：
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
### [首部过滤模块开发](https://github.com/fanux/lhttp/blob/master/doc/DEVELOP.md)

## Partners
[![](https://yunbi.com/logos/logo.svg)](https://yunbi.com)

