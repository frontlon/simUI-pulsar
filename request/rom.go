package request

/**
 * 读取游戏列表 GetGameList
**/
type GetGameList struct {
	Theme         string `json:"theme"`         //主题
	ShowHide      uint8  `json:"showHide"`      //是否隐藏
	ShowSubGame   uint8  `json:"showSubGame"`   //是否显示子游戏
	Platform      uint32 `json:"platform"`      //平台
	Catname       string `json:"catname"`       //分类
	CatnameLike   int    `json:"catnameLike"`   //分类模糊搜索 0精确查找 1模糊查找
	Keyword       string `json:"keyword"`       //关键字
	Letter        string `json:"letter"`        //字母索引
	Page          int    `json:"page"`          //分页数
	BaseType      string `json:"baseType"`      //资料 - 游戏类型
	BasePublisher string `json:"basePublisher"` //资料 - 发布者
	BaseYear      string `json:"baseYear"`      //资料 - 发型年份
	BaseCountry   string `json:"baseCountry"`   //资料 - 国家
	BaseTranslate string `json:"baseTranslate"` //资料 - 语言
	BaseVersion   string `json:"baseVersion"`   //资料 - 版本
	BaseProducer  string `json:"baseProducer"`  //资料 - 制作商
	Score         string `json:"score"`         //评分
	Complete      string `json:"complete"`      //是否通关
	SimpleMode    string `json:"simpleModel"`   //返回的数据结构类型
}
