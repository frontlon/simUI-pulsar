package request

type UpdatePlatform struct {
	Id       uint32   `json:"id"`
	Name     string   `json:"name"`
	Icon     string   `json:"icon"`
	Tag      string   `json:"tag"`
	RomExts  []string `json:"romExts"`
	RootPath string   `json:"rootPath"`
	RomPath  []string `json:"romPath"`
	HideName uint8    `json:"hideName"`
}
