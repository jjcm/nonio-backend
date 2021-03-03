package handlers

import (
	"net/http"
	"soci-backend/models"
)

// GetTokenDetails This is a protected route. This will only work if you have a
// valid JWT token in the request headers, and if the token is still valid. You
// can use this route to see what is currently stored in a token
func GetTokenDetails(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"email":    r.Context().Value("user_email").(string),
		"username": r.Context().Value("user_username").(string),
		"id":       r.Context().Value("user_id").(int),
	}
	SendResponse(w, data, 200)
}

// GetFinancials gets the current cash and subscription amount for a user.
func GetFinancials(w http.ResponseWriter, r *http.Request) {
	// get the user from context
	user := models.User{}
	user.FindByID(r.Context().Value("user_id").(int))

	financialData, err := user.GetFinancialData()
	if err != nil {
		sendSystemError(w, err)
		return
	}

	SendResponse(w, &financialData, 200)
}
