### UML

```go
+----------------------+       +------------------------------+              +---------------------+
|     BaseProcessor    |       |    Context                   |              |      LHTTPMsg       |
+----------------------+       +------------------------------+              +---------------------+
|+OnOpen(c *Context)   |       |  -WebsocketConn              |              |-rawMessage str      |
|+OnMessage(c *Context)|------>|  -req LHTTPMsg               |<>----------->|-command    str      |
|+OnClose(c *Context)  |       |  -resp *LHTTPMsg             |              |-headers map[str]str |
+----------------------+       |  -upstreamURL URL            |              |-body       str      |
                               |  -multiparts *multipartBlock |              +---------------------+
                               +------------------------------+
                               |+SetCommand(str)              |
                               |+GetCommand(str)              |
                               |+GetHeader(str)str            |
                               |+AddHeader(key str, value str)|
                               |+GetBody()str                 |
                               |+Send(body string)            |
                               |+GetMultipart()*multipartBlock|
                               |            ...               |
                               +------------------------------+
```

#### multipartBlock is a list
```go
+---------------+         +---------------+
|multipartBlock |    +--->|multipartBlock |
+---------------+    |    +---------------+
|-headers       |    |    |-headers       |
|-body          |    |    |-body          |
|-nextBlock----------+    |-nextBlock--------    ...    
+---------------+         +---------------+
```
