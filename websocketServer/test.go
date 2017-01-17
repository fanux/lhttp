package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fanux/lhttp"
)

//ChatProcessor is
type ChatProcessor struct {
	*lhttp.BaseProcessor
}

//OnMessage is
func (p *ChatProcessor) OnMessage(h *lhttp.WsHandler) {
	log.Print("on OnMessage: ", h.GetBody())
	h.AddHeader("content-type", "image/png")
	h.SetCommand("auth")
	h.Send(h.GetBody())
}

//SubPubProcessor is
type SubPubProcessor struct {
	*lhttp.BaseProcessor
}

//UpstreamProcessor is
type UpstreamProcessor struct {
	*lhttp.BaseProcessor
}

//UploadProcessor is
type UploadProcessor struct {
	*lhttp.BaseProcessor
}

//OnMessage is
func (*UploadProcessor) OnMessage(ws *lhttp.WsHandler) {
	for m := ws.GetMultipart(); m != nil; m = m.GetNext() {
		log.Print("multibody:", m.GetBody(), " headers:", m.GetHeaders())
	}
}

func main() {
	lhttp.Regist("chat", &ChatProcessor{&lhttp.BaseProcessor{}})
	lhttp.Regist("subpub", &SubPubProcessor{&lhttp.BaseProcessor{}})
	lhttp.Regist("upstream", &UpstreamProcessor{&lhttp.BaseProcessor{}})
	lhttp.Regist("upload", &UploadProcessor{&lhttp.BaseProcessor{}})

	http.Handle("/echo", lhttp.Handler(lhttp.StartServer))
	http.Handle("/", lhttp.Handler(lhttp.StartServer))
	http.HandleFunc("/https", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", "world")
	})
	http.ListenAndServe(":8581", nil)
}
