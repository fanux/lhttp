package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

/*
for test using HTTP publish message to who subscribe mike channel
1. first start websocket server
2. start websocket client, input subscribe command, then client will subscribe "mike" channel
3. run this test, then client will receive "hello mike"
*/

func main() {
	resp, err := http.Post("http://localhost:8081/publish", "text/plain", strings.NewReader("LHTTP/1.0 chat\r\npublish:mike\r\n\r\nhello mike"))
	if err != nil {
		fmt.Println("error")
		return
	}
	fmt.Println(ioutil.ReadAll(resp.Body))
}
