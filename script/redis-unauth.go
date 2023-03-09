package script

import (
	"bytes"
)

func RedisUnauth(ip string, port string) bool {
	if port == "" {
		port = "6379"
	}

	addr := ip + ":" + port

	payload := []byte("*1\r\n$4\r\ninfo\r\n")
	resp, err := Tcp(addr, payload)
	if err != nil {
		return false
	}
	if bytes.Contains(resp, []byte("redis_version")) {
		return true
	}
	return false
}
