package routes

import (
	"github.com/gin-gonic/gin"
	//"github.com/hakuuww/go-gin/middlewares"
)

func SetupRouter() *gin.Engine {
	server := gin.Default()

	RegisterAPIRoutes(server)

	return server
}