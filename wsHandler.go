package lhttp

import (
	"log"
	"strings"
)

var (
	connID = "connid"
	sign   = "sign"
	//headers max num not size
	headerMax               = 20
	version                 = "1.0"
	protocolName            = "LHTTP"
	protocolNameWithVersion = "LHTTP/1.0"
	protocolLength          = 9
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
	m.message = "LHTTP/1.0 "
	m.message += m.command + "\r\n"

	for k, v := range m.headers {
		m.message += k + ":" + v + "\r\n"
	}
	m.message += "\r\n" + m.body

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
	i := strings.Index(s, "\r\n")
	message.command = s[protocolLength+1 : i]

	//parse hearders
	k := 0
	headers := s[i+2:]
	var key string
	var value string
	//traverse once
	for j, ch := range headers {
		if ch == ':' {
			key = headers[k:j]
			k = j + 1
		} else if headers[j:j+2] == "\r\n" {
			value = headers[k:j]
			k = j + 2

			message.headers[key] = value
		}
		if headers[k:k+2] == "\r\n" {
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

	//one connection set id map sevel connections
	connSetID string
}

func (req *WsHandler) subscribeCallback(s string) {
	log.Println("===========", s)

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
	resp := "LHTTP/1.0 "
	if req.resp.command != "" {
		resp = resp + req.resp.command + "\r\n"
	} else {
		resp = resp + req.message.command + "\r\n"
	}

	if req.resp.headers != nil {
		for k, v := range req.resp.headers {
			resp = resp + k + ":" + v + "\r\n"
		}
	} else {
		for k, v := range req.message.headers {
			resp = resp + k + ":" + v + "\r\n"
		}
	}
	resp += "\r\n" + body

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

func registAllHeadFilter() {
	RegistHeadFilter(MQ_PRIORITY, &mqHeadFilter{})
}

func StartServer(ws *Conn) {
	registAllHeadFilter()

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

		//head filter before process message
		for _, h := range headFilterHandler[:PRIORITY_BEFORE_REQUEST] {
			if h != nil {
				h.HeaderFilter(wsHandler)
			}
		}

		wsHandler.callbacks = getProcessor(wsHandler.message.command)
		//log.Print("callbacks:", wsHandler.callbacks.OnMessage)
		//just call once
		if openFlag == 0 {
			if wsHandler.callbacks.OnOpen != nil {
				wsHandler.callbacks.OnOpen(wsHandler)
			} else {
				//log.Print("error on open is null")
			}
			openFlag = 1
		}
		if wsHandler.callbacks.OnMessage != nil {
			wsHandler.callbacks.OnMessage(wsHandler)
		} else {
			//log.Print("error onmessage is null ")
		}

		//head filter after process message
		for _, h := range headFilterHandler[PRIORITY_BEFORE_REQUEST:] {
			if h != nil {
				h.HeaderFilter(wsHandler)
			}
		}
	}
	defer func() {
		if wsHandler.callbacks.OnClose != nil {
			wsHandler.callbacks.OnClose(wsHandler)
		} else {
			//log.Print("error on close is null")
		}
		ws.Close()
	}()
}
