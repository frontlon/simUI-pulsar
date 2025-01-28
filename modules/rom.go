package modules

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"simUI/components"
	"simUI/config"
	"simUI/constant"
	"simUI/db"
	"simUI/request"
	"simUI/utils"
	"sort"
	"strings"
	"time"
)

type RomDetail struct {
	Info       *db.Rom        //rom信息
	DocContent string         //简介内容
	PicList    []string       //展示图片列表
	Sublist    []*db.Rom      //子游戏
	Extra      RomDetailExtra //附近信息
}

type RomDetailExtra struct {
	RomPath      string //rom真实路径
	PlatformName string //平台名称
}

type GameThumb struct {
	Type   string //图片类别
	Path   string //图片路径
	Master int    //是否为主图
	Index  int    //索引号
}

type GameAudio struct {
	Name string //音频名称
	Path string //音频路径
}

// 读取rom游戏列表
func GetGameList(data *request.GetGameList) (any, error) {
	romList := []*db.Rom{}
	m := &db.Rom{}
	if data.Letter == "" {
		romList, _ = m.Get(data.ShowHide, data.Page, data.Platform, data.Catname, data.CatnameLike, data.Keyword, data.BaseType, data.BasePublisher, data.BaseYear, data.BaseCountry, data.BaseTranslate, data.BaseVersion, data.BaseProducer, data.Score, data.Complete)
	} else {
		//按拼音查询
		romList, _ = m.GetByPinyin(data.ShowHide, data.Page, data.Platform, data.Catname, data.CatnameLike, data.Letter)
	}

	if len(romList) == 0 {
		return []string{}, nil
	}

	platformMap := (&db.Platform{}).GetPlatformVOMap(false)

	//定义展示图名称
	romNameMap := map[uint32]map[string]string{}
	for _, v := range romList {
		if _, ok := romNameMap[v.Platform]; ok {
			romNameMap[v.Platform][v.RomName] = "1"
		} else {
			romNameMap[v.Platform] = map[string]string{
				v.RomName: "1",
			}
		}
	}

	switch data.SimpleMode {
	case "simple":
		//格式转换
		datas := (&db.Rom{}).ConvertRomListSimple(romList, data.ShowSubGame)

		//读取展示图
		thumbsMap := map[uint32]map[string]string{}
		if data.Platform == 0 {
			for pfId, _ := range platformMap {
				blockThumbType := "thumb"
				if _, ok := platformMap[pfId]; ok {
					if data.Theme == constant.UI_PLAYNITE {
						blockThumbType = platformMap[pfId].Ui.Playnite.BlockThumbType
					} else if data.Theme == constant.UI_TINY {
						blockThumbType = platformMap[pfId].Ui.Tiny.BlockThumbType
					} else {
						blockThumbType = platformMap[pfId].Ui.Default.BlockThumbType
					}
				}
				thumbsMap[pfId] = components.GetMasterRomThumbs(blockThumbType, pfId, romNameMap[pfId])
			}
		} else {
			blockThumbType := "thumb"
			if _, ok := platformMap[data.Platform]; ok {
				if data.Theme == constant.UI_PLAYNITE {
					blockThumbType = platformMap[data.Platform].Ui.Playnite.BlockThumbType
				} else if data.Theme == constant.UI_TINY {
					blockThumbType = platformMap[data.Platform].Ui.Tiny.BlockThumbType
				} else {
					blockThumbType = platformMap[data.Platform].Ui.Default.BlockThumbType
				}
			}
			thumbsMap[data.Platform] = components.GetMasterRomThumbs(blockThumbType, data.Platform, romNameMap[data.Platform])
		}
		for k, v := range datas {
			if _, ok := thumbsMap[v.Platform]; !ok {
				continue
			}
			if _, ok := thumbsMap[v.Platform][v.RomName]; !ok {
				continue
			}
			datas[k].ThumbPic = utils.WailsPathEncode(thumbsMap[v.Platform][v.RomName], true)
		}
		return datas, nil
	case "simplest":
		//格式转换
		thumbsMap := map[uint32]map[string]string{}
		datas := (&db.Rom{}).ConvertRomListSimplest(romList)

		blockThumbType := "thumb"
		if _, ok := platformMap[data.Platform]; ok {
			if data.Theme == constant.UI_PLAYNITE {
				blockThumbType = platformMap[data.Platform].Ui.Playnite.BlockThumbType
			} else if data.Theme == constant.UI_TINY {
				blockThumbType = platformMap[data.Platform].Ui.Tiny.BlockThumbType
			} else {
				blockThumbType = platformMap[data.Platform].Ui.Default.BlockThumbType
			}
		}
		thumbsMap[data.Platform] = components.GetMasterRomThumbs(blockThumbType, data.Platform, romNameMap[data.Platform])

		for k, v := range datas {
			if _, ok := thumbsMap[v.Platform]; !ok {
				continue
			}
			if _, ok := thumbsMap[v.Platform][v.RomName]; !ok {
				continue
			}

			datas[k].ThumbPic = utils.WailsPathEncode(thumbsMap[v.Platform][v.RomName], true)

		}

		return datas, nil
	}
	return nil, errors.New("SimpleMode is error")

}

