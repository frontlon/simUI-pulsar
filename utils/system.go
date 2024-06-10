package utils

import (
	"fmt"
	"net"
)

// 检查端口是否占用
func IsPortInUse(port int) bool {
	ln, err := net.Listen("tcp", ":"+fmt.Sprint(port))
	if err != nil {
		return true
	}
	defer ln.Close()
	return false
}
