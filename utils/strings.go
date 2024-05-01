package utils

import (
	"strconv"
	"strings"
)

const (
	Teststr = "http://zhtcpic.storage.bbzhtc.com/pos/103/1/20220107164545288.jpg"
)

// 只取尾部
func GetTailbySlit(str, sep string) string {
	str = strings.ReplaceAll(strings.ReplaceAll(str, "\n", ""), " ", "")
	return strings.Split(str, sep)[len(strings.Split(str, sep))-1]
}

// 去掉头部
func CutHeadbySplit(str, sep string) string {
	// 使用逗号分隔字符串
	str = strings.ReplaceAll(strings.ReplaceAll(str, "\n", ""), " ", "")
	parts := strings.Split(str, sep)

	// 检查是否有足够的元素
	if len(parts) > 1 {
		// 删除第一个元素
		parts = parts[3:]
	} else {
		// 如果只有一个元素，则将切片设为空切片
		parts = []string{}
	}

	// 将结果合并为一个字符串
	result := strings.Join(parts, sep)

	return result
}
func StrToInt(str string) int {
	if i, err := strconv.Atoi(str); err != nil {
		return 0
	} else {
		return i
	}

}
