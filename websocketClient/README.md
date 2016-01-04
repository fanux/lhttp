send data to server:
```go
LHTTP/1.0 chat\r\n
content-type:json\r\n
content-length:10\r\n
\r\n
jsonbody
```

if server dont add header and set command,client will receive:
```go
LHTTP/1.0 chat\r\n
content-type:json\r\n
content-length:10\r\n
\r\n
jsonbody
```
else client will recevie new content,in this case ,client will receive:
```go
LHTTP/1.0 auth\r\n
content-type:image/png\r\n
\r\n
jsonbody
```
