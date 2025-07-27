package main

import (
	"log"
	"messenger/database"
	"messenger/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	hub := handlers.NewHub()
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	raw, err := database.OpenUsers()
	if err != nil {
		log.Fatalf("an error occured during opening database %v", err)
	}

	db := handlers.NewDB(raw, hub)
	routs(router, db)

	router.Run(":8080")
}
