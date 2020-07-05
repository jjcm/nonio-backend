package handlers

import "net/http"

// GetTokenDetails This is a protected route. This will only work if you have a
// valid JWT token in the request headers, and if the token is still valid. You
// can use this route to see what is currently stored in a token
func GetTokenDetails(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"email":    r.Context().Value("user_email").(string),
		"username": r.Context().Value("user_username").(string),
		"name":     r.Context().Value("user_name").(string),
		"lastlog":  r.Context().Value("user_last_login").(string),
		"id":       r.Context().Value("user_id").(int),
	}
	SendResponse(w, data, 200)
	return
}
