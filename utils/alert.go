package utils

import "github.com/ncruces/zenity"

// 调用系统dialog框 - 错误
func DialogError(title, msg string) {
	zenity.Error(msg, zenity.Title(title), zenity.ErrorIcon)
}

// 调用系统dialog框 - 信息
func DialogInfo(title, msg string) {
	zenity.Info(msg, zenity.Title(title), zenity.InfoIcon)

}

// 调用系统dialog框 - 警告
func DialogWarn(title, msg string) {
	zenity.Warning(msg, zenity.Title(title), zenity.WarningIcon)
}
