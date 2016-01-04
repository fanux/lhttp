package main

import (
	"log"
	"mine/Lhttp"
	"net/http"
)

type ChatProcessor struct {
	*lhttp.BaseProcessor
}

func (p *ChatProcessor) OnMessage(h *lhttp.WsHandler) {
	log.Print("on OnMessage: ", h.GetBody())
	h.AddHeader("content-type", "image/png")
	h.SetCommand("auth")
	h.Send(h.GetBody())
}

func main() {
	lhttp.Regist("chat", &ChatProcessor{&lhttp.BaseProcessor{}})
	http.Handle("/echo", lhttp.Handler(lhttp.StartServer))
	http.ListenAndServe(":8081", nil)
}
