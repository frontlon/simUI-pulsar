package components

import (
	"archive/zip"
	"io"
	"os"
	"simUI/constant"
	"simUI/utils"
)

var zipMethod = 1                           //0仅存储;1压缩
var TmpOutputIniPath = "./cache/config.ini" //ini文件路径

/**
 * 将资源文件添加到压缩文件中
 **/
func CompressZip(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		if prefix != "" {
			prefix = prefix + "/" + info.Name()
		} else {
			prefix = info.Name()
		}
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			defer f.Close()
			if err != nil {
				return err
			}
			err = CompressZip(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		if prefix != "" {
			header.Name = prefix + "/" + header.Name
		}
		if zipMethod == 1 {
			header.Method = zip.Deflate //压缩
		} else {
			header.Method = zip.Store //仅存储
		}

		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(writer, file)
		if err != nil {
			return err
		}
	}
	return nil
}

// 检查分享文件
func CheckZip(zipFile string) (bool, error) {
	exists := false
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return false, err
	}
	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			continue
		}
		defer rc.Close()

		f := utils.GetFileNameAndExt(file.Name)
		if f == constant.SHARE_FILE_NAME {
			exists = true
			break
		}
	}
	return exists, nil
}

// 分享文件解压
func DecompressZip(zipFile, rootPath, romPath string) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()
	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		m := utils.GetFilePath(file.Name)
		f := rootPath + "/" + file.Name
		r := rootPath
		if m == constant.RES_DIR["rom"] {
			f = romPath + "/" + file.Name
			r = romPath
		}
		err = os.MkdirAll(r, 0755)
		if err != nil {
			return err
		}
		w, err := os.Create(f)
		if err != nil {
			return err
		}
		defer w.Close()
		_, err = io.Copy(w, rc)
		if err != nil {
			return err
		}
		w.Close()
		rc.Close()
	}
	return nil
}
