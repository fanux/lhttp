package lhttp

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/nats-io/nats"
)

type MqHandler interface{}

type Mq struct {
	conn *nats.EncodedConn
}

func (mq *Mq) Publish(key string, v MqHandler) error {
	return mq.conn.Publish(key, v)
}

func (mq *Mq) Subscribe(key string, v MqHandler) (*nats.Subscription, error) {
	return mq.conn.Subscribe(key, v)
}

/*
func (mq *Mq) Unsubscribe(sub *nats.Subscription) error {
	return sub.Unsubscribe()
}
*/

func (mq *Mq) Unsubscribe(channel string) error {
	sub, _ := mq.Subscribe(channel, nil)
	return sub.Unsubscribe()
}

var mq Mq

func NewMq() *Mq {
	return &mq
}

type httpPublisher struct{}

func (*httpPublisher) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Print("an error occur")
	}
	log.Print("http publish body: ", string(body))

	bodyStr := string(body)

	message := buildMessage(bodyStr)

	channels, ok := message.headers[HEADER_KEY_PUBLISH]
	if !ok {
		log.Print("cant get Publish header")
		return
	}

	for _, c := range strings.Split(channels, " ") {
		mq.Publish(c, bodyStr)
	}

	req.Body.Close()
}

func init() {
	nc, _ := nats.Connect(nats.DefaultURL)
	c, err := nats.NewEncodedConn(nc, nats.DEFAULT_ENCODER)
	if err != nil {
		log.Print("mq init error")
	} else {
		mq.conn = c
	}

	//handle http publish message
	http.Handle("/publish", &httpPublisher{})
}
