package routes

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/hakuuww/hermione/handlers"
)

func RegisterAPIRoutes(router *gin.Engine, dg *discordgo.Session) {
	router.POST("/upload",handlers.FileUploadMiddleware(), func(c *gin.Context) {
		handlers.UploadHandler(c,dg)
	})
}
