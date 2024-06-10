package modules

import (
	"encoding/json"
	"errors"
	"regexp"
	"simUI/components"
	"simUI/constant"
	"simUI/db"
	"simUI/request"
	"simUI/utils"
	"strings"
)

func GetAllPlatform() []*db.PlatformVO {
	return (&db.Platform{}).GetPlatformVOList(false)
}

func GetPlatformById(id uint32) (*db.PlatformVO, error) {
	return (&db.Platform{}).GetVOById(id, false), nil
}

// 读取全部平台 - 原始信息
func GetPlatformOriginal() ([]*db.Platform, error) {
	return (&db.Platform{}).GetAll()
}

// 读取平台UI设置
func GetPlatformUi(id uint32, theme string) (any, error) {

	ui := db.PlatformUI{}
	if id == 0 {
		conf := (&db.Config{}).GetConfig(false)
		ui = conf.PlatformUi
	} else {
		platformInfo := (&db.Platform{}).GetVOById(id, false)
		if platformInfo != nil {
			ui = platformInfo.Ui
		} else {
			conf := (&db.Config{}).GetConfig(false)
			ui = conf.PlatformUi
		}
	}

	switch theme {
	case constant.UI_DEFAULT:
		return ui.Default, nil
	case constant.UI_PLAYNITE:
		return ui.Playnite, nil
	case constant.UI_TINY:
		return ui.Tiny, nil
	default:
		return nil, errors.New("Skin Not Exists")
	}
}

// 读取全部平台标签
func GetAllPlatformTag() ([]string, error) {
	tags := (&db.Platform{}).GetAllPlatformTag()
	return tags, nil
}

// 添加一个新平台
func AddPlatform(name string) (*db.Platform, error) {

	conf := (&db.Config{}).GetConfig(false)
	ui := conf.PlatformUi
	uiStr, _ := json.Marshal(ui)

	exts := ".20,.40,.48,.58,.60,.78,.2hd,.2mg,.32x,.3ds,.3dsx,.68k,.7z,.7zip,.88d,.98d,.a0,.a26,.a52,.a78,.abs,.adf,.adz,.agb,.app,.arc,.atr,.atx,.aus,.axf,.b0,.b5t,.b6t,.bat,.bin,.bml,.bps,.bs,.bsx,.bwt,.car,.cas,.cbn,.ccd,.cci,.cdi,.cdm,.cdt,.cfg,.cgb,.ch8,.chd,.ciso,.cmd,.cof,.col,.com,.conf,.crt,.cso,.csw,.cue,.cxi,.d13,.d1m,.d2m,.d4m,.d64,.d6z,.d71,.d7z,.d80,.d81,.d82,.d88,.d8z,.d98,.dat,.dci,.dcm,.dff,.dim,.dmg,.dms,.do,.dol,.dsk,.dup,.dx2,.elf,.eur,.exe,.fd,.fdd,.fdi,.fds,.fig,.fm2,.fs-uae,.g41,.g4z,.g64,.g6z,.gb,.gba,.gbc,.gbs,.gbz,.gcm,.gcz,.gd3,.gd7,.gdi,.gen,.gg,.gz,.hdd,.hdf,.hdi,.hdm,.hdn,.hdz,.img,.int,.ipf,.ips,.iso,.isz,.j64,.jag,.jap,.jma,.k7,.kcr,.ldb,.lha,.lnx,.lst,.lzx,.m3u,.m5,.m7,.mb,.md,.mdf,.mds,.mdx,.mgd,.mgh,.mgw,.msa,.mv,.mx1,.mx2,.n64,.nca,.ndd,.nds,.nes,.nez,.nfse,.ngc,.ngp,.ngpc,.nhd,.npc,.nrg,.nro,.nsf,.nsfe,.nso,.nsp,.o,.obx,.p,.p00,.p64,.pal,.pbp,.pc2,.pce,.pdi,.po,.prg,.pro,.prof,.prx,.psexe,.pzx,.qd,.rar,.raw,.ri,.rom,.rpx,.rvz,.rzx,.sap,.sc,.scl,.scp,.sfc,.sg,.sgb,.sgg,.sgx,.sk,.smc,.smd,.sms,.sna,.st,.stx,.swc,.swf,.t64,.t81,.tap,.tar,.tfd,.tgc,.thd,.toc,.trd,.tzx,.u1,.uae,.ufo,.unf,.unif,.ups,.usa,.uze,.v64,.vb,.vboy,.vec,.vms,.voc,.vpk,.vsf,.wad,.wav,.wbfs,.wia,.woz,.ws,.wsc,.wud,.wux,.x64,.x6z,.xbe,.xci,.xdf,.xex,.xfd,.xml,.z64,.z80,.zip"

	m := &db.Platform{
		Name:    name,
		RomExts: exts,
		Pinyin:  utils.TextToPinyin(name),
		Ui:      string(uiStr),
	}

	id, err := m.Add()

	if err != nil {
		return nil, err
	}

	m.Id = id

	return m, nil
}

