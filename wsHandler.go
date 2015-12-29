package lhttp

type WsResponse struct {
	command string
	headers map[string]string
	body    string

	req *WsHandler
}

//response to client
func (rep *WsResponse) response() {
}

type WsHandler struct {
	callbacks HandlerCallbacks

	//TODO change real websocket connection
	conn string

	//one connection set id map sevel connections
	connSetID string

	resp *WsResponse
}

func (req *WsHandler) SetCommand(s string) {
}
func (req *WsHandler) GetCommand() string {
	return ""
}
func (req *WsHandler) GetHeader(hkey string) string {
	return ""
}
func (req *WsHandler) AddHeader(hkey, hvalue string) {
}
func (req *WsHandler) GetBody() string {
	return ""
}
func (req *WsHandler) Send(string) {
}

type HandlerCallbacks interface {
	onOpen(*WsHandler)
	onClose(*WsHandler)
	onMessage(*WsHandler)
}

type BaseProcessor struct {
}

func (*BaseProcessor) onOpen(*WsHandler) {
}
func (*BaseProcessor) onClose(*WsHandler) {
}
func (*BaseProcessor) onMessage(*WsHandler) {
}