// 读取rom游戏数量
func GetGameCount(data *request.GetGameList) (int, error) {
	count := 0
	if data.Letter == "" {
		count, _ = (&db.Rom{}).Count(data.ShowHide, data.Platform, data.Catname, data.CatnameLike, data.Keyword, data.BaseType, data.BasePublisher, data.BaseYear, data.BaseCountry, data.BaseTranslate, data.BaseVersion, data.BaseProducer, data.Score, data.Complete)
	} else {
		//按拼音查询
		count, _ = (&db.Rom{}).CountByPinyin(data.ShowHide, data.Page, data.Platform, data.Catname, data.CatnameLike, data.Letter)
	}
	return count, nil
}

// 根据菜单读取游戏数
func GetGameCountByMenu(platform uint32) (map[string]int64, error) {
	return (&db.Rom{}).CountGroupMenu(platform)
}

/**
 * 设置通关状态
 **/
func UpdateGameComplete(id uint64, status uint8) error {

	//更新csv文件
	rom, _ := (&db.Rom{}).GetById(id)

	//更新数据
	(&db.Rom{}).UpdateOneField(id, "complete", status)

	//更新csv
	go func() {
		config.SetRomSettingOneField(rom.Platform, rom.RomName, "Complete", status, true)
	}()

	return nil
}

/**
 * 设置游戏评分
 **/
func UpdateGameScore(id uint64, score float64) error {

	//更新csv文件
	rom, _ := (&db.Rom{}).GetById(id)

	platformInfo := (&db.Platform{}).GetVOById(rom.Platform, false)

	if platformInfo.RomBaseFile == "" {
		return errors.New(config.Cfg.Lang["romBaseFileNotFound"])
	}

	rombase := config.GetRomBaseById(rom.Platform, rom.RomName)

	create := &config.RomBase{}

	if rombase != nil {
		create = rombase
		create.Score = utils.ToString(score)
	} else {
		create = &config.RomBase{
			RomName: rom.Name,
			Score:   utils.ToString(score),
		}
	}
	if err := config.AddRomBase(rom.Platform, create); err != nil {
		return err
	}

	//更新数据库
	if err := (&db.Rom{}).UpdateOneField(id, "score", score); err != nil {
		return err
	}

	return nil
}

// 读取rom游戏列表
func GetSubGameList(id uint64) ([]*db.Rom, error) {
	//游戏游戏详细数据
	info, err := (&db.Rom{}).GetById(id)

	if err != nil {
		return nil, err
	}

	//子游戏列表
	sub, _ := (&db.Rom{}).GetSubRom(info.Platform, info.RomName)
	return sub, nil
}

