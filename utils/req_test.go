package utils

import (
	"fmt"
	"github.com/imroc/req/v3"
	"log"
	"testing"
)

func TestHelloWorld(t *testing.T) {
	t.Log("hello world")
	_, client := InitReq2("")
	client.R().Method = "POST"
	res, err := client.R().Get("http://www.baidu.com")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}

func InitReq2(proxy string) (error, *req.Client) {
	var err error
	client := req.C()
	if proxy != "" {
		client.SetProxyURL(proxy)
	}
	return err, client
}
