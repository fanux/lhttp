package lhttp

import (
	"container/list"
	"log"
	"net/url"
	"strings"
)

type WsMessage struct {
	//message raw data
	message string

	//message command type
	command string
	//message headers
	headers map[string]string
	//message body
	body string
}

//fill message by command headers and body
func (m *WsMessage) serializeMessage() string {
	m.message = protocolNameWithVersion + " "
	m.message += m.command + CRLF

	for k, v := range m.headers {
		m.message += k + ":" + v + CRLF
	}
	m.message += CRLF + m.body

	return m.message
}

//parse websocket body
func buildMessage(data string) *WsMessage {
	//TODO optimise ,to use builder pattern
	s := data
	message := &WsMessage{message: data}
	message.headers = make(map[string]string, headerMax)
	//parse message

	//parse start line
	i := strings.Index(s, CRLF)
	message.command = s[protocolLength+1 : i]

	//parse hearders
	k := 0
	headers := s[i+2:]
	var key string
	var value string
	//traverse once
	for j, ch := range headers {
		if ch == ':' && key == "" {
			key = headers[k:j]
			k = j + 1
		} else if headers[j:j+2] == CRLF {
			value = headers[k:j]
			k = j + 2

			message.headers[key] = value
			log.Print("parse head key:", key, " value:", value)
			key = ""
		}
		if headers[k:k+2] == CRLF {
			k += 2
			break
		}
	}

	//set body
	message.body = headers[k:]

	return message
}

type WsHandler struct {
	callbacks HandlerCallbacks

	//websocket connection
	conn *Conn

	//receive message
	message *WsMessage

	resp WsMessage

	upstreamURL *url.URL
	//one connection set id map sevel connections
	//connSetID string

	//save multipars datas, it is a list
	multiparts *multipartBlock
}

func (req *WsHandler) GetMultipart() *multipartBlock {
	return req.multiparts
}

//define subscribe callback as a WsHandler method is very very very importent
func (req *WsHandler) subscribeCallback(s string) {
	Message.Send(req.conn, s)
}

func (req *WsHandler) SetCommand(s string) {
	req.resp.command = s
}

func (req *WsHandler) GetCommand() string {
	return req.message.command
}
func (req *WsHandler) GetHeader(hkey string) string {
	if value, ok := req.message.headers[hkey]; ok {
		return value
	} else {
		return ""
	}
}

//if header already exist,update it
func (req *WsHandler) AddHeader(hkey, hvalue string) {
	req.resp.headers = make(map[string]string, headerMax)
	req.resp.headers[hkey] = hvalue
}

func (req *WsHandler) GetBody() string {
	return req.message.body
}

//if response is nil, use request to fill it
func (req *WsHandler) setResponse() {
	if req.resp.command == "" {
		req.resp.command = req.message.command
	}
	if req.resp.headers == nil {
		req.resp.headers = req.message.headers
	}
	if req.resp.body == "" {
		req.resp.body = req.message.body
	}
}

//if you want change command or header ,using SetCommand or AddHeader
func (req *WsHandler) Send(body string) {
	resp := protocolNameWithVersion + " "
	if req.resp.command != "" {
		resp = resp + req.resp.command + CRLF
	} else {
		resp = resp + req.message.command + CRLF
	}

	if req.resp.headers != nil {
		for k, v := range req.resp.headers {
			resp = resp + k + ":" + v + CRLF
		}
	} else {
		for k, v := range req.message.headers {
			resp = resp + k + ":" + v + CRLF
		}
	}
	resp += CRLF + body

	req.resp.message = resp

	//log.Print("send message:", string(req.message.message))

	Message.Send(req.conn, req.resp.message)

	req.resp = WsMessage{command: "", headers: nil, body: ""}
}

type HandlerCallbacks interface {
	OnOpen(*WsHandler)
	OnClose(*WsHandler)
	OnMessage(*WsHandler)
}

type BaseProcessor struct {
}

func (*BaseProcessor) OnOpen(*WsHandler) {
	log.Print("base on open")
}
func (*BaseProcessor) OnMessage(*WsHandler) {
	log.Print("base on message")
}
func (*BaseProcessor) OnClose(*WsHandler) {
	log.Print("base on close")
}

func StartServer(ws *Conn) {
	openFlag := 0

	//init WsHandler,set connection and connsetid
	wsHandler := &WsHandler{conn: ws}

	for {
		var data string
		err := Message.Receive(ws, &data)
		//log.Print("receive message:", string(data))
		if err != nil {
			break
		}

		if len(data) <= protocolLength {
			//TODO how to provide other protocol
			log.Print("TODO provide other protocol")
			continue
		}

		if data[:protocolLength] != protocolNameWithVersion {
			//TODO how to provide other protocol
			log.Print("TODO provide other protocol")
			continue
		}

		wsHandler.message = buildMessage(data)

		var e *list.Element
		//head filter before process message
		for e = beforeRequestFilterList.Front(); e != nil; e = e.Next() {
			e.Value.(HeadFilterHandler).BeforeRequestFilterHandle(wsHandler)
		}

		wsHandler.callbacks = getProcessor(wsHandler.message.command)
		if wsHandler.callbacks == nil {
			wsHandler.callbacks = &BaseProcessor{}
		}
		//log.Print("callbacks:", wsHandler.callbacks.OnMessage)
		//just call once
		if openFlag == 0 {
			for e = onOpenFilterList.Front(); e != nil; e = e.Next() {
				e.Value.(HeadFilterHandler).OnOpenFilterHandle(wsHandler)
			}
			if wsHandler.callbacks != nil {
				wsHandler.callbacks.OnOpen(wsHandler)
			} else {
				//log.Print("error on open is null")
			}
			openFlag = 1
		}
		if wsHandler.callbacks != nil {
			wsHandler.callbacks.OnMessage(wsHandler)
		} else {
			//log.Print("error onmessage is null ")
		}

		//head filter after process message
		for e = afterRequestFilterList.Front(); e != nil; e = e.Next() {
			e.Value.(HeadFilterHandler).AfterRequestFilterHandle(wsHandler)
		}
	}
	defer func() {
		for e := onCloseFilterList.Front(); e != nil; e = e.Next() {
			e.Value.(HeadFilterHandler).OnCloseFilterHandle(wsHandler)
		}
		if wsHandler.callbacks != nil {
			wsHandler.callbacks.OnClose(wsHandler)
		} else {
			//log.Print("error on close is null")
		}
		ws.Close()
	}()
}
