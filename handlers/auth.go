package handlers

import (
	"messenger/database"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// переименовать все err для фронтенда

type DB struct {
	*database.Storage
	//mongo for chat
}

func (s *DB) HandleRegister(c *gin.Context) {
	var user database.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{
			"error": "not binding",
		})
		return
	}
	err := s.Storage.AddUser(&user)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "database",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Username,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // expires after week
	})
	tokenStr, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "while creating tokn",
		})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", tokenStr, 3600*24*7, "", "", false, true)
	c.JSON(200, gin.H{})

}

func (s *DB) HandleLogin(c *gin.Context) {

	var login struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(400, gin.H{
			"error": "not binding",
		})
		return
	}
	user, err := s.GetUser2Auth(login.Email, login.Password)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "db err",
		})
		return
	}

	if user == nil {
		c.JSON(400, gin.H{
			"error": "no such user is exists",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Username,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // expires after week
	})
	tokenStr, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(400, gin.H{
			"error": "while creating tokn",
		})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("token", tokenStr, 3600*24*7, "", "", false, true)
	c.JSON(200, gin.H{})
}
