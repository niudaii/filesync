package filesync

import (
	"strings"
)

// syncFileBlackList 不需要、禁止同步的文件黑名单
var syncFileBlackList []string

// checkFileIsSyncList 同步文件的白名单校验
func checkFileIsSyncList(filePathName string) bool {
	if strings.Contains(filePathName, "..") {
		return false
	}
	for _, f := range syncFileBlackList {
		if strings.HasPrefix(filePathName, f) {
			return false
		}
	}
	return true
}
