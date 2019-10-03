package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jjcm/soci-backend/httpd/handlers"
)

// CheckToken this acts as a middleware, but I'm not really using any middleware packages
func CheckToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			handlers.CorsAdjustments(&w)
			handlers.SendResponse(w, "", 200)
			return
		}

		token := r.Header.Get("Authorization")
		if strings.TrimSpace(token) == "" || strings.TrimSpace(token) == "Bearer" {
			handlers.SendResponse(w, "Authorization required", http.StatusUnauthorized)
			return
		}

		// if the header starts with "Bearer " then let's trim that junk
		if token[:7] == "Bearer " {
			token = token[7:]
		}

		goodies, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return handlers.HmacSampleSecret, nil
		})

		claims, ok := goodies.Claims.(jwt.MapClaims)
		if !ok || !goodies.Valid || err != nil {
			handlers.SendResponse(w, handlers.MakeError("Error working with your token"), 500)
			return
		}

		// check expiresAt inside the token
		// convert unix timestamp to int
		i := int64(claims["expiresAt"].(float64))
		ts := time.Unix(i, 0)
		now := time.Now()
		if now.After(ts) {
			handlers.SendResponse(w, handlers.MakeError("Your token is expired"), http.StatusUnauthorized)
			return
		}
		secondsRemaining := int(ts.Sub(now).Seconds())
		w.Header().Set("X-Seconds-Remaining", strconv.Itoa(secondsRemaining))

		ctx := context.WithValue(r.Context(), "email", claims["email"])
		next(w, r.WithContext(ctx))
	}
}
