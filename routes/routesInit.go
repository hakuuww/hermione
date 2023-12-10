package routes

import (
	"github.com/gin-gonic/gin"
	//"github.com/hakuuww/go-gin/middlewares"
	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"

)

func SetupRouter(dg *discordgo.Session, fileList *mongo.Collection) *gin.Engine {
	server := gin.Default()

	RegisterAPIRoutes(server, dg, fileList)

	return server
}