package utils

import "github.com/google/uuid"

// 生成uuid
func CreateUUid() string {
	return uuid.New().String()
}
