package controller

import "simUI/modules"

// 启动上传服务
func (a *Controller) StartUploadServer() string {
	return Resp(modules.StartUploadServer(), nil)
}
