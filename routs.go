package main

import (
	"messenger/handlers"
	"messenger/middleware"

	"github.com/gin-gonic/gin"
)

func routs(r *gin.Engine, db *handlers.DB) {
	go db.H.Run()

	{
		api := r.Group("api")
		api.GET("/me", middleware.IsAuth, db.IsLoggedIn)
		api.POST("/login", db.HandleLogin)
		api.POST("/register", db.HandleRegister)
		{
			chat := api.Group("chat", middleware.IsAuth)
			chat.GET("/ws", db.ServeWs) //websocket connection\
			chat.GET("/history/:id", db.History)
		}
	}
	r.StaticFile("/test-chat", "./frontend/test-chat.html")
	r.Static("/frontend", "./frontend/")
	r.NoRoute(func(c *gin.Context) {
		c.File("frontend/index.html")
	})

}