// 更新平台信息
func UpdatePlatform(data request.UpdatePlatform) (*db.Platform, error) {

	//路径转为相对目录
	data.RootPath = utils.ToRelPath(data.RootPath, "")
	data.Icon = utils.ToRelPath(data.Icon, "")

	//处理rom目录
	if len(data.RomPath) > 0 {
		for k, p := range data.RomPath {
			p = utils.ToRelPath(p, "")
			data.RomPath[k] = p
		}
	}
	data.RomPath = utils.SliceRemoveEmpty(data.RomPath)
	data.RomPath = utils.SliceRemoveDuplicate(data.RomPath)

	//处理扩展名
	if len(data.RomExts) > 0 {
		for k, ext := range data.RomExts {
			ext = strings.ToLower(ext)
			ext = strings.Trim(ext, " ")
			if strings.Contains(ext, ".") == false {
				ext = "." + ext
			}
			data.RomExts[k] = ext
		}
	}
	data.RomExts = utils.SliceRemoveEmpty(data.RomExts)
	data.RomExts = utils.SliceRemoveDuplicate(data.RomExts)

	create := map[string]any{
		"name":      data.Name,
		"icon":      data.Icon,
		"tag":       data.Tag,
		"rom_exts":  strings.Join(data.RomExts, ","),
		"root_path": data.RootPath,
		"rom_path":  strings.Join(data.RomPath, ";"),
		"hide_name": data.HideName,
		"pinyin":    utils.TextToPinyin(data.Name),
	}

	err := (&db.Platform{}).UpdateById(data.Id, create)
	if err != nil {
		utils.WriteLog(err.Error())
		return nil, err
	}

	platform, _ := (&db.Platform{}).GetById(data.Id)

	//创建rom目录
	absPath := utils.ToAbsPath(platform.RootPath, "") + "/"
	utils.CreateDir(absPath + "roms")

	//创建资源目录
	paths := components.GetResPath(data.Id, 0)
	for _, v := range paths {
		utils.CreateDir(v)
	}

	//创建资料文件
	absRootPath := utils.ToAbsPath(data.RootPath, "")
	csvFile := absRootPath + "/" + constant.ROMBASE_FILE_NAME + ".csv"
	if !utils.FileExists(csvFile) {
		utils.CreateFile(csvFile, "")
	}

	//重新读取数据
	platform, _ = (&db.Platform{}).GetById(data.Id)

	return platform, nil
}

// 更新平台简介
func UpdatePlatformDesc(id uint32, desc string) error {

	platformInfo := (&db.Platform{}).GetVOById(id, false)

	//处理图片地址
	if desc != "" {
		desc = strings.TrimSpace(desc)
		re := regexp.MustCompile(utils.WailsPathPattern)
		desc = re.ReplaceAllStringFunc(desc, func(match string) string {
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
			return utils.WailsPathEncode(rel, false)
		})
	}

	platform := &db.Platform{
		Id:   id,
		Desc: desc,
	}
	err := platform.UpdateDescById()
	if err != nil {
		utils.WriteLog(err.Error())
		return err
	}
	return nil
}

// 更新平台排序
func UpdatePlatformSort(lists []uint32) error {
	if len(lists) == 0 {
		return nil
	}
	for k, pfId := range lists {
		platform := &db.Platform{
			Id:   pfId,
			Sort: uint32(k + 1),
		}
		err := platform.UpdateSortById()
		if err != nil {
			utils.WriteLog(err.Error())
			return err
		}
	}
	return nil
}

// 删除一个平台
func DelPlatform(platformId uint32) error {

	//删除rom数据
	if err := (&db.Rom{Platform: platformId}).DeleteByPlatform(); err != nil {
		utils.WriteLog(err.Error())
	}

	//删除模拟器
	if err := (&db.Simulator{Platform: platformId}).DeleteByPlatform(); err != nil {
		utils.WriteLog(err.Error())
	}

	//删除平台
	if err := (&db.Platform{Id: platformId}).DeleteById(); err != nil {
		utils.WriteLog(err.Error())
		return err
	}
	return nil
}

// 更新平台UI
func UpdatePlatformUi(id uint32, theme string, data string) error {

	ui := db.PlatformUI{}

	if id == 0 {
		//更新全部UI
		conf := (&db.Config{}).GetConfig(true)
		ui = conf.PlatformUi
	} else {
		//更新单个平台
		pfConf := (&db.Platform{}).GetVOById(id, true)
		ui = pfConf.Ui
	}

	if theme == constant.UI_DEFAULT {
		defUi := db.PlatformUIDefault{}
		_ = json.Unmarshal([]byte(data), &defUi)
		ui.Default = defUi
	} else if theme == constant.UI_PLAYNITE {
		pn := db.PlatformUIPlaynite{}
		_ = json.Unmarshal([]byte(data), &pn)
		ui.Playnite = pn
	} else if theme == constant.UI_TINY {
		ti := db.PlatformUITiny{}
		_ = json.Unmarshal([]byte(data), &ti)
		ui.Tiny = ti
	}

	ui.Default.BackgroundImage = utils.WailsPathDecode(ui.Default.BackgroundImage)
	ui.Playnite.BackgroundImage = utils.WailsPathDecode(ui.Playnite.BackgroundImage)
	ui.Tiny.BackgroundImage = utils.WailsPathDecode(ui.Tiny.BackgroundImage)

	ui.Default.BackgroundMask = utils.WailsPathDecode(ui.Default.BackgroundMask)
	ui.Playnite.BackgroundMask = utils.WailsPathDecode(ui.Playnite.BackgroundMask)
	ui.Tiny.BackgroundMask = utils.WailsPathDecode(ui.Tiny.BackgroundMask)

	uiJson, _ := json.Marshal(ui)

	//更新数据
	var err error
	if id == 0 {
		_ = (&db.Platform{}).ClearPlatformUi()
		err = (&db.Config{}).UpdateOne("PlatformUi", string(uiJson))
	} else {
		err = (&db.Platform{}).UpdateOneField(id, "ui", string(uiJson))
	}
	(&db.Platform{}).DelPlatformVOCache() //清空缓存
	return err
}
