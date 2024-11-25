package middleware

import (
	"a21hc3NpZ25tZW50/model"
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Cookie("session_token")
		if err != nil || cookie == "" {
			if strings.HasPrefix(ctx.Request.URL.Path, "/api") || ctx.GetHeader("Accept") == "application/json" {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"error": "unauthorized: missing session token",
				})
				return
			}

			ctx.Redirect(http.StatusSeeOther, "/client/login")
			ctx.Abort()
			return
		}

		tokenClaims := model.Claims{}
		token, err := jwt.ParseWithClaims(cookie, &tokenClaims, func(token *jwt.Token) (interface{}, error) {
			return model.JwtKey, nil
		})
		
		if err != nil || !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "unauthorized: invalid token",
			})
			return
		}

		if tokenClaims.Email == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized: invalid email",
			})
			return
		}

		ctx.Set("email", tokenClaims.Email)
		ctx.Next()
	}
}