package modules

import (
	"simUI/server"
)

// 启动上传服务
func StartUploadServer() string {
	if server.Addr == "" {
		server.StartHttpServer()
	}
	return server.Addr
}
