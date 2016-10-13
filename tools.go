package lhttp

import (
	"net/http"
	"strings"
)

var (
	lhttpHost = "localhost"
	lhttpPort = "8081"
)

// Publish message to channel using http
func Publish(channelID []string, command string, header map[string]string, body string) (err error) {
	//HTTP Client
	//Connect to Terminal Module
	//URL:"/message?channel=channelID&channel=channelID&cmd=cmdID
	//from default.toml HOST : terminal_host
	msg := "LHTTP/1.0 " + command + "\r\n"
	msg += "publish:" + strings.Join(channelID, " ") + "\r\n"

	for k, v := range header {
		msg += k + ":" + v + "\r\n"
	}
	msg += "\r\n" + body

	url := "http://" + lhttpHost + ":" + lhttpPort + "/publish"
	_, err = http.Post(url, "text/plain", strings.NewReader(msg))

	return
}
