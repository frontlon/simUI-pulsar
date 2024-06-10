package utils

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
)

// base64编码
func Base64Encode(s string) string {
	if s == "" {
		return ""
	}
	encodeString := base64.StdEncoding.EncodeToString([]byte(s))
	return strings.Replace(encodeString, "=", "_", -1)
}

// base64解码
func Base64Decode(s string) string {
	if s == "" {
		return ""
	}

	s = strings.Replace(s, "_", "=", -1)
	decodeBytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return ""
	}
	return string(decodeBytes)
}

// base64转换为文件
func Base64ToFile(base64Str, outputPath string) error {
	// 提取base64数据部分
	data := base64Str
	if idx := strings.Index(base64Str, "base64,"); idx != -1 {
		data = base64Str[idx+7:]
		fmt.Println("=========", idx, idx+7)

	}

	// 解码base64内容
	decodedData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return fmt.Errorf("decode base64 string: %v", err)
	}

	// 写入文件
	err = os.WriteFile(outputPath, decodedData, 0666)
	if err != nil {
		return fmt.Errorf("write file: %v", err)
	}

	return nil
}
