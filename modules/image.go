package modules

import (
	"errors"
	"simUI/components"
	"simUI/db"
	"simUI/utils"
	"strings"
	"time"
)

var imageWidth = 320
var quality = 80

func CreateRomResByBase64(id uint64, resType string, slaveRes uint8, fileType string, base64Str string) (string, error) {

	fileTypes := strings.Split(fileType, "/")
	if fileTypes[0] != "image" && fileTypes[0] != "video" {
		return "", errors.New("文件类型错误")
	}

	rom, err := (&db.Rom{}).GetById(id)
	if err != nil {
		return "", err
	}

	resPath := components.GetResPathByType(resType, rom.Platform)
	pth := ""
	tm := utils.ToString(time.Now().Unix())
	if slaveRes == 0 {
		pth = resPath + "/" + rom.RomName + "." + fileTypes[1]
	} else {
		pth = resPath + "/" + rom.RomName + "/" + tm + "." + fileTypes[1]
	}

	if err = utils.Base64ToFile(base64Str, pth); err != nil {
		return "", err
	}

	//return utils.WailsPathEncode(pth), nil
	return pth, nil

}

func CreateImg() {
	/** 判断是不是图片 */
	/*in := "/Users/frontlon/go/src/wails/sim-ui-pulsar/in.jpg"
	format := utils.GetFileExt(in)
	out := "/Users/frontlon/go/src/wails/sim-ui-pulsar/out" + format

	err := utils.CreateThumbnail(in, out)

	fmt.Println("-=-=-")
	fmt.Println(err)

	if err != nil {
		fmt.Println(err)
	}
	*/
}

//func CreateOptimizedImage(platform uint32, opt string) error {
//
//	paths := components.GetResPath(platform, 0)
//	path := paths[opt]
//
//	platformInfo := (&db.Platform{}).GetVOById(platform, false)
//	outputPath := platformInfo.OptimizedPath
//
//	//先删除原文件
//	utils.DeleteDir(outputPath)
//	utils.CreateDir(outputPath)
//
//	//读取文件总数
//	//files, _ := ioutil.ReadDir(path)
//	//fileCount := len(files)
//	i := 0
//	filepath.Walk(path, func(p string, f os.FileInfo, err error) error {
//		if f == nil {
//			return err
//		}
//		if f.IsDir() { /** 是否是目录 */
//			return nil
//		}
//
//		/** 判断是不是图片 */
//		format := utils.GetFileExt(p)
//
//		outputPath := outputPath + "/" + utils.GetFileNameAndExt(p)
//		if p != "" {
//			err := utils.ImageCompress(
//				func() (io.Reader, error) {
//					return os.Open(p)
//				},
//				func() (*os.File, error) {
//					return os.Open(p)
//				},
//				outputPath,
//				quality,
//				imageWidth,
//				format)
//
//			if err != nil {
//				return err
//			}
//		}
//
//		/*if i%10 == 0 {
//			utils.Loading("[2/3]已生成("+utils.ToString(i)+" / "+utils.ToString(fileCount)+")", config.Cfg.Platform[platform].Name)
//		}*/
//
//		i++
//		return nil
//	})
//
//	//数据更新完成后，页面回调，更新页面DOM
//	/*if _, err := utils.Window.Call("CB_createOptimizedCache"); err != nil {
//		fmt.Print(err)
//	}*/
//
//	return nil
//}
