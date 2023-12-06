package routes

import (
	"github.com/gin-gonic/gin"
	//"github.com/hakuuww/go-gin/middlewares"
	"github.com/bwmarrin/discordgo"

)

func SetupRouter(dg *discordgo.Session) *gin.Engine {
	server := gin.Default()

	RegisterAPIRoutes(server, dg)

	return server
}