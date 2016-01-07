package lhttp

var message = "LHTTP/1.0 CHAT\r\ncontent-type:json\r\ncontent-length:10\r\n\r\njsonbody"

/*
func TestParseUnparse(t *testing.T) {
	//parse
	w := buildMessage(message)
	if w.command != "CHAT" {
		t.Errorf("command error:%s", w.command)
	}
	if w.body != "jsonbody" {
		t.Errorf("error body:%s", w.body)
	}
	if v, ok := w.headers["content-type"]; ok {
		if v != "json" {
			t.Errorf("error head value:%s", w.headers["content-type"])
		}
	} else {
		t.Errorf("error key")
	}
	if v, ok := w.headers["content-length"]; ok {
		if v != "10" {
			t.Errorf("error head value:%s", w.headers["content-length"])
		}
	} else {
		t.Errorf("error key")
	}

	//unparse
	req := &WsHandler{message: w}
	req.Send("hello")

	if string(req.message.message) != "LHTTP/1.0 CHAT\r\ncontent-length:10\r\ncontent-type:json\r\n\r\nhello" {
		t.Errorf("combine message error:%s", string(req.message.message))
	}

	if req.message.serializeMessage() != "LHTTP/1.0 CHAT\r\ncontent-length:10\r\ncontent-type:json\r\n\r\nhello" {
		t.Errorf("serializeMessage error:%s", req.message.serializeMessage)
	}
}
*/
