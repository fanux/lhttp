package lhttp

type HandlerHub struct {
}

func (h *HandlerHub) Get(connSetID string) *WsHandler {
	return &WsHandler{}
}
func (h *HandlerHub) Add(connSetID string, w *WsHandler) {
}
func (h *HandlerHub) Delete(w *WsHandler) {
}
