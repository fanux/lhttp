package lhttp

import (
	"log"
	"strings"
)

var (
	connID = "connid"
	sign   = "sign"
	//headers max num not size
	headerMax = 20
)

type WsMessage struct {
	//message raw data
	message []byte

	//message command type
	command string
	//message headers
	headers map[string]string
	//message body
	body string
}

//parse websocket body
func buildMessage(data []byte) *WsMessage {
	s := string(data)
	message := &WsMessage{message: data}
	message.headers = make(map[string]string, headerMax)
	//parse message

	//parse start line
	i := strings.Index(s, "\r\n")
	message.command = s[:i]

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

func (req *WsHandler) SetCommand(s string) {
	req.resp.command = s
}

func (req *WsHandler) GetCommand() string {
	return req.message.command
}
func (req *WsHandler) GetHeader(hkey string) string {
	return req.message.headers[hkey]
}

//if header already exist,update it
func (req *WsHandler) AddHeader(hkey, hvalue string) {
	req.resp.headers = make(map[string]string, headerMax)
	req.resp.headers[hkey] = hvalue
}

func (req *WsHandler) GetBody() string {
	return req.message.body
}

//if you want change command or header ,using SetCommand or AddHeader
func (req *WsHandler) Send(body string) {
	var resp string
	if req.resp.command != "" {
		resp = req.resp.command + "\r\n"
	} else {
		resp = req.message.command + "\r\n"
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

	req.resp.message = []byte(resp)

	//log.Print("send message:", string(req.message.message))

	Message.Send(req.conn, req.resp.message)

	req.resp = WsMessage{command: "", headers: nil}
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
	//log.Print("start serve")
	openFlag := 0

	//init WsHandler,set connection and connsetid
	//TODO auth
	id := "123"
	wsHandler := &WsHandler{conn: ws, connSetID: id}

	for {
		var data []byte
		err := Message.Receive(ws, &data)
		//log.Print("receive message:", string(data))
		if err != nil {
			break
		}
		wsHandler.message = buildMessage(data)

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