// 右键打开文件夹
func OpenFolder(opt string, id uint64) error {

	info := &db.Rom{}
	platform := &db.PlatformVO{}
	simulator := &db.Simulator{}
	if opt == "sim" {
		//打开模拟器读模拟器数据
		simu, err := (&db.Simulator{}).GetById(uint32(id))
		if err != nil {
			return err
		}
		simulator = simu
	} else {
		//非模拟器读rom数据
		romInfo, err := (&db.Rom{}).GetById(id)
		if err != nil {
			return err
		}
		info = romInfo
		platform = (&db.Platform{}).GetVOById(info.Platform, false)
	}

	fileName := ""
	switch opt {
	case "rom": //rom文件
		slnk := components.GetLinkFileData(utils.ToAbsPath(info.LinkFile, platform.LinkPath))
		fileName = slnk.AbsRomPath
	case "slnk": //链接文件
		fileName = utils.ToAbsPath(info.LinkFile, platform.LinkPath)
	case "platform": //平台根目录
		fileName = platform.RootPath
	case "sim": //模拟器路径
		fileName = utils.ToAbsPath(simulator.Path, "")
	default: //资源文件
		res := components.GetResPath(platform.Id, 0)
		if res[opt] != "" {
			romName := utils.GetFileName(filepath.Base(info.LinkFile)) //读取文件名
			fileName = components.GetMasterRes(opt, info.Platform, romName)
			if fileName == "" {
				fileName = res[opt]
			}
		}
	}
	fmt.Println("打开文件夹", opt, id, fileName)

	if fileName != "" {
		if err := components.OpenFolderByWindow(fileName); err != nil {
			return err
		}
	}
	return nil
}

// 直接定位文件
func OpenFolderByPath(pth string) error {

	if pth == "" {
		return errors.New("目录不能为空")
	}

	if pth != "" {
		if err := components.OpenFolderByWindow(pth); err != nil {
			return err
		}
	}
	return nil
}

// 添加ROM
func AddFileGame(platform uint32, menu, romPath string, romFiles []string) error {

	platformInfo := (&db.Platform{}).GetVOById(platform, false)

	root := romPath + "/"
	folder := ""
	if menu != "" {
		folder = root + menu + "/"
	}

	//如果目录不存在，则放到根目录中
	if !utils.DirExists(folder) {
		folder = root
	}

	for _, romFile := range romFiles {
		romFile = utils.ToAbsPath(romFile, "")
		if !utils.FileExists(romFile) {
			return errors.New(romFile + "不存在")
		}

		//复制rom到simui
		dst := romPath + "/" + utils.GetFileNameAndExt(romFile)
		if err := utils.FileCopy(romFile, dst); err != nil {
			return errors.New(romFile + "已存在")
		}

		filename := utils.GetFileName(romFile)
		p := platformInfo.LinkPath + "/" + filename + components.RomLinkExt //slnk文件路径

		//如果文件已存在，则使用新名称
		if utils.FileExists(p) {
			t := utils.ToString(time.Now().Unix())
			p = folder + filename + t + components.RomLinkExt //slnk文件路径
		}

		//转换为相对路径
		rel := utils.ToRelPath(dst, "")
		if err := components.CreateLinkFile(p, rel, ""); err != nil {
			return err
		}
	}

	return nil
}

// 添加PC或文件夹ROM
func AddPcOrFolderGame(platform uint32, menu, romPath, romFile string, bootParam string, isBatFile uint8) error {

	romFile = utils.ToAbsPath(romFile, "")
	if !utils.IsExist(romFile) {
		return errors.New(romFile + "不存在")
	}
	platformInfo := (&db.Platform{}).GetVOById(platform, false)

	linkRomFile := ""
	if isBatFile == 1 {
		//在rom目录创建bat文件
		batName := utils.GetFileName(romFile) + ".bat"
		linkRomFile = utils.ToAbsPath(batName, romPath)

		//如果文件存在，改个名
		if utils.FileExists(linkRomFile) {
			batName = utils.GetFileName(romFile) + "_" + utils.ToString(time.Now().Unix()) + ".bat"
			linkRomFile = utils.ToAbsPath(batName, romPath)
		}
		batContent := romFile + " " + bootParam
		if err := components.CreateGameBat(linkRomFile, batContent); err != nil {
			return err
		}
	} else {
		linkRomFile = romFile
	}
	linkFileName := utils.GetFileName(romFile) + components.RomLinkExt
	linkFolder := platformInfo.LinkPath + "/"
	if menu != "" {
		linkFolder = platformInfo.LinkPath + menu
	}

	linkFile := linkFolder + linkFileName
	if utils.DirExists(linkFile) {
		//如果文件已存在，则使用新名称
		t := utils.ToString(time.Now().Unix())
		linkFileName = utils.GetFileName(romFile) + t + components.RomLinkExt
		linkFile = linkFolder + linkFileName
	}
	rel := utils.ToRelPath(linkRomFile, "")

	if err := components.CreateLinkFile(linkFile, rel, bootParam); err != nil {
		return err
	}

	return nil
}

