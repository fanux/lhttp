#lhttp http long live server with websocket
###discribe
lhttp is a http like protocol but using websocket provide long live
###protocol
```go
Command\r\n                          --------start line,define command
Header1:value\r\n                    --------headers
Header2:value\r\n
\r\n
body                                 --------message body
```
for example:
```go
chat\r\n
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
    BaseProcessor
}
```
2. regist your processor
```go
Regist('chat',&ChatProcessor{})
```
then if command is 'chat' ChatProcessor will handle it 
3.define your onmessage handle
```go
func (p *ChatProcessor)onMessage(h *WsHandler) {
    h.Send(h.GETbody())
}
```
