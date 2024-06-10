package modules

import (
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"simUI/config"
	"simUI/utils"
)

// 选择文件
func OpenFileDialog(typ string) (string, error) {

	opt := runtime.OpenDialogOptions{
		Filters: getDialogFilters(typ),
	}
	pth, _ := runtime.OpenFileDialog(config.Ctx, opt)
	return utils.ToRelPath(pth, ""), nil //转换为相对路径
}

// 选择多文件
func OpenMultiFileDialog(typ string) ([]string, error) {

	opt := runtime.OpenDialogOptions{
		Filters: getDialogFilters(typ),
	}
	paths, _ := runtime.OpenMultipleFilesDialog(config.Ctx, opt)

	if paths != nil && len(paths) > 0 {
		for k, v := range paths {
			paths[k] = utils.ToRelPath(v, "")
		}
	}

	return paths, nil //转换为相对路径
}

// 选择目录
func OpenDirectoryDialog() (string, error) {
	opt := runtime.OpenDialogOptions{}
	pth, _ := runtime.OpenDirectoryDialog(config.Ctx, opt)
	return utils.ToRelPath(pth, ""), nil //转换为相对路径
}

// 保存文件对话框
func SaveFileDialog(filename string) (string, error) {
	opt := runtime.SaveDialogOptions{
		CanCreateDirectories: true,     //允许用户创建目录
		DefaultFilename:      filename, //默认文件名
	}
	pth, _ := runtime.SaveFileDialog(config.Ctx, opt)
	return utils.ToRelPath(pth, ""), nil //转换为相对路径
}

// 编辑器图片选择框
func OpenFileDialogForEditor() (string, error) {
	filters := []runtime.FileFilter{}
	filters = append(filters, runtime.FileFilter{
		DisplayName: "Image Files (*.png, *.jpg, *.gif, *.ico, *.jpeg, *.webp, *.bmp)",
		Pattern:     "*.png;*.jpg;*.gif;*.ico;*.jpeg;*.webp;*.bmp",
	})
	filters = append(filters, runtime.FileFilter{
		DisplayName: "All Files (*.*)",
		Pattern:     "*.*",
	})

	opt := runtime.OpenDialogOptions{
		Filters: filters,
	}
	pth, _ := runtime.OpenFileDialog(config.Ctx, opt)

	pth = utils.WailsPathEncode(pth, true)

	return pth, nil //转换为相对路径
}

func getDialogFilters(typ string) []runtime.FileFilter {
	filters := []runtime.FileFilter{}

	switch typ {
	case "image":
		filters = append(filters, runtime.FileFilter{
			DisplayName: "Image Files (*.png, *.jpg, *.gif, *.ico, *.jpeg, *.webp, *.bmp)",
			Pattern:     "*.png;*.jpg;*.gif;*.ico;*.jpeg;*.webp;*.bmp",
		})
	case "media":
		filters = append(filters, runtime.FileFilter{
			DisplayName: "Media Files (*.png, *.jpg, *.gif, *.ico, *.jpeg, *.webp, *.bmp, *.wmv, *.mp4, *.avi, *.flv, *.webm)",
			Pattern:     "*.png;*.jpg;*.gif;*.ico;*.jpeg;*.webp;*.bmp;*.wmv;*.mp4;*.avi;*.flv;*.webm",
		})
	case "audio":
		filters = append(filters, runtime.FileFilter{
			DisplayName: "Audio Files (*.mp3, *.dmi, *.wav, *.wma)",
			Pattern:     "*.mp3;*.dmi;*.wav;*.wma",
		})
	case "video":
		filters = append(filters, runtime.FileFilter{
			DisplayName: "Video Files (*.wmv, *.mp4, *.avi, *.flv, *.webm)",
			Pattern:     "*.wmv;*.mp4;*.avi;*.flv;*.webm",
		})
	case "doc":
		filters = append(filters, runtime.FileFilter{
			DisplayName: "Doc Files (*.txt, *.html, *.htm, *.md)",
			Pattern:     "*.txt;*.html;*.htm;*.md",
		})
	case "app":
		filters = append(filters, runtime.FileFilter{
			DisplayName: "Application Files (*.exe, *.lnk, *.bat, *.cmd, *.com)",
			Pattern:     "*.exe;*.lnk;*.bat;*.cmd;*.com",
		})
	case "share":
		filters = append(filters, runtime.FileFilter{
			DisplayName: "Share File (*.shr)",
			Pattern:     "*.shr",
		})
	}

	filters = append(filters, runtime.FileFilter{
		DisplayName: "All Files (*.*)",
		Pattern:     "*.*",
	})

	return filters
}
