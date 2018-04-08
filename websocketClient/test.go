package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "115.28.143.67:8081", "http service address")

var (
	MESSAGE_CHAT = "LHTTP/1.0 chat\r\ncontent-type:json\r\ncontent-length:10\r\n\r\njsonbody"
	MESSAGE_SUB  = "LHTTP/1.0 subpub\r\nsubscribe:channel_mike\r\n\r\n"
	MESSAGE_PUB  = "LHTTP/1.0 subpub\r\npublish:channel_mike\r\n\r\nhello mike, I am jack"
)

func sendSubscribe(c *websocket.Conn) {
	/*
		fmt.Println("input your channels split by blank:")
		in := make([]byte, 1024)
		os.Stdin.Read(in)
	*/
	msg := "LHTTP/1.0 subpub\r\nsubscribe:camera_123" + "\r\n\r\n"
	c.WriteMessage(websocket.TextMessage, []byte(msg))
}
func sendUnsubscribe(c *websocket.Conn) {
	/*
		fmt.Println("input your channels split by blank:")
		in := make([]byte, 1024)
		os.Stdin.Read(in)
	*/
	msg := "LHTTP/1.0 subpub\r\nunsubscribe:camera_123" + "\r\n\r\n"
	c.WriteMessage(websocket.TextMessage, []byte(msg))
}
func sendPublish(c *websocket.Conn) {
	/*
		fmt.Println("input your channels split by blank:")
		in := make([]byte, 1024)
		os.Stdin.Read(in)
	*/
	msg := "LHTTP/1.0 subpub\r\npublish:camera_123" + "\r\n\r\nhello world"
	c.WriteMessage(websocket.TextMessage, []byte(msg))
}

func sendUpstream(c *websocket.Conn) {
	msg := "LHTTP/1.0 upstream\r\nupstream:GET http://115.28.143.67:8080/v1/user/10000/reports?AlgID=1\r\n\r\n"
	c.WriteMessage(websocket.TextMessage, []byte(msg))
	msg = "LHTTP/1.0 upstream\r\nupstream:POST http://115.28.143.67:8080/v1/camera/10000/alarms\r\n\r\n{}"
	c.WriteMessage(websocket.TextMessage, []byte(msg))
}

func sendMultipart(c *websocket.Conn) {
	msg := "LHTTP/1.0 upload\r\nmultipart:0 14\r\n\r\nk1:v1\r\n\r\nbody1k2:v2\r\n\r\nbody2"
	c.WriteMessage(websocket.TextMessage, []byte(msg))
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer c.Close()
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case _ = <-ticker.C:
			//err := c.WriteMessage(websocket.TextMessage, []byte(MESSAGE_CHAT))
			fmt.Println("input your commands:multipart subscribe publish unsubscribe upstream")
			input := make([]byte, 1024)
			os.Stdin.Read(input)
			if strings.HasPrefix(string(input), "subscribe") {
				sendSubscribe(c)
			}
			if strings.HasPrefix(string(input), "unsubscribe") {
				sendUnsubscribe(c)
			}
			if strings.HasPrefix(string(input), "publish") {
				sendPublish(c)
			}
			if strings.HasPrefix(string(input), "upstream") {
				sendUpstream(c)
			}
			if strings.HasPrefix(string(input), "multipart") {
				sendMultipart(c)
			}
			//err := c.WriteMessage(websocket.TextMessage, []byte(MESSAGE_SUB))
			//err2 := c.WriteMessage(websocket.TextMessage, []byte(MESSAGE_PUB))
			//if err != nil || err2 != nil {
			//	log.Println("write:", err)
			//	return
			//}
		case <-interrupt:
			log.Println("interrupt")
			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			c.Close()
			return
		}
	}
}
