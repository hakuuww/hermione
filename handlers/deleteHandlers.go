package handlers

import (
	"context"
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// DeleteFile deletes a document from the MongoDB collection based on the filename parameter.
func DeleteFile(c *gin.Context, fileList *mongo.Collection) {
	// Get the filename parameter from the path
	filename := c.Param("filename")

	filename = strings.TrimLeft(filename, "/")

	// Define the filter to find the document by filename
	filter := bson.M{"fileName": filename}

	// Attempt to delete the document from the collection
	result, err := fileList.DeleteOne(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting document"})
		return
	}

	// Check if the document was found and deleted
	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Document not found"})
		return
	}

	// Document successfully deleted
	c.JSON(http.StatusOK, gin.H{"message": "Document deleted successfully"})
}