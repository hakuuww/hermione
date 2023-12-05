package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func uploadFile(c *gin.Context) {



	c.JSON(http.StatusOK, gin.H{"upload": "success"})
}

