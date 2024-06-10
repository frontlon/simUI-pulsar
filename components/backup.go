package components

import (
	"simUI/constant"
	"simUI/utils"
	"strings"
)

// 复制文件夹时，生成 src 和 dst 两个路径
func GetSrcAndDstPath(root, p, resType string) (string, string) {
	if root == "" || p == "" {
		return "", ""
	}

	if resType == "simulator" {
		p = utils.GetFilePath(p)
	}

	//绝对路径不拷贝
	src := strings.Replace(p, root, "", 1)
	if utils.IsAbsPath(src) {
		return p, ""
	}
	src = root + src
	dst := constant.ROOT_PATH + p

	//目录不存在，则新建
	fileType, _ := utils.CheckFileOrDir(dst)
	if fileType == 2 {
		utils.CreateDir(dst)

		//已存在，不复制
		if !utils.IsDirEmpty(dst) {
			return "", ""
		}
	} else {
		//已存在，不复制
		if utils.FileExists(dst) {
			return "", ""
		}
	}
	return src, dst
}