// 读取rom详情
func GetGameDetail(id uint64) (*RomDetail, error) {

	res := &RomDetail{}
	//游戏游戏详细数据
	info, err := (&db.Rom{}).GetById(id)

	if err != nil {
		return res, err
	}
	platformInfo := (&db.Platform{}).GetVOById(info.Platform, false)

	//子游戏列表
	sub, _ := (&db.Rom{}).GetSubRom(info.Platform, info.RomName)

	//附加信息
	linkp := utils.ToAbsPath(info.LinkFile, platformInfo.LinkPath)
	linkData := components.GetLinkFileData(linkp)
	extra := RomDetailExtra{
		RomPath:      linkData.RelRomPath,
		PlatformName: platformInfo.Name,
	}

	res.Info = info
	res.Sublist = sub
	res.Extra = extra

	//读取文档内容
	if platformInfo.DocPath != "" {
		docFileName := ""
		for _, v := range constant.DOC_EXTS {
			docFileName = platformInfo.DocPath + "/" + info.RomName + v
			res.DocContent = GetDocContent(docFileName)
			if res.DocContent != "" {
				res.DocContent = strings.Trim(res.DocContent, "\t")
				res.DocContent = strings.Trim(res.DocContent, "\n\r")
				res.DocContent = strings.Trim(res.DocContent, "\r")
				res.DocContent = strings.Trim(res.DocContent, "\n")
				res.DocContent = strings.ReplaceAll(res.DocContent, "\r\n", "<br/>")
				res.DocContent = strings.ReplaceAll(res.DocContent, "\r", "<br/>")
				res.DocContent = strings.ReplaceAll(res.DocContent, "\n", "<br/>")
				break
			}
		}
	}

	return res, nil
}

// 读取rom展示图列表
func GetGameThumbs(id uint64, imgType string) ([]GameThumb, error) {

	//游戏游戏详细数据
	info, err := (&db.Rom{}).GetById(id)
	pics := []GameThumb{}
	if err != nil {
		return pics, err
	}

	//图片列表
	resPaths := map[string]string{}
	if imgType == "" {
		resPaths = components.GetResPath(info.Platform, 2)
	} else {
		resPaths[imgType] = components.GetResPathByType(imgType, info.Platform)
	}
	idx := 0
	for typ, pth := range resPaths {
		//读主资源
		r := components.GetMasterRes(typ, info.Platform, info.RomName)
		if r != "" {
			s := GameThumb{
				Type:   typ,
				Path:   utils.WailsPathEncode(r, true),
				Master: 1,
				Index:  idx,
			}
			pics = append(pics, s)
			idx++
		}

		//读子资源
		slaves, _ := components.GetSlaveRes(pth, strings.TrimSpace(info.RomName))

		//子资源文件名排序
		sort.Slice(slaves, func(i, j int) bool {
			return strings.ToLower(slaves[i]) < strings.ToLower(slaves[j])
		})

		if len(slaves) > 0 {
			for _, v := range slaves {
				s := GameThumb{
					Type:   typ,
					Path:   utils.WailsPathEncode(v, true),
					Master: 0,
					Index:  idx,
				}
				pics = append(pics, s)
				idx++
			}
		}
	}
	return pics, nil
}

