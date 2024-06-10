package constant

// 是否为开发环境
var DEV = false

var DB_ADD_MAX_NUM = 999 //数据库每次查询/写入的最大数量

var VERSION_NO = "" //当前软件版本号

// doc文档支持的扩展名
var DOC_EXTS = []string{".txt", ".html", ".htm", ".md"}

// 支持的图片类型
var MEDIA_EXTS = []string{".png", ".jpg", ".gif", ".webp", ".jpeg", ".ico", ".mp4", ".webm", "avif"}

// 支持的音频类型
var AUDIO_EXTS = []string{".mp3", ".dmi", ".wav", ".wma"}

// 可直接运行的doc文档支持的扩展名
var FILE_EXTS = []string{
	".html", ".htm", ".mht", ".mhtml", ".url",
	".chm", ".pdf", ".doc", ".docx", ".ppt", ".pptx", ".xls", ".xlsx", ".rtf",
	".exe", ".com", ".cmd", ".bat", ".lnk",
}

// 可直接运行的扩展名
var RUN_EXTS = []string{".exe", ".cmd", ".bat"}

// 可直接通过explorer运行的扩展名
var EXPLORER_EXTS = []string{".lnk"}

/*EXPLORER_EXTS = []string{
	".lnk", ".html", ".htm", ".mht", ".mhtml", ".url",
	".chm", ".doc", ".docx", ".ppt", ".pptx", ".xls", ".xlsx", ".rtf",
} //通过explorer运行的扩展名*/

// 主题标识配置
var UI_DEFAULT = "Default"
var UI_PLAYNITE = "Playnite"
var UI_TINY = "Tiny"

// 平台资源文件夹定义
var RES_DIR = map[string]string{
	"rom":        "roms",
	"thumb":      "thumbs",
	"snap":       "snaps",
	"poster":     "poster",
	"packing":    "packing",
	"title":      "title",
	"cassette":   "cassette",
	"icon":       "icon",
	"gif":        "gif",
	"background": "background",
	"video":      "video",
	"doc":        "docs",
	"strategy":   "strategies",
	"audio":      "audio",
	"file":       "files",
	"upload":     "uploads",
	"link":       "links",
}

// 默认展示图排序
var DefaultThumbOrders = []string{
	"thumb", "snap", "poster", "packing", "title", "cassette", "icon", "gif", "background", "video",
}

// 默认过滤器
var DefaultListColumns = []string{
	"BaseNameEn", "BaseNameJp", "BaseType", "BaseYear", "BasePublisher", "BaseProducer",
	"BaseCountry", "BaseTranslate", "BaseVersion", "Score", "Complete", "Menu",
}
