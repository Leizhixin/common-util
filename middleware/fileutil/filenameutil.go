package fileutil

import "strings"

func GetSuffixByName(fileName string) string {
	arr := strings.Split(fileName, ".")
	return arr[len(arr)-1]
}
