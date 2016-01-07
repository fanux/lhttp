// you can define your handle to processing your private header before or after process message
package lhttp

import (
	"log"
	"strings"
)

/*
  |0|1|2|3|4|5|6|7|8|9|10|11|12|13|14|15|16|17|18|19|

  |<---before-------->|<---after process message--->|
  |<------->|<------->|<-framework-->|<--user use-->|
  |framework|user use |   use        |              |
  |use      |         |

  we use RegistHeadFilter(MQ_PRIORITY,&mqHeadFilter{}) to heandle mq headers(publish/subscribe...)
  if you define your private head filter,priority must [5,9] (before handle message) or
  [15,19] after handle message
*/
var headFilterHandler [20]HeadFilterHandle

var (
	PRIORITY_BEFORE_REQUEST = 10
	PRIORITY_AFTER_REQUEST  = 20
)

var (
	//handle subscribe/publish/unsubscribe header
	MQ_PRIORITY = 5
)

var (
	HEADER_KEY_PUBLISH     = "publish"
	HEADER_KEY_SUBSCRIBE   = "subscribe"
	HEADER_KEY_UNSUBSCRIBE = "unsubscribe"
)

type HeadFilterHandle interface {
	HeaderFilter(*WsHandler)
}

func RegistHeadFilter(priority int, h HeadFilterHandle) {
	if headFilterHandler[priority] == nil {
		headFilterHandler[priority] = h
	} else {
		log.Print("regist head filter error")
	}
}

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
