package script

import "bytes"

func Memcache(ip string, port string) bool {
	if port != "11211" {
		port = "11211"
	}

	addr := ip + ":" + port
	payload := []byte("stats\n")
	resp, err := Tcp(addr, payload)
	if err != nil {
		return false
	}
	if bytes.Contains(resp, []byte("STAT pid")) {
		return true
	}
	return false
}
