package components

import (
	"github.com/adrg/sysfont"
	"simUI/constant"
	"simUI/db"
	"simUI/utils"
	"strings"
)

// 读取用户字体
func GetUserFontList() []db.PlatformUIFont {

	files, err := utils.ScanCurrentDir(constant.FONT_PATH)

	list := []db.PlatformUIFont{}

	if err != nil || len(files) == 0 {
		return list
	}

	for _, f := range files {

		if f.IsDir() {
			continue
		}

		ext := strings.ToLower(utils.GetFileExt(f.Name()))
		if _, ok := constant.FONT_EXTS[ext]; !ok {
			continue
		}

		fonts := db.PlatformUIFont{
			Type:   2,
			Family: utils.GetFileName(f.Name()),
			Format: constant.FONT_EXTS[ext],
			Src:    utils.ToRelPath(constant.FONT_PATH+f.Name(), ""),
		}
		list = append(list, fonts)
	}

	return list
}

/*
读取系统字体列表
*/
func GetSystemFontList() []db.PlatformUIFont {

	// 创建一个字体查找器
	finder := sysfont.NewFinder(nil)

	// 获取系统所有字体
	fonts := finder.List()

	if fonts == nil || len(fonts) == 0 {
		return []db.PlatformUIFont{}
	}

	// 打印字体名称
	list := []string{}
	for _, font := range fonts {
		if font.Family == "" {
			continue
		}
		list = append(list, font.Family)
	}

	list = utils.SliceRemoveDuplicate(list)

	resp := []db.PlatformUIFont{}
	for _, family := range list {
		f := db.PlatformUIFont{
			Type:   1,
			Family: family,
		}
		resp = append(resp, f)

	}
	return resp
}
