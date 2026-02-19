package router

import (
	"red-packet/handler"
	"red-packet/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", handler.Register)
			auth.POST("/login", handler.Login)
		}

		user := api.Group("/user").Use(middleware.Auth())
		{
			user.GET("/profile", handler.GetProfile)
		}
	}

	return r
}
