package components

import (
	"simUI/constant"
	"simUI/db"
)

// 读取全部资源目录
// typ:0全部 1图片 2图片+视频
func GetResPath(platformId uint32, typ int) map[string]string {

	thumb, snap, poster, packing, title, cassette := "", "", "", "", "", ""
	icon, gif, background, video := "", "", "", ""
	doc, strategy, files, upload := "", "", "", ""

	info := (&db.Platform{}).GetVOById(platformId, false)

	if info != nil {
		thumb = info.ThumbPath
		snap = info.SnapPath
		poster = info.PosterPath
		packing = info.PackingPath
		title = info.TitlePath
		cassette = info.CassettePath
		icon = info.IconPath
		gif = info.GifPath
		background = info.BackgroundPath
		video = info.VideoPath
		doc = info.DocPath
		strategy = info.StrategyPath
		//audio = info.AudioPath
		files = info.FilesPath
		upload = info.UploadPath
	}

	res := map[string]string{
		"thumb":      thumb,
		"snap":       snap,
		"poster":     poster,
		"packing":    packing,
		"title":      title,
		"cassette":   cassette,
		"icon":       icon,
		"gif":        gif,
		"background": background,
	}
	if typ == 0 || typ == 2 {
		res["video"] = video
	}
	if typ == 0 {
		res["doc"] = doc
		res["strategy"] = strategy
		//res["audio"] = audio
		res["file"] = files
		res["upload"] = upload
	}
	return res
}

// 读取资源类型名
func GetResExts() map[string][]string {
	res := map[string][]string{}

	picExt := constant.MEDIA_EXTS
	docExt := constant.DOC_EXTS
	fileExt := constant.FILE_EXTS
	//audioExt := constant.AUDIO_EXTS

	res["thumb"] = picExt
	res["snap"] = picExt
	res["poster"] = picExt
	res["packing"] = picExt
	res["title"] = picExt
	res["cassette"] = picExt
	res["icon"] = picExt
	res["gif"] = picExt
	res["background"] = picExt
	res["video"] = picExt
	res["doc"] = docExt
	res["strategy"] = docExt
	res["files"] = fileExt
	//res["audio"] = audioExt
	return res
}

// 根据图片类型 读取 图片路径
func GetResPathByType(typ string, platformId uint32) string {
	resPath := GetResPath(platformId, 0)
	return resPath[typ]
}
