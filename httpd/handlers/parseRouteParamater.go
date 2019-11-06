package handlers

import "strings"

func parseRouteParamater(fullpath, pattern string) string {
	return strings.TrimPrefix(fullpath, pattern)
}
