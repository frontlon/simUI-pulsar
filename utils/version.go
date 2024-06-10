package utils

import (
	"fmt"
	"strings"
)

// 比较版本号
func CompareVersion(v1, v2 string) int {
	ver1 := strings.Split(v1, ".") // 将版本号按照"."分割为切片
	ver2 := strings.Split(v2, ".")

	for i := 0; i < len(ver1) || i < len(ver2); i++ {
		num1 := 0
		num2 := 0

		if i < len(ver1) {
			fmt.Sscanf(ver1[i], "%d", &num1)
		}

		if i < len(ver2) {
			fmt.Sscanf(ver2[i], "%d", &num2)
		}

		if num1 > num2 {
			return 1
		} else if num1 < num2 {
			return -1
		}
	}

	return 0
}
