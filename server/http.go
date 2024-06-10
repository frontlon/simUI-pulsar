package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"simUI/components"
	"simUI/config"
	"simUI/db"
	"simUI/utils"
	"time"
)

var Addr = ""

func StartHttpServer() {

	if Addr != "" {
		return
	}

	port := 0
	start := 10544

	for i := start; i <= start+10; i++ {
		//检查端口占用
		if utils.IsPortInUse(i) {
			continue
		}

		port = i
		break
	}

	if port == 0 {
		fmt.Println("no port available")
		return
	}

	Addr = "http://127.0.0.1:" + utils.ToString(port)

	go func() {
		//启动http服务
		mux := http.NewServeMux()
		mux.HandleFunc("/uploadFile", uploadFileHandler)
		err := http.ListenAndServe(":"+utils.ToString(port), mux)
		if err != nil {
			Addr = ""
			fmt.Println("ListenAndServe start fail", err)
		}
	}()

	return
}

// 上传文件到本地磁盘
func uploadFileHandler(w http.ResponseWriter, r *http.Request) {

	// 设置CORS头部允许所有源
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != "POST" {
		w.Write(resp("", errors.New("非POST方法")))
		return
	}

	// 获取POST参数
	romId := uint64(utils.ToInt(r.PostFormValue("id")))
	typ := r.PostFormValue("type")
	file, finfo, err := r.FormFile("file")
	if err != nil {
		w.Write(resp("", errors.New("没有选择文件")))
		return
	}
	defer file.Close()
	if romId == 0 {
		w.Write(resp("", errors.New("romid不能为空")))
		return
	}
	if typ == "" {
		w.Write(resp("", errors.New("type不能为空")))
		return
	}

	vo, err := (&db.Rom{}).GetById(romId)
	if err != nil {
		w.Write(resp("", err))
		return
	}

	res := components.GetResPathByType(typ, vo.Platform)
	if res == "" {
		w.Write(resp("", errors.New(config.Cfg.Lang["noSetThumbDir"])))
		return
	}

	//检查是否存在主图
	master := components.GetMasterRes(typ, vo.Platform, vo.RomName)
	dst := "" //目标路径
	ext := utils.GetFileExt(finfo.Filename)
	if master == "" {
		//设为主图
		dst = res + "/" + vo.RomName + ext
	} else {
		//设为子图
		nano := utils.ToString(time.Now().UnixNano())
		dst = res + "/" + vo.RomName + "/" + vo.RomName + "_" + nano + ext
	}

	//检查目录是否存在
	pth := utils.GetFilePath(dst)
	if !utils.DirExists(pth) {
		utils.CreateDir(pth)
	}

	// 创建输出文件
	out, err := os.Create(dst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// 将上传的文件写入输出文件
	if _, err = io.Copy(out, file); err != nil {
		w.Write(resp("", err))
		return
	}

	w.Write(resp(utils.WailsPathEncode(dst, true), nil))
}

func resp(data any, err error) []byte {

	r := map[string]any{
		"code": 0,
		"data": data,
		"err":  "",
	}

	if err != nil {
		r = map[string]any{
			"code": 0,
			"err":  err.Error(),
		}
	}

	js, _ := json.Marshal(r)
	return js
}
