package handlers

import "net/http"

// GetTokenDetails This is a protected route. This will only work if you have a
// valid JWT token in the request headers, and if the token is still valid. You
// can use this route to see what is currently stored in a token
func GetTokenDetails(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"email": r.Context().Value("email").(string),
	}
	SendResponse(w, data, 200)
	return
}
