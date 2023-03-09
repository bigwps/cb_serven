package script

import "bytes"

func RsyncUnauth(ip string, port string) bool {
	if port != "873" {
		port = "873"
	}
	addr := ip + ":" + port
	payload := []byte("info\r\n")

	resp, err := Tcp(addr, payload)
	if err != nil {
		return false
	}
	if bytes.Contains(resp, []byte("@RSYNCD")) {
		return true
	}
	return false
}