// 读取游戏攻略内容
func GetGameDoc(t string, id uint64, toHtml uint8) (string, error) {

	//游戏游戏详细数据
	info, err := (&db.Rom{}).GetById(id)

	platformInfo := (&db.Platform{}).GetVOById(info.Platform, false)

	if err != nil {
		return "", err
	}

	//如果没有执行运行的文件，则读取文档内容
	strategy := ""
	for _, v := range constant.DOC_EXTS {
		strategyFileName := ""
		if t == "strategy" {
			strategyFileName = platformInfo.StrategyPath + "/" + info.RomName + v
		} else if t == "doc" {
			strategyFileName = platformInfo.DocPath + "/" + info.RomName + v
		}
		strategy = GetDocContent(strategyFileName)
		if strategy != "" {
			break
		}
	}

	//攻略中资源相对路径都改成绝对路径
	if strategy != "" {
		re := regexp.MustCompile(utils.WailsPathPattern)
		strategy = re.ReplaceAllStringFunc(strategy, func(match string) string {
			pth := utils.WailsPathDecode(match)
			pth = utils.ToAbsPath(pth, "")
			return utils.WailsPathEncode(pth, true)
		})
	}

	//攻略替换操作按键
	if t == "strategy" {
		strategy = strings.ReplaceAll(strategy, "\r\n", "<br/>")
		strategy = strings.ReplaceAll(strategy, "\r", "<br/>")
		strategy = strings.ReplaceAll(strategy, "\n", "<br/>")
		strategy = strings.ReplaceAll(strategy, "_+", "<img class='ctl' src='/images/ctl/+.png'/>")
		strategy = strings.ReplaceAll(strategy, "_1", "<img class='ctl' src='/images/ctl/1.png'/>")
		strategy = strings.ReplaceAll(strategy, "_2", "<img class='ctl' src='/images/ctl/2.png'/>")
		strategy = strings.ReplaceAll(strategy, "_3", "<img class='ctl' src='/images/ctl/3.png'/>")
		strategy = strings.ReplaceAll(strategy, "_4", "<img class='ctl' src='/images/ctl/4.png'/>")
		strategy = strings.ReplaceAll(strategy, "_5", "<img class='ctl' src='/images/ctl/5.png'/>")
		strategy = strings.ReplaceAll(strategy, "_6", "<img class='ctl' src='/images/ctl/6.png'/>")
		strategy = strings.ReplaceAll(strategy, "_7", "<img class='ctl' src='/images/ctl/7.png'/>")
		strategy = strings.ReplaceAll(strategy, "_8", "<img class='ctl' src='/images/ctl/8.png'/>")
		strategy = strings.ReplaceAll(strategy, "_9", "<img class='ctl' src='/images/ctl/9.png'/>")
		strategy = strings.ReplaceAll(strategy, "_A", "<img class='ctl' src='/images/ctl/A.png'/>")
		strategy = strings.ReplaceAll(strategy, "_B", "<img class='ctl' src='/images/ctl/B.png'/>")
		strategy = strings.ReplaceAll(strategy, "_C", "<img class='ctl' src='/images/ctl/C.png'/>")
		strategy = strings.ReplaceAll(strategy, "_D", "<img class='ctl' src='/images/ctl/D.png'/>")
		strategy = strings.ReplaceAll(strategy, "_E", "<img class='ctl' src='/images/ctl/E.png'/>")
		strategy = strings.ReplaceAll(strategy, "_F", "<img class='ctl' src='/images/ctl/F.png'/>")
		strategy = strings.ReplaceAll(strategy, "_N", "<img class='ctl' src='/images/ctl/N.png'/>")
		strategy = strings.ReplaceAll(strategy, "_S", "<img class='ctl' src='/images/ctl/S.png'/>")
		strategy = strings.ReplaceAll(strategy, "_P", "<img class='ctl' src='/images/ctl/P.png'/>")
		strategy = strings.ReplaceAll(strategy, "_K", "<img class='ctl' src='/images/ctl/K.png'/>")

	}

	return strategy, nil
}

