###Lhttp filter module develop
####define your filter, for example message queue filter: 
message queue filter care about "publish subscribe unsubscribe" headers.

```go
//if client send message include subscribe/publish/unsubscribe header
//this filter work,use nats as a message queue client
type mqHeadFilter struct {
}

func (*mqHeadFilter) HeaderFilter(ws *WsHandler) {
	var value string
	var channels []string

	if value = ws.GetHeader(HEADER_KEY_SUBSCRIBE); value != "" {
		channels = strings.Split(value, " ")
		for _, c := range channels {
			mq.Subscribe(c, ws.subscribeCallback)
			log.Print("subscribe channel: ", c)
		}
	}

	if value = ws.GetHeader(HEADER_KEY_PUBLISH); value != "" {
		channels = strings.Split(value, " ")
		for _, c := range channels {
			ws.setResponse()
			ws.resp.serializeMessage()
			mq.Publish(c, ws.resp.message)

			log.Print("publish channel: ", c, "message:", ws.resp.message)
		}
	}

	if value = ws.GetHeader(HEADER_KEY_UNSUBSCRIBE); value != "" {
		channels = strings.Split(value, " ")
		for _, c := range channels {
			mq.Unsubscribe(c)
		}
	}
}
```
####regist your filter
```go
RegistHeadFilter(MQ_PRIORITY, &mqHeadFilter{})
```
####something about filter priority
```go
  |0|1|2|3|4|5|6|7|8|9|10|11|12|13|14|15|16|17|18|19|

  |<---before-------->|<---after process message--->|
  |<------->|<------->|<-framework-->|<--user use-->|
  |framework|user use |   use        |              |
  |use      |         |

  we use RegistHeadFilter(MQ_PRIORITY,&mqHeadFilter{}) to heandle mq headers(publish/subscribe...)
  if you define your private head filter,priority must [5,9] (before handle message) or
  [15,19] after handle message
```
