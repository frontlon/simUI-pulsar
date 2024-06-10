package utils

import (
	"errors"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"os"
	"strings"
)

func CreateThumbnail(src string, dst string, imgWidth int, quality int) error {

	format := strings.ToLower(GetFileExt(src))

	//检查图片格式是否支持
	allowFormats := map[string]uint8{".gif": 1, ".jpg": 1, ".jpeg": 1, ".png": 1}
	if _, ok := allowFormats[format]; !ok {
		return nil
	}

	if !FileExists(src) {
		return errors.New("file not exists")
	}

	//读取文件
	fa, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fa.Close()

	config, _, err := image.DecodeConfig(fa)
	if err != nil {
		return err
	}

	if config.Width < imgWidth {
		return errors.New("no need dst scale")
	}

	fb, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fb.Close()
	origin, _, err := image.Decode(fb)
	if err != nil {
		return err
	}

	// 做等比缩放
	width := uint(imgWidth)
	height := uint(imgWidth * config.Height / config.Width)

	canvas := resize.Thumbnail(width, height, origin, resize.Lanczos3)

	if FileExists(dst) {
		FileDelete(dst)
	}

	file_out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer file_out.Close()

	//生成图片
	err = jpeg.Encode(file_out, canvas, &jpeg.Options{quality})
	if err != nil {
		return err
	}

	return nil
}
