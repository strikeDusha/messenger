package middleware

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func IsAuth(c *gin.Context) {
	cookie, err := c.Cookie("token")
	if err != nil {
		c.AbortWithStatus(401)
		return
	}
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("SECRET")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		c.AbortWithStatus(401)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["exp"].(float64) < float64(time.Now().Unix()) {
			c.AbortWithStatus(401)
			return
		} else {
			c.Set("username", claims["sub"])
			c.Next()
		}
	} else {
		c.AbortWithStatus(401)
		return
	}

}
