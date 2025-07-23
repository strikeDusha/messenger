package main

import (
	"messenger/handlers"
	"messenger/middleware"

	"github.com/gin-gonic/gin"
)

func routs(r *gin.Engine, db *handlers.DB) {
	r.Static("/frontend", "./frontend/")
	{
		api := r.Group("api")
		api.GET("/me", middleware.IsAuth, db.IsLoggedIn)
		api.POST("/login", db.HandleLogin)
		api.POST("/register", db.HandleRegister)
		{
			chat := api.Group("chat", middleware.IsAuth)
			chat.GET("/", handlers.ChatApp) //uses mongo later
		}
	}
	r.NoRoute(func(c *gin.Context) {
		c.File("./frontend/index.html")
	})

}
