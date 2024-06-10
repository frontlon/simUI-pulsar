package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"simUI/utils"
	"strings"
	"time"
)

var ROM_PAGE_NUM = 100 //首页每页加载rom数量

type Rom struct {
	Id            uint64
	Pid           string  // 父游戏ID，用于识别父子游戏关联
	Menu          string  // 菜单名称
	Name          string  // 别名
	Platform      uint32  // 平台
	LinkFile      string  // slnk路径
	RomName       string  // rom文件名
	SimId         uint32  // 正在使用的模拟器id
	Score         float64 // 评分
	Hide          uint8   // 是否隐藏
	Size          string  // rom文件大小
	BaseNameEn    string  // 英文名
	BaseNameJp    string  // 日文名
	BaseType      string  // 游戏类型，如RPG
	BaseYear      string  // 游戏年份
	BaseProducer  string  // 游戏出品公司
	BasePublisher string  // 游戏发行公司
	BaseCountry   string  // 游戏国家
	BaseTranslate string  // 汉化组
	BaseVersion   string  // 版本
	BaseOtherA    string  // 其他信息A
	BaseOtherB    string  // 其他信息B
	BaseOtherC    string  // 其他信息C
	BaseOtherD    string  // 其他信息D
	Pinyin        string  // 拼音索引
	InfoMd5       string  // 信息md5，用于标记rom信息或资料是否被修改
	RunNum        uint64  // 运行次数
	RunLasttime   int64   // 最后运行时间
	Complete      uint8   // 通关状态(0未通关;1已通关;2完美通关)
	SimSetting    string  // rom模拟器独立设置
}

func (*Rom) TableName() string {
	return "rom"
}

