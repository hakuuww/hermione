package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hakuuww/discordFileServer/handlers"
)

func RegisterAPIRoutes(router *gin.Engine) {
	router.POST("/upload", handlers.uploadFile)
}
