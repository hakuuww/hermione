package routes

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/hakuuww/hermione/handlers"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterAPIRoutes(router *gin.Engine, dg *discordgo.Session, fileList *mongo.Collection) {
	router.POST("/upload", handlers.FileUploadMiddleware(fileList), func(c *gin.Context) {
		handlers.UploadHandler(c, dg, fileList)
	})

	// Download route
	router.GET("/download/:filename", handlers.DownloadMiddleware1(fileList), func(c *gin.Context) {
		handlers.DownloadHandler(c, dg)
	})

	router.GET("/allFileInfo", func(c *gin.Context) {
		handlers.DisplayFiles(c, fileList)
	})

	router.GET("/search/:filename", func(c *gin.Context) {
		handlers.SearchFile(c, fileList)
	})

	router.DELETE("/delete/:filename", func(c *gin.Context) {
		handlers.DeleteFile(c, fileList)
	})

}
