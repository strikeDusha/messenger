package main

import (
	"log"
	"messenger/database"
	"messenger/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	raw, err := database.OpenUsers()
	if err != nil {
		log.Fatalf("an error occured during opening database %v", err)
	}
	db := &handlers.DB{Storage: raw}
	routs(router, db)
	router.Run(":8080")
}
