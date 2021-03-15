package util

import (
	"os"
	"regexp"
	"strings"
)

/**
 * 获取系统参数 arg=value
 */
func GetArg(name, defaultValue string) string {
	args := os.Args
	pattern := name + "="
	for _, arg := range args {
		if match, _ := regexp.MatchString(pattern, arg); match {
			value := strings.Split(arg, "=")
			return strings.TrimSpace(value[1])
		}
	}
	return strings.TrimSpace(defaultValue)
}
