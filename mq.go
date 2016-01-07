package lhttp

import (
	"log"

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

func init() {
	nc, _ := nats.Connect(nats.DefaultURL)
	c, err := nats.NewEncodedConn(nc, nats.DEFAULT_ENCODER)
	if err != nil {
		log.Print("mq init error")
	} else {
		mq.conn = c
	}
}
