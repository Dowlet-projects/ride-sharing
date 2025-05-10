// // handlers/middleware.go
// // Authentication middleware for JWT

// package handlers

// import (
// 	"context"
// 	"net/http"
// 	"strings"

// 	"ride-sharing/models"
// 	"ride-sharing/utils"

// 	"github.com/dgrijalva/jwt-go"
// )

// // authMiddleware authenticates JWT tokens
// func (app *App) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		tokenStr := r.Header.Get("Authorization")
// 		if tokenStr == "" {
// 			utils.RespondError(w, http.StatusUnauthorized, "Authorization header missing")
// 			return
// 		}
// 		// Remove "Bearer " prefix
// 		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

// 		claims := &models.Claims{}
// 		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
// 			return []byte(app.Config.JWTSecret), nil
// 		})
// 		if err != nil || !token.Valid {
// 			utils.RespondError(w, http.StatusUnauthorized, "Invalid token")
// 			return
// 		}

// 		// Add claims to context
// 		ctx := context.WithValue(r.Context(), "claims", claims)
// 		next(w, r.WithContext(ctx))
// 	}
// }


package handlers

import (
	"context"
	"net/http"
	"strings"

	"ride-sharing/models"
	"ride-sharing/utils"

	"github.com/dgrijalva/jwt-go"
)

// authMiddleware authenticates JWT tokens
func (app *App) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			utils.RespondError(w, http.StatusUnauthorized, "Authorization header missing")
			return
		}
		// Remove "Bearer " prefix
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		claims := &models.Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(app.Config.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			utils.RespondError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}