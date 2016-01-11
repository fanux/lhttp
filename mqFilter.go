package lhttp

import (
	"log"
	"strings"
)

//if client send message include subscribe/publish/unsubscribe header
//this filter work,use nats as a message queue client
type mqHeadFilter struct {
	*HeadFilterBase
}

func (*mqHeadFilter) AfterRequestFilterHandle(ws *WsHandler) {
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
