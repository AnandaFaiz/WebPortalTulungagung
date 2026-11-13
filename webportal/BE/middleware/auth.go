package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// 1. Coba ambil dari Cookie
		cookie, err := c.Cookie("auth_token")
		if err == nil {
			tokenString = cookie
		}

		// 2. Jika tidak ada, coba ambil dari Header Authorization
		if tokenString == "" {
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"statusMessage": "Belum login"})
			c.Abort()
			return
		}

		// 3. Verifikasi Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"statusMessage": "Token tidak valid atau kadaluarsa"})
			c.Abort()
			return
		}

		// 4. Simpan data user ke Context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", int(claims["id"].(float64)))
			c.Set("user_role", claims["role"].(string))
		}

		c.Next()
	}
}