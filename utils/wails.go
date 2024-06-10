package utils

import (
	"strings"
)

const wailsPathPrefixHead = "/ASSET/"
const wailsPathPrefixEnd = ".asset"
const WailsPathPattern = `(/ASSET/.+?\.asset)` //正则

// wails 路径编码
func WailsPathEncode(p string, toAbs bool) string {
	if p == "" {
		return p
	}

	if strings.HasPrefix(p, wailsPathPrefixHead) {
		return p
	}

	if p == "" {
		return ""
	}
	if toAbs && !IsAbsPath(p) {
		p = ToAbsPath(p, "")
	}
	return wailsPathPrefixHead + Base64Encode(p) + wailsPathPrefixEnd
}

// wails 路径解码
func WailsPathDecode(p string) string {
	if p == "" {
		return p
	}
	if !strings.HasPrefix(p, wailsPathPrefixHead) {
		return p
	}

	p = strings.Replace(p, wailsPathPrefixHead, "", 1)
	p = strings.Replace(p, wailsPathPrefixEnd, "", 1)
	return Base64Decode(p)
}

// 检查是否为 wails路径
func WailsPathCheck(p string) bool {
	if p == "" {
		return false
	}
	return strings.Contains(p, wailsPathPrefixHead)
}
