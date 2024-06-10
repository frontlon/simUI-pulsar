package components

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/getlantern/elevate"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"simUI/db"
	"simUI/utils"
	"strings"
)

var LAST_PROCESS int = 0

/**
 * 运行游戏
 **/
func RunGame(exeFile string, cmd []string) error {

	switch runtime.GOOS {
	case "darwin":
		if exeFile == "explorer" {
			exeFile = "open"
		} else {
			if utils.IsDir(exeFile) {
				exeFile += "/Contents/MacOS/" + getDarwinAppName(exeFile)
			}
		}
	case "windows":
		/*if exeFile == "explorer" {
			exeFile = cmd[0]
		}*/

	case "linux":
		//ROOTPATH = ""
	}

	processPid, err := RunProgram(exeFile, cmd)
	if err != nil {
		return err
	}

	//保存进程id
	LAST_PROCESS = processPid

	return nil
}

// 在窗口中打开文件夹
func OpenFolderByWindow(fileName string) error {

	isDir := utils.IsDir(fileName)
	exeFile := ""
	cmd := []string{}
	switch runtime.GOOS {

	case "darwin":
		exeFile = "open"
		if isDir == true {
			cmd = []string{fileName}
		} else {
			cmd = []string{"_R", fileName}
		}
	case "windows":
		exeFile = "explorer"
		fileName = strings.ReplaceAll(fileName, `/`, `\`)
		if isDir == true {
			cmd = []string{fileName}
		} else {
			cmd = []string{"/select,", "/n,", fileName}
		}
	case "linux":
		exeFile = "open"
		if isDir == true {
			cmd = []string{fileName}
		} else {
			cmd = []string{"-R,", fileName}
		}
	}

	if exeFile != "" {
		_, err := RunProgram(exeFile, cmd)
		return err
	}
	return errors.New("Program Is Not Exists")
}

/**
 * 关闭游戏
 **/
func KillGame() error {

	if LAST_PROCESS == 0 {
		return nil
	}

	switch runtime.GOOS {
	case "darwin":
		c := exec.Command("kill", utils.ToString(LAST_PROCESS))
		c.Start()
	case "windows":
		c := exec.Command("taskkill.exe", "/T", "/PID", utils.ToString(LAST_PROCESS))
		c.Start()
	case "linux":
		c := exec.Command("kill", utils.ToString(LAST_PROCESS))
		c.Start()
	}

	LAST_PROCESS = 0
	return nil
}

// 从info.plist中读取应用程序名称
func getDarwinAppName(p string) string {

	fi, err := os.Open(p + "/Contents/Info.plist")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return ""
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	isset := false
	for {
		a, _, c := br.ReadLine()
		str := string(a)
		if isset == true {
			str = strings.Replace(str, " ", "", -1)
			str = strings.Replace(str, "\t", "", -1)
			str = strings.Replace(str, "<string>", "", -1)
			str = strings.Replace(str, "</string>", "", -1)
			return str
		}
		key := strings.Index(str, "CFBundleExecutable")
		if key > -1 {
			isset = true
		}
		if c == io.EOF {
			break
		}
	}
	return ""
}

// 播放音频
func PlayAudio(params []string) error {

	//检测rom文件是否存在
	/*if utils.FileExists(config.Cfg.Default.MusicPlayer) == false {
		return errors.New(config.Cfg.Lang["MusicPlayerNotFound"])
	}

	if err := exec.Command(config.Cfg.Default.MusicPlayer, params...).Start(); err != nil {
		return err
	}*/
	return nil
}

func RunProgram(exeFile string, cmd []string) (int, error) {
	processId := 0

	config := (&db.Config{}).GetConfig(false)

	if err := os.Chdir(filepath.Dir(exeFile)); err != nil {
		return 0, err
	}

	if config.AdminRunGame == 1 {
		fmt.Println("amin")
		//管理员方式运行
		result := elevate.Command(exeFile, cmd...)
		if err := result.Start(); err != nil {
			return 0, err
		}
		processId = result.Process.Pid
	} else {
		fmt.Println("normal")

		//普通模式运行
		result := exec.Command(exeFile, cmd...)
		if err := result.Start(); err != nil {
			fmt.Println("err", err)
			return 0, err
		}
		processId = result.Process.Pid
	}
	return processId, nil
}
