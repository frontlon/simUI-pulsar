package utils

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"simUI/constant"
	"strings"
)

/*
从路径中读取文件名，不包含扩展名
*/
func GetFileName(p string) string {
	if p == "" {
		return ""
	}
	base := filepath.Base(p)
	str := strings.TrimSuffix(base, path.Ext(base))
	return str
}

/*
从路径中读取文件名扩展名
返回 .ext
*/
func GetFileExt(p string) string {
	if p == "" {
		return ""
	}
	return path.Ext(p)
}

/*
获取文件的路径部分
*/
func GetFilePath(p string) string {
	if p == "" {
		return ""
	}
	paths, _ := filepath.Split(p)
	return paths
}

/*
获取文件的绝对路径
*/
func GetFileAbsPath(p string) string {
	if p == "" {
		return ""
	}
	dir := filepath.Dir(p)
	return strings.Replace(dir, "./", "", 1)
}

/*
路径转换为相对路径
如果prefix为空，则默认为软件根目录
*/
func ToRelPath(p string, prefix string) string {

	// \file.exe这种路径在windows下，直接返回
	if strings.HasPrefix(p, "\\") && runtime.GOOS == "windows" {
		return p
	}

	if prefix == "" {
		prefix = constant.ROOT_PATH
	}
	if p == "" {
		return ""
	}

	p = ReplcePathSeparator(p) //替换 \

	if !IsAbsPath(p) {
		return p
	}

	prefix = ReplcePathSeparator(prefix) //替换 \
	p = strings.TrimRight(p, `/`)        //删除结尾 /
	p = strings.Replace(p, prefix, "", 1)

	if p != "" && p[0] == '/' {
		p = strings.TrimLeft(p, `/`) //删除开头 /
	}

	return p
}

/*
路径转换为绝对路径
*/
func ToAbsPath(p, prefix string) string {
	if p == "" {
		return ""
	}
	if !filepath.IsAbs(p) {

		// \file.exe这种路径在windows下，读取指定盘符根目录
		if strings.HasPrefix(p, "\\") && runtime.GOOS == "windows" {
			//volume := filepath.VolumeName(constant.ROOT_PATH)
			//return volume + p
			return p
		}

		if prefix == "" {
			prefix = constant.ROOT_PATH
		}

		if !strings.HasSuffix(prefix, "/") && p[0] != '/' {
			prefix = prefix + "/"
		}

		return prefix + p
	}
	return p
}

/*
检查路径是否为绝对路径
*/
func IsAbsPath(p string) bool {
	if p == "" {
		return false
	}
	return filepath.IsAbs(p)
}

