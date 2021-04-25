package handlers

import (
	"net/http"
	"soci-backend/httpd/utils"
	"soci-backend/models"
	"strings"
)

// GetTags - get tags out of the database, 100 at a time, optional offset
func GetTags(w http.ResponseWriter, r *http.Request) {
	tags, err := models.GetTags(0, 100)
	if err != nil {
		SendResponse(w, utils.MakeError(err.Error()), 500)
		return
	}

	output := map[string]interface{}{
		"tags": tags,
	}
	SendResponse(w, output, 200)
}

// GetTagsByPrefix - search for tags beginning with a string
func GetTagsByPrefix(w http.ResponseWriter, r *http.Request) {
	prefix := strings.TrimSpace(utils.ParseRouteParameter(r.URL.Path, "/tags/"))
	if prefix == "" {
		GetTags(w, r)
		return
	}

	tags, err := models.GetTagsByPrefix(prefix)
	if err != nil {
		SendResponse(w, utils.MakeError(err.Error()), 500)
		return
	}

	output := map[string]interface{}{
		"tags": tags,
	}
	SendResponse(w, output, 200)
}
