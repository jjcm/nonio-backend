package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"soci-backend/httpd/utils"
	"soci-backend/models"
)

// ChannelCreate - POST /community/channel/create (moderators only)
func ChannelCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		SendResponse(w, utils.MakeError("you can only POST to this route"), 405)
		return
	}

	type requestPayload struct {
		Community string `json:"community"`
		Kind      string `json:"kind"`
		Slug      string `json:"slug"`
		Name      string `json:"name"`
	}
	var payload requestPayload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		SendResponse(w, utils.MakeError("invalid JSON"), 400)
		return
	}

	communityURL := strings.TrimSpace(strings.TrimPrefix(payload.Community, "@"))
	if communityURL == "" {
		SendResponse(w, utils.MakeError("community is required"), 400)
		return
	}

	c := models.Community{}
	if err := c.FindByURL(communityURL); err != nil {
		sendNotFound(w, errors.New("community not found"))
		return
	}

	userID := r.Context().Value("user_id").(int)
	mods, err := c.GetModerators()
	if err != nil {
		sendSystemError(w, err)
		return
	}
	isMod := false
	for _, mod := range mods {
		if mod.ID == userID {
			isMod = true
			break
		}
	}
	if !isMod {
		SendResponse(w, utils.MakeError("only moderators can create channels"), 403)
		return
	}

	u := models.User{}
	u.FindByID(userID)
	ch, err := u.CreateChannel(c.ID, payload.Kind, payload.Slug, payload.Name)
	if err != nil {
		SendResponse(w, map[string]string{"error": err.Error()}, 400)
		return
	}

	SendResponse(w, ch, 201)
}

// GetChannels - GET /community/channels?community=... (member-gated)
func GetChannels(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		SendResponse(w, utils.MakeError("you can only GET this route"), 405)
		return
	}

	communityURL := strings.TrimSpace(r.URL.Query().Get("community"))
	communityURL = strings.TrimPrefix(communityURL, "@")
	if communityURL == "" {
		SendResponse(w, utils.MakeError("community is required"), 400)
		return
	}

	c := models.Community{}
	if err := c.FindByURL(communityURL); err != nil {
		sendNotFound(w, errors.New("community not found"))
		return
	}

	userID := r.Context().Value("user_id").(int)
	u := models.User{}
	u.FindByID(userID)
	subs, err := u.GetSubscribedCommunities()
	if err != nil {
		sendSystemError(w, err)
		return
	}
	isMember := false
	for _, sub := range subs {
		if sub.URL == c.URL {
			isMember = true
			break
		}
	}
	if !isMember {
		SendResponse(w, utils.MakeError("you must be a member of this community to list channels"), 403)
		return
	}

	channels, err := models.GetChannelsByCommunityID(c.ID)
	if err != nil {
		sendSystemError(w, err)
		return
	}

	SendResponse(w, map[string]interface{}{"channels": channels}, 200)
}

// ChannelMessageCreate - POST /community/channel/message (send message to text channel)
func ChannelMessageCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		SendResponse(w, utils.MakeError("you can only POST to this route"), 405)
		return
	}

	type requestPayload struct {
		Community string `json:"community"`
		Channel   string `json:"channel"` // slug of the text channel
		Content   string `json:"content"`
		ImageURL  string `json:"imageUrl"`
	}
	var payload requestPayload
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&payload); err != nil {
		SendResponse(w, utils.MakeError("invalid JSON"), 400)
		return
	}

	communityURL := strings.TrimSpace(strings.TrimPrefix(payload.Community, "@"))
	channelSlug := strings.TrimSpace(payload.Channel)
	if communityURL == "" || channelSlug == "" {
		SendResponse(w, utils.MakeError("community and channel are required"), 400)
		return
	}

	c := models.Community{}
	if err := c.FindByURL(communityURL); err != nil {
		sendNotFound(w, errors.New("community not found"))
		return
	}

	ch := models.CommunityChannel{}
	if err := ch.FindByCommunityAndSlug(c.ID, channelSlug); err != nil {
		sendNotFound(w, errors.New("channel not found"))
		return
	}
	if ch.Kind != models.ChannelKindText {
		SendResponse(w, utils.MakeError("channel is not a text channel"), 400)
		return
	}

	userID := r.Context().Value("user_id").(int)
	u := models.User{}
	u.FindByID(userID)
	subs, err := u.GetSubscribedCommunities()
	if err != nil {
		sendSystemError(w, err)
		return
	}
	isMember := false
	for _, sub := range subs {
		if sub.URL == c.URL {
			isMember = true
			break
		}
	}
	if !isMember {
		SendResponse(w, utils.MakeError("you must be a member of this community to send messages"), 403)
		return
	}

	msg, err := u.CreateChannelMessage(ch.ID, payload.Content, payload.ImageURL)
	if err != nil {
		SendResponse(w, map[string]string{"error": err.Error()}, 400)
		return
	}

	SendResponse(w, msg, 201)
}

// GetChannelMessages - GET /community/channel/messages?community=...&channel=...&before=...&limit=...
func GetChannelMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		SendResponse(w, utils.MakeError("you can only GET this route"), 405)
		return
	}

	communityURL := strings.TrimSpace(r.URL.Query().Get("community"))
	communityURL = strings.TrimPrefix(communityURL, "@")
	channelSlug := strings.TrimSpace(r.URL.Query().Get("channel"))
	if communityURL == "" || channelSlug == "" {
		SendResponse(w, utils.MakeError("community and channel are required"), 400)
		return
	}

	c := models.Community{}
	if err := c.FindByURL(communityURL); err != nil {
		sendNotFound(w, errors.New("community not found"))
		return
	}

	ch := models.CommunityChannel{}
	if err := ch.FindByCommunityAndSlug(c.ID, channelSlug); err != nil {
		sendNotFound(w, errors.New("channel not found"))
		return
	}
	if ch.Kind != models.ChannelKindText {
		SendResponse(w, utils.MakeError("channel is not a text channel"), 400)
		return
	}

	userID := r.Context().Value("user_id").(int)
	u := models.User{}
	u.FindByID(userID)
	subs, err := u.GetSubscribedCommunities()
	if err != nil {
		sendSystemError(w, err)
		return
	}
	isMember := false
	for _, sub := range subs {
		if sub.URL == c.URL {
			isMember = true
			break
		}
	}
	if !isMember {
		SendResponse(w, utils.MakeError("you must be a member of this community to read messages"), 403)
		return
	}

	params := &models.ChannelMessageQueryParams{ChannelID: ch.ID, Limit: 50}
	if b := r.URL.Query().Get("before"); b != "" {
		if id, err := strconv.Atoi(b); err == nil && id > 0 {
			params.BeforeID = id
		}
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 {
			params.Limit = n
		}
	}

	msgs, err := models.GetChannelMessages(params)
	if err != nil {
		sendSystemError(w, err)
		return
	}

	SendResponse(w, map[string]interface{}{"messages": msgs}, 200)
}

// ChannelMessages - dispatch GET or POST for /community/channel/messages
func ChannelMessages(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetChannelMessages(w, r)
	case http.MethodPost:
		ChannelMessageCreate(w, r)
	default:
		SendResponse(w, utils.MakeError("method not allowed"), 405)
	}
}
