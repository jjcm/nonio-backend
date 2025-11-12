package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"soci-backend/httpd/utils"
	"soci-backend/models"
)

// PostEncodingComplete - marks a post as no longer encoding
// This endpoint is called by the video CDN when encoding is complete
func PostEncodingComplete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		SendResponse(w, utils.MakeError("you can only POST to the encoding complete route"), 405)
		return
	}

	type requestPayload struct {
		URL string `json:"url"`
	}

	var payload requestPayload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		SendResponse(w, utils.MakeError("invalid request payload"), 400)
		return
	}

	if strings.TrimSpace(payload.URL) == "" {
		SendResponse(w, utils.MakeError("url is required"), 400)
		return
	}

	p := models.Post{}
	if err := p.MarkEncodingComplete(payload.URL); err != nil {
		Log.Error("Failed to mark encoding complete")
		sendSystemError(w, err)
		return
	}

	SendResponse(w, map[string]string{"status": "success"}, 200)
}

