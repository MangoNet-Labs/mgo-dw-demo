package middleware

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"user/common/response"
)

type AuthUser struct {
	UserID     int64  `json:"user_id"`
	MgoAddress string `json:"mgo_address"`
	SolAddress string `json:"sol_address"`
}

func NewJwtMiddleware(secret string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			tokenStr := r.Header.Get("Token")
			if tokenStr == "" {
				response.FailJson(w, "missing Token header", 401)
				return
			}

			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				response.FailJson(w, "invalid token", 401)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				userId, ok := claims["userId"].(float64)
				if !ok {
					response.FailJson(w, "invalid token claims", 401)
					return
				}
				mgoAddress, _ := claims["mgoAddress"].(string)
				solAddress, _ := claims["solAddress"].(string)
				ctx := context.WithValue(r.Context(), "authUser", &AuthUser{
					UserID:     int64(userId),
					MgoAddress: mgoAddress,
					SolAddress: solAddress,
				})
				next(w, r.WithContext(ctx))
				return
			}
			response.FailJson(w, "invalid token structure", 401)
		}
	}
}

func GetAuthUser(ctx context.Context) (*AuthUser, bool) {
	user, ok := ctx.Value("authUser").(*AuthUser)
	return user, ok
}
