package utils

import (
	"fmt"
	"strings"
)

func MergeConnInfos(connInfos ...string) string {
	connInfo := make(map[string]string)

	for _, info := range connInfos {
		fields := strings.Split(info, " ")
		for _, field := range fields {
			keyValue := strings.Split(field, "=")
			connInfo[keyValue[0]] = keyValue[1]
		}
	}

	var ret []string
	for key, value := range connInfo {
		ret = append(ret, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(ret, " ")
}
