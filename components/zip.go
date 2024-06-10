package components

import (
	"archive/zip"
	"github.com/axgle/mahonia"
	"io"
	"os"
	"path/filepath"
	"simUI/constant"
	"simUI/utils"
	"strings"
)

/*
 zip解压
*/

func UnzipRom(zipFile string) (string, error) {

	if strings.ToLower(filepath.Ext(zipFile)) != ".zip" {
		return "", nil
	}

	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return "", err
	}
	defer zipReader.Close()

	//拼接解压路径
	zipfileName := utils.GetFileName(zipFile)
	fpath := constant.CACHE_UNZIP_PATH + zipfileName + "/"

	if !utils.IsDirEmpty(fpath) {
		return fpath, nil
	}

	if !utils.DirExists(fpath) {
		if err = utils.CreateDir(fpath); err != nil {
		}
	}

	for _, f := range zipReader.File {
		//解决中文文件名乱码问题
		enc := mahonia.NewDecoder("gbk")
		f.Name = enc.ConvertString(f.Name)

		//开始解压
		if f.FileInfo().IsDir() {
			if err = utils.CreateDir(fpath + f.Name); err != nil {
				return "", err
			}
		} else {
			srcFile, err := f.Open()
			if err != nil {
				return "", err
			}
			defer srcFile.Close()

			dstPath := fpath + f.Name
			dstFile, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return "", err
			}
			defer dstFile.Close()

			_, err = io.Copy(dstFile, srcFile)
			if err != nil {
				return "", err
			}

		}
	}
	return fpath, nil
}

/*
清理解压缓存
*/
func ClearZipRom() error {
	err := os.RemoveAll(constant.CACHE_UNZIP_PATH)
	if err != nil {
		return err
	}
	return nil
}
