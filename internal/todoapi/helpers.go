package todoapi

import (
	"strconv"
	"strings"
)

func getId(route, path string) *int64 {
	param := strings.TrimPrefix(path, route)
	if len(param) > 0 {
		id, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			return nil
		}
		return &id
	}
	return nil
}
