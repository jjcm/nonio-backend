package handlers

import "strings"

func parseRouteParameter(fullpath, pattern string) string {
	return strings.TrimPrefix(fullpath, pattern)
}