// 更新游戏攻略内容
func SetGameDoc(t string, id uint64, content string) error {

	//游戏游戏详细数据
	info, err := (&db.Rom{}).GetById(id)

	if err != nil {
		return err
	}

	platformInfo := (&db.Platform{}).GetVOById(info.Platform, false)

	//如果没有执行运行的文件，则读取文档内容
	newExt := ""
	Filename := ""
	for _, v := range constant.DOC_EXTS {
		strategyFileName := ""
		if t == "strategy" {
			strategyFileName = platformInfo.StrategyPath + "/" + info.RomName + v
			newExt = platformInfo.StrategyPath + "/" + info.RomName + ".txt"
		} else if t == "doc" {
			strategyFileName = platformInfo.DocPath + "/" + info.RomName + v
			newExt = platformInfo.DocPath + "/" + info.RomName + ".txt"
		}

		if utils.FileExists(strategyFileName) {
			Filename = strategyFileName
			break
		}
	}

	if Filename == "" {
		Filename = newExt
	}

	if !utils.FileExists(Filename) {
		if err := utils.CreateFile(Filename, ""); err != nil {
			return err
		}
	}

	if !utils.IsUTF8(content) {
		content = utils.ToUTF8(content)
	}

	content = strings.TrimSpace(content)

	//替换图片路径为相对路径
	re := regexp.MustCompile(utils.WailsPathPattern)
	content = re.ReplaceAllStringFunc(content, func(match string) string {
		src := utils.WailsPathDecode(match)
		ext := utils.GetFileExt(src)
		//检查图片是否在软件目录，不在的话复制进来
		if !strings.Contains(src, constant.ROOT_PATH) {
			if !utils.FileExists(src) {
				dst := platformInfo.UploadPath + "/" + utils.CreateUUid() + ext
				utils.FileCopy(src, dst)
				src = dst
			}
		}
		//资源转换成相对路径
		rel := utils.ToRelPath(src, "")
		return utils.WailsPathEncode(rel, true)
	})

	content = strings.ReplaceAll(content, constant.ROOT_PATH, "")
	if err := utils.OverlayWriteFile(Filename, content); err != nil {
		return err
	}

	return nil
}

// 删除游戏攻略内容
func DelGameDoc(t string, id uint64) error {

	//游戏游戏详细数据
	info, err := (&db.Rom{}).GetById(id)

	if err != nil {
		return err
	}

	//如果没有执行运行的文件，则读取文档内容
	res := components.GetResPath(info.Platform, 0)
	Filename := ""
	for _, v := range constant.DOC_EXTS {
		strategyFileName := res[t] + "/" + info.RomName + v

		if utils.FileExists(strategyFileName) {
			Filename = strategyFileName
			break
		}
	}

	if Filename != "" && utils.FileExists(Filename) {
		if err := utils.FileDelete(Filename); err != nil {
			return err
		}
	}

	return nil
}

/**
 * 读取游戏介绍文本
 **/
func GetDocContent(f string) string {
	if f == "" {
		return ""
	}
	text, err := ioutil.ReadFile(f)
	content := ""
	if err != nil {
		return content
	}
	content = string(text)

	if !utils.IsUTF8(content) {
		content = utils.ToUTF8(content)
	}

	return content
}

// 根据romid读取资料
func GetRomBase(id uint64) (*config.RomBase, error) {
	rom, _ := (&db.Rom{}).GetById(id)
	return config.GetRomBaseById(rom.Platform, rom.RomName), nil
}

