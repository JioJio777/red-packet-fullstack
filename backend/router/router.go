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
			user.GET("/red-packets/sent", handler.GetSentRedPackets)
			user.GET("/red-packets/received", handler.GetReceivedRedPackets)
		}

		rp := api.Group("/red-packets").Use(middleware.Auth())
		{
			rp.POST("", handler.SendRedPacket)
			rp.POST("/:id/claim", handler.ClaimRedPacket)
			rp.GET("/:id", handler.GetRedPacketDetail)
			rp.GET("/:id/records", handler.GetRedPacketRecords)
		}
	}

	return r
}
