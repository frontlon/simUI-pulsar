package controller

import (
	"simUI/modules"
)

// 检测升级
func (a *Controller) CheckUpgrade(newVersion string) string {
	return Resp(modules.CheckUpgrade(newVersion))
}

// 跳过更新
func (a *Controller) JumpUpgrade(version string) string {
	return Resp("", modules.JumpUpgrade(version))
}

// 下载新版本
func (a *Controller) DownloadNewVersion(url string) string {
	return Resp(modules.DownloadNewVersion(url))
}

// 安装新版本
func (a *Controller) InstallUpgrade(version string, unzipPath string) string {
	return Resp("", modules.InstallUpgrade(version, unzipPath))
}
