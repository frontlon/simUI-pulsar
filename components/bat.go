package components

import (
	"simUI/utils"
)

func CreateGameBat(batPath, gamePath string) error {
	relPath := utils.GetRelPathByTowPath(batPath, gamePath)
	f := utils.GetFileNameAndExt(gamePath)

	bat := ""
	bat += "cd %~dp0\r\n"
	bat += `cd "` + relPath + `"` + "\r\n"
	bat += `start "" "` + f + `"` + "\r\n"
	bat = utils.Utf8ToGbk(bat)
	return utils.CreateFile(batPath, bat)

}
