// you can define your handle to processing your private header before or after process message
package lhttp

import (
	"container/list"
	"log"
)

var (
	onOpenFilterList        *list.List
	beforeRequestFilterList *list.List
	afterRequestFilterList  *list.List
	onCloseFilterList       *list.List
)

type HeadFilterHandler interface {
	OnOpenFilterHandle(*WsHandler)
	BeforeRequestFilterHandle(*WsHandler)
	AfterRequestFilterHandle(*WsHandler)
	OnCloseFilterHandle(*WsHandler)
}

//define your filter need combine base
/*
type YourFilter struct {
	*HeadFilterBase
}
*/

type HeadFilterBase struct{}

func (*HeadFilterBase) BeforeRequestFilterHandle(ws *WsHandler) {
	log.Print("head base filter before request")
}

func (*HeadFilterBase) AfterRequestFilterHandle(ws *WsHandler) {
	log.Print("head base filter after request")
}

func (*HeadFilterBase) OnOpenFilterHandle(ws *WsHandler) {
	log.Print("head base filter on open")
}

func (*HeadFilterBase) OnCloseFilterHandle(ws *WsHandler) {
	log.Print("head base filter on close")
}

func RegistHeadFilter(h HeadFilterHandler) {
	onOpenFilterList.PushBack(h)
	beforeRequestFilterList.PushBack(h)
	afterRequestFilterList.PushBack(h)
	onCloseFilterList.PushBack(h)
}

func init() {
	onOpenFilterList = list.New()
	beforeRequestFilterList = list.New()
	afterRequestFilterList = list.New()
	onCloseFilterList = list.New()

	RegistHeadFilter(&HeadFilterBase{})
	RegistHeadFilter(&mqHeadFilter{})
	RegistHeadFilter(&upstreamHeadFilter{})
	RegistHeadFilter(&multipartFilter{})
}