// 编辑rom基础信息
func SetRomBase(id uint64, d map[string]string, alias string) (string, error) {

	rom, _ := (&db.Rom{}).GetById(id)
	if alias == rom.RomName {
		alias = ""
	}

	platformInfo := (&db.Platform{}).GetVOById(rom.Platform, false)

	if platformInfo.RomBaseFile == "" {
		return "", errors.New(config.Cfg.Lang["romBaseFileNotFound"])
	}

	romBase := &config.RomBase{
		RomName:   rom.RomName,
		Name:      alias,
		NameEN:    d["NameEN"],
		NameJP:    d["NameJP"],
		Type:      d["Type"],
		Year:      d["Year"],
		Producer:  d["Producer"],
		Publisher: d["Publisher"],
		Country:   d["Country"],
		Translate: d["Translate"],
		Version:   d["Version"],
		OtherA:    d["OtherA"],
		OtherB:    d["OtherB"],
		OtherC:    d["OtherC"],
		OtherD:    d["OtherD"],
		Score:     d["Score"],
	}

	//写入配置文件
	if err := config.AddRomBase(rom.Platform, romBase); err != nil {
		return "", err
	}

	//更新到数据库

	name := alias
	if name == "" {
		name = rom.Name
	}

	infoMd5 := (&db.Rom{
		Name:          name,
		LinkFile:      rom.LinkFile,
		Score:         utils.ToFloat64(d["Score"]),
		Size:          rom.Size,
		BaseNameEn:    d["NameEN"],
		BaseNameJp:    d["NameJP"],
		BaseType:      d["Type"],
		BaseYear:      d["Year"],
		BaseProducer:  d["Producer"],
		BasePublisher: d["Publisher"],
		BaseCountry:   d["Country"],
		BaseTranslate: d["Translate"],
		BaseVersion:   d["Version"],
		BaseOtherA:    d["OtherA"],
		BaseOtherB:    d["OtherB"],
		BaseOtherC:    d["OtherC"],
		BaseOtherD:    d["OtherD"],
	}).CreateInfoMd5()

	dbRom := &db.Rom{
		RomName:       rom.RomName,
		Name:          name,
		BaseType:      d["Type"],
		BaseYear:      d["Year"],
		BaseProducer:  d["Producer"],
		BasePublisher: d["Publisher"],
		BaseCountry:   d["Country"],
		BaseTranslate: d["Translate"],
		BaseVersion:   d["Version"],
		Score:         utils.ToFloat64(d["Score"]),
		BaseNameEn:    d["NameEN"],
		BaseNameJp:    d["NameJP"],
		BaseOtherA:    d["OtherA"],
		BaseOtherB:    d["OtherB"],
		BaseOtherC:    d["OtherC"],
		BaseOtherD:    d["OtherD"],
		InfoMd5:       infoMd5,
	}
	if err := dbRom.UpdateRomBase(); err != nil {
		return "", err
	}

	return "", nil
}

// 批量编辑rom基础信息
func BatchSetRomBase(platform uint32, data map[string]*db.RomSimpleVO) error {

	//编辑资料
	rombaseList, _ := config.GetRomBase(platform, false)

	for _, d := range data {
		rombaseList[d.RomName] = &config.RomBase{
			RomName:   d.RomName,
			Name:      d.Name,
			NameEN:    d.BaseNameEn,
			NameJP:    d.BaseNameJp,
			Type:      d.BaseType,
			Year:      d.BaseYear,
			Producer:  d.BaseProducer,
			Publisher: d.BasePublisher,
			Country:   d.BaseCountry,
			Translate: d.BaseTranslate,
			Version:   d.BaseVersion,
			OtherA:    d.BaseOtherA,
			OtherB:    d.BaseOtherB,
			OtherC:    d.BaseOtherC,
			OtherD:    d.BaseOtherD,
			Score:     utils.ToString(d.Score),
		}
	}

	//写入配置文件
	if err := config.CoverRomBaseFile(platform, rombaseList); err != nil {
		return err
	}

	//更新数据库
	for _, v := range data {

		m := &db.Rom{
			RomName:       v.RomName,
			Name:          v.Name,
			BaseNameEn:    v.BaseNameEn,
			BaseNameJp:    v.BaseNameJp,
			BaseType:      v.BaseType,
			BaseYear:      v.BaseYear,
			BaseProducer:  v.BaseProducer,
			BasePublisher: v.BasePublisher,
			BaseCountry:   v.BaseCountry,
			BaseTranslate: v.BaseTranslate,
			BaseVersion:   v.BaseVersion,
			BaseOtherA:    v.BaseOtherA,
			BaseOtherB:    v.BaseOtherB,
			BaseOtherC:    v.BaseOtherC,
			BaseOtherD:    v.BaseOtherD,
		}
		m.InfoMd5 = m.CreateInfoMd5()
		err := m.UpdateRomBase()
		if err != nil {
			return err
		}
	}

	return nil
}