// 写入
func (m *Rom) Add() (uint64, error) {
	result := getDb().Create(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return m.Id, result.Error
}

// 批量插入
func (m *Rom) BatchAdd(romlist []*Rom, showLoading int) {

	if len(romlist) == 0 {
		return
	}

	count := len(romlist)
	if showLoading == 1 {
		utils.Loading("[2/3]开始写入缓存(1/"+utils.ToString(count)+")", "")
	}
	tx := getDb().Begin()
	for _, v := range romlist {

		tx.Create(&v)
	}
	tx.Commit()
}

//根据条件，查询多条数据

func (m *Rom) Get(showHide uint8, pages int, platform uint32, menu string, menuLike int, keyword string, baseType string, basePublisher string, baseYear string, baseCountry string, baseTranslate string, baseVersion string, baseProducer string, score string, complete string) ([]*Rom, error) {

	volist := []*Rom{}
	where := map[string]interface{}{}
	likeWhere := "1=1 "

	if platform != 0 {
		where["platform"] = platform
	}

	where["pid"] = ""

	if showHide == 0 {
		where["hide"] = showHide
	}

	if menu != "" {
		if menu == "hide" {
			where["hide"] = 1
		} else {
			if menuLike == 0 {
				where["menu"] = menu //精确搜索菜单
			} else {
				likeWhere += ` AND menu LIKE "` + menu + `%"` //模糊搜索菜单
			}
		}
	}

	if baseType != "" {
		likeWhere += ` AND base_type LIKE "%` + baseType + `%"`
	}
	if baseProducer != "" {
		likeWhere += ` AND base_producer LIKE "%` + baseProducer + `%"`
	}
	if basePublisher != "" {
		likeWhere += ` AND base_publisher LIKE "%` + basePublisher + `%"`
	}
	if baseYear != "" {
		likeWhere += ` AND base_year LIKE "` + baseYear + `%"`
	}
	if baseCountry != "" {
		where["base_country"] = baseCountry
	}
	if baseTranslate != "" {
		where["base_translate"] = baseTranslate
	}
	if baseVersion != "" {
		where["base_version"] = baseVersion
	}

	if score != "" {
		where["score"] = score
	}

	if complete != "" {
		where["complete"] = complete
	}
	if keyword != "" {
		likeWhere += `AND (name LIKE "%` + keyword + `%" or rom_name LIKE "%` + keyword + `%" or link_file LIKE "%` + keyword + `%" or pinyin LIKE "%` + keyword + `%")`
	}

	confSort := (&Platform{}).GetPlatformUiRomSort(platform) //读取排序方式

	sort := "pinyin ASC"
	switch confSort {
	case 1:
		sort = "pinyin ASC"
	case 2:
		sort = "pinyin DESC"
	case 3:
		sort = "score ASC,pinyin ASC"
	case 4:
		sort = "score DESC,pinyin ASC"
	case 5:
		sort = "base_year ASC,pinyin ASC"
	case 6:
		sort = "base_year DESC,pinyin ASC"
	case 7:
		sort = "run_num DESC,pinyin ASC"
	case 8:
		sort = "run_lasttime DESC,pinyin ASC"
	}
	mod := getDb().Select("*").Where(where).Where(likeWhere).Order(sort)

	if pages > -1 {
		offset := pages * ROM_PAGE_NUM
		mod = mod.Offset(offset).Limit(ROM_PAGE_NUM)
	}

	if menu == "" {
		//全部目录，过滤到重复模块
		mod = mod.Group("rom_name")
	}

	result := mod.Find(&volist)

	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return volist, result.Error
}

// 读取没有子游戏的rom
func (m *Rom) GetNotSubRom(pages int, showHide uint8, platform uint32, menu string, keyword string) ([]*Rom, error) {

	volist := []*Rom{}
	where := map[string]interface{}{}

	if platform != 0 {
		where["platform"] = platform
	}

	if showHide == 0 {
		where["hide"] = showHide
	}

	if menu != "" {
		if menu == "hide" {
			where["hide"] = 1
		} else {
			where["menu"] = menu
		}

	}

	likeWhere := ""
	if keyword != "" {
		likeWhere = `(name LIKE "%` + keyword + `%" or link_file LIKE "%` + keyword + `%" or pinyin LIKE "%` + keyword + `%")`
	}

	where["pid"] = ""

	offset := pages * ROM_PAGE_NUM

	//查出重复rom
	repeat := []*Rom{}
	parentName := []string{}
	getDb().Select("pid").Where("platform = ? AND pid != ''", platform).Where(likeWhere).Group("pid").Find(&repeat)
	for _, v := range repeat {
		parentName = append(parentName, v.Pid)
	}

	//查询rom数据
	if len(parentName) > 0 {
		getDb().Select("*").Where(where).Where(likeWhere).Where("rom_name not in (?)", parentName).Order("pinyin ASC").Offset(offset).Limit(ROM_PAGE_NUM).Find(&volist)
	} else {
		getDb().Select("*").Where(where).Where(likeWhere).Order("pinyin ASC").Offset(offset).Limit(ROM_PAGE_NUM).Find(&volist)
	}

	return volist, nil
}

// 读取子rom
func (*Rom) GetSubRom(platform uint32, pid string) ([]*Rom, error) {

	volist := []*Rom{}

	if platform == 0 || pid == "" {
		return volist, nil
	}

	result := getDb().Select("id,name,pid,link_file,rom_name").Where("platform=? AND pid=?", platform, pid).Group("rom_name").Order("pinyin ASC").Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return volist, result.Error
}

// 批量查询子游戏
func (*Rom) GetSubRomByFileNames(names []string) (map[string][]*Rom, error) {

	volist := []*Rom{}

	if len(names) == 0 {
		return map[string][]*Rom{}, nil
	}

	result := getDb().Select("id,name,pid,link_file").Where("pid in (?)", names).Order("pinyin ASC").Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	//封装成map
	create := map[string][]*Rom{}
	for _, v := range volist {
		create[v.Pid] = append(create[v.Pid], v)
	}

	return create, result.Error
}

// 根据id查询一条数据
func (*Rom) GetById(id uint64) (*Rom, error) {

	vo := &Rom{}

	result := getDb().Where("id=?", id).First(&vo)
	if result.Error != nil {
		fmt.Println(result.Error)
		return nil, result.Error
	}

	return vo, result.Error
}

// 根据id查询一条数据
func (*Rom) GetByIds(ids []uint64) ([]*Rom, error) {

	vo := []*Rom{}

	result := getDb().Where("id in(?)", ids).Find(&vo)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return vo, result.Error
}

// 根据拼音筛选
func (*Rom) GetByPinyin(showHide uint8, pages int, platform uint32, menu string, menuLike int, letter string) ([]*Rom, error) {
	where := map[string]interface{}{}
	likeWhere := ""
	if platform != 0 {
		where["platform"] = platform
	}

	if menu != "" {
		if menu == "hide" {
			where["hide"] = 1
		} else {
			if menuLike == 0 {
				where["menu"] = menu //精确搜索菜单
			} else {
				likeWhere += `menu LIKE "` + menu + `%"` //模糊搜索菜单
			}
		}
	}

	where["hide"] = showHide

	where["pid"] = ""
	offset := pages * ROM_PAGE_NUM
	volist := []*Rom{}
	result := getDb().Select("*").Order("pinyin ASC").Limit(ROM_PAGE_NUM).Offset(offset)

	if menu == "" {
		//全部目录，过滤到重复模块
		result = result.Group("rom_name")
	}

	if letter == "#" {

		//查询0-9数字rom
		subWhere := "pinyin LIKE '0%'"
		for i := 1; i <= 9; i++ {
			subWhere += " OR pinyin LIKE '" + utils.ToString(i) + "%'"
		}
		result.Where(where).Where(subWhere).Where(likeWhere).Find(&volist)
	} else {
		letter = strings.ToLower(letter)
		result.Where(where).Where(likeWhere).Where("pinyin LIKE ?", letter+"%").Find(&volist)
	}

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return volist, result.Error
}

// 读取一个平台下的所有md5
func (m *Rom) GetIdentByPlatform(platform uint32) ([]string, error) {
	volist := []*Rom{}
	result := getDb().Select("info_md5").Where("platform=?", platform).Find(&volist)
	if result.Error != nil {
		return []string{}, nil
	}

	md5List := []string{}
	for _, v := range volist {
		md5List = append(md5List, v.InfoMd5)
	}
	return md5List, result.Error
}

// 根据平台id，查询数据
func (*Rom) GetByPlatform(platform uint32) ([]*Rom, error) {

	volist := []*Rom{}
	result := getDb().Select("*").Where("platform = ?", platform).Order("pinyin ASC").Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return volist, result.Error
}

// 根据平台id，查询主游戏
func (*Rom) GetMasterRomByPlatform(platform uint32) ([]*Rom, error) {

	volist := []*Rom{}
	result := getDb().Select("*").Where("platform = ? AND pid = ''", platform).Order("pinyin ASC").Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return volist, result.Error
}

func (*Rom) GetByPid(platform uint32, romName string) (*Rom, error) {
	vo := &Rom{}
	result := getDb().Where("platform = ? AND pid = ?", platform, romName).First(&vo)
	if result.Error != nil {
		return nil, result.Error
	}
	return vo, result.Error
}

func (*Rom) GetByRomName(platform uint32, romName string) ([]*Rom, error) {
	volist := []*Rom{}
	result := getDb().Where("platform = ? AND rom_name = ?", platform, romName).First(&volist)
	if result.Error != nil {
		return nil, result.Error
	}
	return volist, result.Error
}

// 根据目录名读取rom
func (*Rom) GetByMenu(platform uint32, menu string) ([]*Rom, error) {

	volist := []*Rom{}

	result := getDb().Where("platform = ? AND menu = ?", platform, menu).Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return volist, result.Error
}

// 根据目录名读取rom，包含子目录
func (*Rom) GetMenuDeep(platform uint32, menu string) ([]*Rom, error) {

	volist := []*Rom{}

	result := getDb().Where("platform = ? AND menu LIKE ?", platform, menu).Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return volist, result.Error
}

// 根据目录名列表读取rom
func (*Rom) GetMasterByMenus(platform uint32, menus []string) ([]*Rom, error) {

	volist := []*Rom{}

	result := getDb().Where("platform = ? AND pid = '' and menu in (?)", platform, menus).Find(&volist)
	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return volist, result.Error
}

// 根据目录列表读取rom
func (r *Rom) GetMasterAndSubByMenus(platform uint32, menus []string) ([]*Rom, error) {

	vo := []*Rom{}
	t := r.TableName()
	subSql := `select pid from ` + t + ` where platform = ? and menu in (?) and pid =''`
	sql := `select * from ` + t + ` where (platform = ? and menu in (?)) or pid in(` + subSql + `)`
	err := getDb().Raw(sql, platform, menus, platform, menus).Scan(&vo).Error
	if err != nil {
		return nil, err
	}

	return vo, err
}

// 根据父id列表读取rom
func (*Rom) GetMasterAndSubByMasterIds(ids []uint64) ([]*Rom, error) {

	vo := []*Rom{}

	subSql := `select pid from rom where id in (?) and pid =''`
	sql := `select * from rom where id in (?) or pid in(` + subSql + `)`

	err := getDb().Raw(sql, ids, ids).Scan(&vo).Error
	if err != nil {
		return nil, err
	}

	return vo, err
}

// 根据满足条件的rom数量
func (m *Rom) Count(showHide uint8, platform uint32, menu string, menuLike int, keyword string, baseType string, basePublisher string, baseYear string, baseCountry string, baseTranslate string, baseVersion string, baseProducer string, score string, complete string) (int, error) {
	count := 0
	where := map[string]interface{}{}
	likeWhere := "1=1 "
	if platform != 0 {
		where["platform"] = platform
	}

	where["pid"] = ""

	if showHide == 0 {
		where["hide"] = showHide
	}
	if menu != "" {
		if menu == "hide" {
			where["hide"] = 1
		} else {
			if menuLike == 0 {
				where["menu"] = menu //精确搜索菜单
			} else {
				likeWhere += ` AND menu LIKE "` + menu + `%"` //模糊搜索菜单
			}
		}
	}

	if baseType != "" {
		likeWhere += ` AND base_type LIKE "%` + baseType + `%"`
	}
	if baseProducer != "" {
		likeWhere += ` AND base_producer LIKE "%` + baseProducer + `%"`
	}
	if basePublisher != "" {
		likeWhere += ` AND base_publisher LIKE "%` + basePublisher + `%"`
	}
	if baseYear != "" {
		likeWhere += ` AND base_year LIKE "` + baseYear + `%"`
	}
	if baseCountry != "" {
		where["base_country"] = baseCountry
	}
	if baseTranslate != "" {
		where["base_translate"] = baseTranslate
	}
	if baseTranslate != "" {
		where["base_version"] = baseVersion
	}
	if score != "" {
		where["score"] = score
	}
	if complete != "" {
		where["complete"] = complete
	}
	if keyword != "" {
		likeWhere += `AND (name LIKE "%` + keyword + `%" or link_file LIKE "%` + keyword + `%" or pinyin LIKE "%` + keyword + `%")`
	}

	result := getDb().Table(m.TableName()).Where(where).Where(likeWhere).Group("rom_name").Count(&count)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return count, result.Error
}

// 根据拼音筛选
func (m *Rom) CountByPinyin(showHide uint8, pages int, platform uint32, menu string, menuLike int, letter string) (int, error) {
	where := map[string]interface{}{}
	likeWhere := ""
	if platform != 0 {
		where["platform"] = platform
	}

	if menu != "" {
		if menu == "hide" {
			where["hide"] = 1
		} else {
			if menuLike == 0 {
				where["menu"] = menu //精确搜索菜单
			} else {
				likeWhere += `menu LIKE "` + menu + `%"` //模糊搜索菜单
			}
		}
	}

	where["hide"] = showHide

	where["pid"] = ""
	count := 0

	result := &gorm.DB{}
	if letter == "#" {

		//查询0-9数字rom
		subWhere := "pinyin LIKE '0%'"
		for i := 1; i <= 9; i++ {
			subWhere += " OR pinyin LIKE '" + utils.ToString(i) + "%'"
		}
		result = getDb().Table(m.TableName()).Where(where).Where(subWhere).Where(likeWhere).Group("rom_name").Count(&count)
	} else {
		letter = strings.ToLower(letter)
		result = getDb().Table(m.TableName()).Where(where).Where(likeWhere).Where("pinyin LIKE ?", letter+"%").Group("rom_name").Count(&count)
	}

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	return count, result.Error
}

// 根据菜单读取游戏数
func (m *Rom) CountGroupMenu(platform uint32) (map[string]int64, error) {
	volist := []*struct {
		Count int64
		Menu  string
	}{}
	result := getDb().Table(m.TableName()).Select("menu,count(*) as count").Where("platform=?", platform).Group("menu").Find(&volist)
	if result.Error != nil {
		return nil, result.Error
	}
	data := map[string]int64{}
	for _, v := range volist {
		data[v.Menu] = v.Count
	}
	return data, nil
}

// 更新名称
func (m *Rom) UpdateName() error {
	create := map[string]interface{}{
		"name":     m.Name,
		"pinyin":   m.Pinyin,
		"info_md5": m.InfoMd5,
	}

	result := getDb().Table(m.TableName()).Where("platform = ? AND rom_name = ?", m.Platform, m.RomName).Updates(create)
	return result.Error
}

// 更新隐藏状态
func (m *Rom) UpdateHide() error {
	result := getDb().Table(m.TableName()).Where("rom_name = ?", m.RomName).Update("hide", m.Hide)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

// 更新模拟器
func (m *Rom) UpdateSimIdByIds(romIds []uint64, simId uint32) error {
	result := getDb().Table(m.TableName()).Where("id in (?)", romIds).Update("sim_id", simId)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

// 更新rom模拟器配置
func (m *Rom) UpdateSimSettingById(romId uint64, data string) error {
	result := getDb().Table(m.TableName()).Where("id = ?", romId).Update("sim_setting", data)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

// rom改名
func (m *Rom) UpdateRomName(platform uint32, oldName, newName string) error {

	update := map[string]interface{}{
		"rom_name": newName,
	}

	result := getDb().Table(m.TableName()).Where("platform = ? AND rom_name=?", platform, oldName).Update(update)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

// 更新游戏资料
func (m *Rom) UpdateRomBase() error {

	create := map[string]string{
		"name":           m.Name,
		"base_type":      m.BaseType,
		"base_year":      m.BaseYear,
		"base_producer":  m.BaseProducer,
		"base_publisher": m.BasePublisher,
		"base_country":   m.BaseCountry,
		"base_translate": m.BaseTranslate,
		"base_version":   m.BaseVersion,
		"score":          utils.ToString(m.Score),
		"base_name_en":   m.BaseNameEn,
		"base_name_jp":   m.BaseNameJp,
		"base_other_a":   m.BaseOtherA,
		"base_other_b":   m.BaseOtherB,
		"base_other_c":   m.BaseOtherC,
		"base_other_d":   m.BaseOtherD,
		"info_md5":       m.InfoMd5,
	}

	result := getDb().Table(m.TableName()).Where("rom_name=?", m.RomName).Updates(create)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

// 更新运行次数和最后运行时间
func (m *Rom) UpdateRunNumAndTime(id uint64) error {
	result := getDb().Table(m.TableName()).Where("id=?", id).Update("run_num", gorm.Expr("run_num + 1"))
	getDb().Table(m.TableName()).Where("id=?", id).Update("run_lasttime", time.Now().Unix())
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

// 更新一个游戏的关系（设为子游戏/解绑子游戏）
func (m *Rom) UpdatePidByRomName(platform uint32, romName, newName string) error {
	result := getDb().Table(m.TableName()).Where("platform = ? AND rom_name = ?", platform, romName).Update("pid", newName)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

// 解绑一个游戏的所有子游戏
func (m *Rom) UpdatePidByPid(platform uint32, oldPid, newPid string) error {
	result := getDb().Table(m.TableName()).Where("platform = ? AND pid = ?", platform, oldPid).Update("pid", newPid)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

// 更新一个字段
func (m *Rom) UpdateOneField(id uint64, field string, val any) error {
	result := getDb().Table(m.TableName()).Where("id = ?", id).Update(field, val)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

// 更新rom的链接文件地址
func (m *Rom) UpdateRoomLinkById() error {

	update := map[string]string{
		"menu":      m.Menu,
		"link_file": m.LinkFile,
	}

	result := getDb().Table(m.TableName()).Where("id = ?", m.Id).Updates(update)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

// 替换字段中的内容
func (m *Rom) ReplaceFieldData(platform uint32, field, from, to string) error {
	result := getDb().Table(m.TableName()).Where("platform", platform).Update(field, gorm.Expr("REPLACE("+field+", ?, ?)", from, to))
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

// 删除一个rom
func (m *Rom) DeleteById(id uint64) error {
	result := getDb().Where("id=? ", id).Delete(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

// 删除by rom_name
func (m *Rom) DeleteByName(romName string) error {
	result := getDb().Where("rom_name=? ", romName).Delete(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

// 删除一个游戏的所有子游戏
func (m *Rom) DeleteSubRom(platform uint32, pid string) error {
	result := getDb().Where("platform = ? AND pid=? ", platform, pid).Delete(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

// 删除一个平台下的所有rom数据
func (m *Rom) DeleteByPlatform() error {
	result := getDb().Where("platform=? ", m.Platform).Delete(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

func (m *Rom) DeleteByInfoMd5s(platform uint32, md5s []string) {
	if len(md5s) == 0 {
		return
	}
	result := getDb().Where("platform = ? AND info_md5 in (?)", platform, md5s).Delete(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return
}

// 删除不存在的平台下的所有rom
func (m *Rom) DeleteByNotPlatform(platforms []string) error {
	result := getDb().Where("platform not in (?)", platforms).Delete(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

// 清空表数据
func (m *Rom) Truncate() error {
	result := getDb().Delete(&m)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

// 有空游戏统计信息
func (m *Rom) TruncateGameStat() error {

	create := map[string]string{
		"run_lasttime": "0",
		"run_num":      "0",
	}

	result := getDb().Table(m.TableName()).Updates(create)
	if result.Error != nil {
		fmt.Println(result.Error)
	}
	return result.Error
}

// 计算rom信息md5，主要用于确定资料文件是否被手动修改
func (m *Rom) CreateInfoMd5() string {
	data := []string{
		m.Menu,
		m.Name,
		utils.ToString(m.Platform),
		m.Pid,
		m.LinkFile,
		m.RomName,
		utils.ToString(m.Hide),
		m.Size,
		utils.ToString(m.RunNum),
		utils.ToString(m.RunLasttime),
		utils.ToString(m.Score),
		utils.ToString(m.Complete),
		utils.ToString(m.SimId),
		m.BaseType,
		m.BaseYear,
		m.BaseProducer,
		m.BasePublisher,
		m.BaseCountry,
		m.BaseTranslate,
		m.BaseVersion,
		m.BaseNameEn,
		m.BaseNameJp,
		m.BaseOtherA,
		m.BaseOtherB,
		m.BaseOtherC,
		m.BaseOtherD,
		m.SimSetting,
		m.Pinyin,
	}
	return utils.Md5(strings.Join(data, ","))
}
