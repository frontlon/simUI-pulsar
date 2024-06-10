package components

import (
	"errors"
	"os"
	"path/filepath"
	"simUI/utils"
	"strings"
)

// rom链接文件类型
var RomLinkExt = ".slnk"
var RomLinkSplit = "||"

type Slnk struct {
	RelRomPath string   //rom相对路径
	AbsRomPath string   //rom绝对路径
	Params     []string //启动参数
	RomName    string   //rom名称
}

// 读取所有rom链接
func ReadAllRomLinks(LinkPath string) (map[string]string, error) {
	result := map[string]string{}
	if err := filepath.Walk(LinkPath,
		func(p string, f os.FileInfo, err error) error {

			if f == nil {
				return nil
			}

			if f.IsDir() {
				return nil
			}
			if utils.GetFileExt(f.Name()) != RomLinkExt {
				return nil
			}
			romPath, _ := utils.ReadFile(p, false)
			result[utils.GetFileName(romPath)] = p
			return nil
		}); err != nil {
		return nil, err

	}
	return result, nil
}

// 读取一个rom的所有链接文件
func ReadRomLinksByRom(LinkPath, romName string) ([]string, error) {
	result := []string{}
	if err := filepath.Walk(LinkPath,
		func(p string, f os.FileInfo, err error) error {

			if f == nil {
				return nil
			}

			if f.IsDir() {
				return nil
			}
			if utils.GetFileExt(f.Name()) != RomLinkExt {
				return nil
			}
			romPath, _ := utils.ReadFile(p, false)

			fileRomName := utils.GetFileName(romPath)

			if fileRomName == romName {
				result = append(result, p)
			}

			return nil
		}); err != nil {
		return nil, err

	}
	return result, nil
}

// 读取rom链接文件
func GetLinkFileData(p string) *Slnk {

	if p == "" {
		return &Slnk{}
	}

	content, err := utils.ReadFile(p, false)

	arr := strings.Split(content, RomLinkSplit)
	romPath := arr[0]
	param := []string{}
	if len(arr) > 1 {
		param = strings.Split(arr[1], " ")
	}
	if err != nil {
		return &Slnk{}
	}

	data := &Slnk{
		RelRomPath: romPath,
		RomName:    utils.GetFileName(romPath),
		AbsRomPath: utils.ToAbsPath(romPath, ""),
		Params:     param,
	}
	return data
}

// 创建rom链接文件
func CreateLinkFile(p string, f, params string) error {
	if utils.FileExists(p) {
		return nil
	}
	content := f
	if params != "" {
		content = content + RomLinkSplit + params
	}

	err := utils.CreateFile(p, content)
	if err != nil {
		return nil
	}
	return nil
}

// 设置rom链接文件
func UpdateLinkFileData(p, romPath, param string) error {

	if p == "" {
		return nil
	}

	if !utils.FileExists(p) {
		return errors.New("链接文件不存在")
	}

	content := romPath + RomLinkSplit + param
	return utils.OverlayWriteFile(p, content)
}
