package components

import (
	"simUI/constant"
	"simUI/utils"
)

// 创建升级配置文件
func UpgradeCreateBat() (string, error) {
	content := `@echo off
cd %~dp0
IF "%~1"=="" (
	echo param error.
	exit
)
set "currDir=%~dp0"
set "zipDir=%~1"
IF NOT EXIST "%zipDir%" (
	echo unzip dir is not exists.
	exit
)

echo start upgrade...
ping 127.0.0.1 -n 5 > nul
XCOPY /S /y "%zipDir%*" "%currDir%"
start "" "%currDir%simui-pulsar.exe"
RMDIR /s /q "%zipDir%"
DEL "%~f0"
`
	p := constant.ROOT_PATH + "upgrade.bat"
	if err := utils.CreateFile(p, content); err != nil {
		return "", err
	}
	return p, nil
}
