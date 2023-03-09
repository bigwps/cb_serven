package script

import (
	"net"
	"time"
)

func Tcp(addr string, payload []byte) ([]byte, error) {

	// 连接超时
	timeout := time.Duration(6) * time.Second
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	// 读写超时
	conn.SetReadDeadline(time.Now().Add(timeout))
	// 发送数据
	_, err = conn.Write(payload)
	if err != nil {
		return nil, err
	}
	// 创建一个名为 buf 的字节切片，长度为 20480。它将在下一行代码中用于存储从连接读取的数据：
	buf := make([]byte, 20480)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], err
}
