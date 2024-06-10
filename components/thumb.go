package components

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"simUI/config"
	"simUI/constant"
	"simUI/db"
	"simUI/utils"
	"strings"
	"time"
)

type DownThumbs struct {
	Width  int
	Height int
	Ext    string
	ImgUrl string
	Type   string
}

// 创建图片缩略图
func CreateThumbnail(src string, dst string) error {

	var imgWidth = 320
	var quality = 80

	return utils.CreateThumbnail(src, dst, imgWidth, quality)
}

// 备份老图片
func BackupPic(p string) error {

	//备份图片
	date := time.Now().Format("2006-01-02")

	folder := constant.CACHE_THUMB_PATH + date
	//检测bak文件夹是否存在，不存在则创建bak目录
	if !utils.DirExists(folder) {
		if err := utils.CreateDir(folder); err != nil {
			return err
		}
	}

	//开始备份
	fileName := utils.GetFileName(p)
	fileExt := utils.GetFileExt(p)
	bakFile := folder + "/" + fileName + "_" + utils.ToString(time.Now().UnixNano()) + fileExt //生成备份文件名
	if err := utils.FileCopy(p, bakFile); err != nil {
		return err
	}

	return nil
}

// 百度图片搜索
func SearchThumbsForBaidu(keyword string, page int) (map[string][]DownThumbs, error) {
	size := 30

	configData := (&db.Config{}).GetConfig(false)

	postUrl := configData.SearchEnginesBaidu
	num := page * size
	keyword = url.QueryEscape(keyword)
	postUrl = strings.ReplaceAll(postUrl, "{$keyword}", keyword)
	postUrl = strings.Replace(postUrl, "{$NumIndex}", utils.ToString(num), 1)
	postUrl = strings.Replace(postUrl, "{$pageNum}", utils.ToString(size), 1)

	//整理http请求体
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, postUrl, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// 添加请求头
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	respMap := make(map[string]interface{})
	bstr := string(b)
	bstr = strings.ReplaceAll(bstr, `\'`, "'")
	err = json.Unmarshal([]byte(bstr), &respMap)

	//被识别为蜘蛛，被拦截
	if _, ok := respMap["antiFlag"]; ok {
		return nil, errors.New(respMap["message"].(string))
	}

	//请求成功，开始组装消息体
	respList := map[string][]DownThumbs{}
	list := []DownThumbs{}
	if _, ok := respMap["data"]; ok {
		for _, v := range respMap["data"].([]interface{}) {
			vo := v.(map[string]interface{})

			if _, ok = vo["thumbURL"]; !ok {
				continue
			}

			width := 0
			height := 0
			ext := ""
			if _, ok = vo["width"]; ok {
				width = utils.ToInt(vo["width"].(float64))
			}
			if _, ok = vo["height"]; ok {
				height = utils.ToInt(vo["height"].(float64))
			}
			if _, ok = vo["type"]; ok {
				ext = vo["type"].(string)
			}

			stu := DownThumbs{
				Width:  width,
				Height: height,
				Ext:    ext,
				ImgUrl: vo["thumbURL"].(string),
			}
			list = append(list, stu)
		}
	}

	if len(list) > 0 {
		respList["百度图片搜索结果"] = list
	}

	return respList, nil
}

// hfsDB图片搜索
func SearchThumbsForHfsDB(keyword string, page int) (map[string][]DownThumbs, error) {

	if !utils.IsEnglish(keyword) {
		keyword = utils.BaiduTranslate(keyword, "en")
	}
	limit := 10
	offset := page * limit
	body, err := GetHfsDbGameList(keyword, limit, offset)
	if err != nil {
		fmt.Println(err)
		return nil, errors.New(config.Cfg.Lang["tipHfsDbError"])
	}

	respList := map[string][]DownThumbs{}

	for _, v := range body.Results {
		medias := []DownThumbs{}

		if len(v.Medias) == 0 {
			continue
		}
		for _, m := range v.Medias {
			if m.IsImage == true {
				md := DownThumbs{
					Width:  m.ResX,
					Height: m.ResY,
					Ext:    m.Extension,
					ImgUrl: m.File,
					Type:   m.Type,
				}
				medias = append(medias, md)
			}
		}
		respList[v.NameEn] = medias
	}
	return respList, nil
}
