package handlers

import "strings"

func parseRouteParamater(fullpath, pattern string) string {
	return strings.TrimLeft(fullpath, pattern)
}
