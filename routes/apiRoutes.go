package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hakuuww/hermione/handlers"
)

func RegisterAPIRoutes(router *gin.Engine) {
	router.POST("/upload",handlers.FileUploadMiddleware(), handlers.UploadHandler )
}
