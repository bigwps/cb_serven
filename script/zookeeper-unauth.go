package script

import "bytes"

func Zookeeper(ip string, port string) bool {
	if port != "2181" {
		port = "2181"
	}
	addr := ip + ":" + port
	payload := []byte("envidddfdsfsafafaerwrwerqwe")
	resp, err := Tcp(addr, payload)
	if err != nil {
		return false
	}
	if bytes.Contains(resp, []byte("Environment")) {
		return true
	}
	return false
}
