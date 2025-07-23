package handlers

import "github.com/gin-gonic/gin"

func (s *DB) IsLoggedIn(c *gin.Context) {
	un, _ := c.Get("username")
	usr, err := s.Storage.GetByUsername(un.(string))
	if err != nil || usr == nil {
		c.JSON(404, gin.H{"error": "user not found"})
		return
	}

	c.JSON(200, gin.H{
		"pn": usr.PublicName,
		"un": usr.Username,
	})
}
