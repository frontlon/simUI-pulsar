package controller

import (
	"simUI/modules"
)

// 弹出文件选择窗口
func (a *Controller) OpenFileDialog(typ string) string {
	return Resp(modules.OpenFileDialog(typ))
}

// 弹出多文件选择窗口
func (a *Controller) OpenMultiFileDialog(typ string) string {
	return Resp(modules.OpenMultiFileDialog(typ))
}

// 弹出目录选择窗口
func (a *Controller) OpenDirectoryDialog() string {
	return Resp(modules.OpenDirectoryDialog())
}

// 保存文件对话框
func (a *Controller) SaveFileDialog(filename string) string {
	return Resp(modules.SaveFileDialog(filename))
}

// 编辑器图片选择框
func (a *Controller) OpenFileDialogForEditor() string {
	return Resp(modules.OpenFileDialogForEditor())
}