/**
 * 读取过滤器列表
 **/
func GetFilter(platform uint32) (map[string][]map[string]any, error) {

	volist := []*db.Filter{}
	if platform == 0 {
		volist, _ = (&db.Filter{}).GetAll()
	} else {
		volist, _ = (&db.Filter{}).GetByPlatform(platform)
	}

	//填充数据
	filterList := map[string][]map[string]any{}
	for _, v := range volist {

		val := map[string]any{}
		if v.Type == "complete" { //通关状态
			name := ""
			if v.Name == "0" {
				name = config.Cfg.Lang["notCleared"]
			} else if v.Name == "1" {
				name = config.Cfg.Lang["passed"]
			} else if v.Name == "2" {
				name = config.Cfg.Lang["perfectClear"]
			}
			val["label"] = name
			val["value"] = v.Name
		} else {
			val["label"] = v.Name
			val["value"] = v.Name
		}

		if _, ok := filterList[v.Type]; ok {
			filterList[v.Type] = append(filterList[v.Type], val)
		} else {
			filterList[v.Type] = []map[string]any{val}
		}
	}

	//补全不存在的数据
	types := []string{"base_type", "base_year", "base_producer", "base_publisher", "base_country", "base_translate", "base_version", "score", "complete"}
	for _, t := range types {
		if _, ok := filterList[t]; !ok {
			filterList[t] = []map[string]any{}
		}
	}

	//顶部加入全部
	langMap := map[string]string{
		"base_type":      config.Cfg.Lang["baseType"],
		"base_year":      config.Cfg.Lang["baseYear"],
		"base_producer":  config.Cfg.Lang["baseProducer"],
		"base_publisher": config.Cfg.Lang["basePublisher"],
		"base_country":   config.Cfg.Lang["baseCountry"],
		"base_translate": config.Cfg.Lang["baseTranslate"],
		"base_version":   config.Cfg.Lang["baseVersion"],
		"score":          config.Cfg.Lang["rating"],
		"complete":       config.Cfg.Lang["complete"],
	}
	for t, lst := range filterList {
		all := map[string]any{
			"label": langMap[t],
			"value": "",
		}
		filterList[t] = append([]map[string]any{all}, lst...)
	}

	return filterList, nil
}

/**
 * 读取rom主资源文件
 **/
func GetRomMasterResFile(dir string, filename string) string {
	for _, ext := range constant.MEDIA_EXTS {
		f := dir + filename + ext
		if utils.FileExists(f) {
			return f
		}
	}
	return ""
}

/**
 * 读取rom子资源文件列表
 **/
func GetRomSlaveResFile(dir string, prefix string) ([]string, error) {
	prefix += "__"

	files := []string{}

	extMap := utils.SliceToMap(constant.MEDIA_EXTS)

	err := filepath.Walk(dir, func(p string, fi os.FileInfo, err error) error { //遍历目录
		if err != nil {                                                         //忽略错误
			return err
		}

		if fi.IsDir() { // 忽略目录
			return nil
		}

		if strings.Index(fi.Name(), prefix) == -1 {
			return nil
		}

		romExt := strings.ToLower(utils.GetFileExt(fi.Name())) //获取文件后缀
		if _, ok := extMap[romExt]; ok {
			files = append(files, p)
		}

		return nil
	})

	return files, err
}

/**
 * 读取没有子游戏的主游戏列表
 **/
func GetGameListNotSubGame(data *request.GetGameList) ([]*db.Rom, error) {
	lists, err := (&db.Rom{}).GetNotSubRom(data.Page, data.ShowHide, data.Platform, data.Catname, data.Keyword)
	return lists, err
}