/*
从路径中读取文件名+扩展名
*/
func GetFileNameAndExt(p string) string {
	if p == "" {
		return ""
	}
	return filepath.Base(p)
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	if path == "" {
		return false
	}
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

/**
 * 检测文件是否存在（文件夹也返回false）
 **/
func FileExists(path string) bool {

	if path == "" {
		return false
	}

	finfo, err := os.Stat(path)
	isset := false
	if err != nil || finfo.IsDir() == true {
		isset = false
	} else {
		isset = true
	}
	return isset
}

/**
 * 检测路径是否存在
 **/
func DirExists(p string) bool {
	if p == "" {
		return false
	}
	ff, err := os.Stat(p)
	if err != nil {
		return false
	}
	if ff.IsDir() == false {
		return false
	}
	return true
}

/**
 * 判断文件夹或文件是否存在
 **/
func IsExist(f string) bool {
	if f == "" {
		return false
	}
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

/**
 * 扫描目录和子目录，读取所有文件
 **/
func ScanDirAndSubDir(dir string) ([]string, error) {

	files := []string{}

	err := filepath.Walk(dir, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		if err != nil { //忽略错误
			return err
		}

		if fi.IsDir() { // 忽略目录
			return nil
		}

		files = append(files, filename)

		return nil
	})

	return files, err
}

/**
 * 扫描当前目录文件，不包含子目录
 **/
func ScanCurrentDir(dir string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(dir)

	if err != nil {
		return nil, err
	}
	return files, nil
}

/*
*
检查文件夹是否为空
*/
func IsDirEmpty(dirname string) bool {
	f, err := os.Open(dirname)
	if err != nil {
		// 打开文件夹失败，则认为该目录为空
		return true
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err != nil {
		// 如果无法读取目录项，则认为该目录为空
		return true
	}

	// 如果读取到了目录项，则认为该目录不为空
	return false
}

/**
 *  根据两个路径，读取出from对应to的相对路径
 **/
func GetRelPathByTowPath(from string, to string) string {

	separator := string(os.PathSeparator)
	root, _ := os.Getwd()
	root += separator

	root = strings.ReplaceAll(root, `/`, separator)
	root = strings.ReplaceAll(root, `\`, separator)
	from = strings.ReplaceAll(from, `/`, separator)
	from = strings.ReplaceAll(from, `\`, separator)
	to = strings.ReplaceAll(to, `/`, separator)
	to = strings.ReplaceAll(to, `\`, separator)

	newrel := ""

	if strings.Contains(to, root) == true {

		from = strings.Replace(GetFilePath(from), root, "", -1)
		to = strings.Replace(GetFilePath(to), root, "", -1)

		fromArr := strings.Split(from, separator)
		toArr := strings.Split(to, separator)

		repeatArr := []string{}

		//获取重复路径
		for k, _ := range fromArr {
			if fromArr[k] == toArr[k] && fromArr[k] != "" {
				repeatArr = append(repeatArr, fromArr[k])
			} else {
				break
			}
		}

		if len(repeatArr) > 0 {
			repeat := strings.Join(repeatArr, separator) + separator
			from = strings.Replace(from, repeat, "", -1)
			to = strings.Replace(to, repeat, "", -1)
		}
		from2 := strings.Split(from, separator)
		for _, v := range from2 {
			if v == "" {
				continue
			}
			newrel += ".." + separator
		}

		newrel += to

	} else {
		//不在simui目录内，直接读取全路径
		newrel = GetFilePath(to)
	}

	return newrel
}

// 将路径拆解为路径和参数，第一行路径，第二行参数
func ParsePathFile(input string) (string, []string) {
	input = strings.ReplaceAll(input, "\r\n", "\n")
	arr := strings.Split(input, "\n")
	path := arr[0]
	args := []string{}
	if len(arr) > 1 && arr[1] != "" {
		cmds := strings.Split(arr[1], " ")
		for _, v := range cmds {
			if strings.Trim(v, " ") == "" {
				continue
			}
			args = append(args, v)
		}
	}
	return path, args
}

// 检查一个路径是文件还是文件夹 0文件不存在 1文件 2文件夹
func CheckFileOrDir(p string) (int, error) {
	info, err := os.Stat(p)
	if err != nil {
		return 0, err
	}
	if info.IsDir() {
		return 2, nil
	} else {
		return 1, nil
	}
}

// 把路径分隔符 "\" 统一为 "/"
func ReplcePathSeparator(p string) string {
	return strings.Replace(p, `\`, `/`, -1)
}

// 替换路径中的文件名
func ReplaceFileNameByPath(p, filename string) string {
	return filepath.Dir(p) + "/" + filename + path.Ext(p)
}

// 目录切分为数组
func SplitMenu(p string) []string {
	p = ReplcePathSeparator(p)
	menus := strings.Split(p, "/")
	arr := []string{}
	for _, v := range menus {
		if v != "" {
			arr = append(arr, v)
		}
	}
	return arr
}

// 替换路径中的特殊符号到全角
func ReplaceFileNameSpecial(opt, name string) string {
	if name == "" {
		return name
	}

	symbols := map[string]string{
		"\\": "＼",
		"/":  "／",
		":":  "：",
		"?":  "？",
		"<":  "《",
		">":  "》",
		"|":  "｜",
	}

	if opt == "encode" {
		for en, de := range symbols {
			name = strings.ReplaceAll(name, en, de)
		}
	} else {
		for en, de := range symbols {
			name = strings.ReplaceAll(name, de, en)
		}
	}
	return name
}
