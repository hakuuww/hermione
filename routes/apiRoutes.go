package routes

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/hakuuww/hermione/handlers"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterAPIRoutes(router *gin.Engine, dg *discordgo.Session, fileList *mongo.Collection) {
	router.POST("/upload",handlers.FileUploadMiddleware(fileList), func(c *gin.Context) {
		handlers.UploadHandler(c,dg,fileList)
	})
}
