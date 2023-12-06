package handlers

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

// ChunkSize is the size of each chunk in bytes (25MB)
const ChunkSize = int64(25 * 1024 * 1024)

const channelID string = "1182079083249680404"

// FileUploadMiddleware is a middleware that breaks the uploaded file into chunks and stores them in Gin context.
func FileUploadMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Println(file.Filename)

		filePath := "uploads/" + filepath.Base(file.Filename)
		err = c.SaveUploadedFile(file, filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Open the uploaded file
		uploadedFile, err := os.Open(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer uploadedFile.Close()

		// Calculate the number of chunks
		fileInfo, _ := uploadedFile.Stat()
		fileSize := fileInfo.Size()
		numChunks := (fileSize + ChunkSize - 1) / ChunkSize

		// Create a buffer to store each chunk
		chunks := make([][]byte, numChunks)

		// Read the file into chunks
		for i := int64(0); i < numChunks; i++ {
			offset := i * ChunkSize
			chunkSize := ChunkSize
			if offset+ChunkSize > fileSize {
				chunkSize = fileSize - offset
			}
			chunk := make([]byte, chunkSize)
			_, err := uploadedFile.ReadAt(chunk, offset)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			chunks[i] = chunk
		}

		// Store chunks in Gin context
		c.Set("fileChunks", chunks)
		c.Next()
	}
}

// UploadHandler is the handler function that sends each chunk to a specific channel in Discord.
func UploadHandler(c *gin.Context, dg *discordgo.Session) {
	// Retrieve chunks from Gin context
	chunks, exists := c.Get("fileChunks")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File chunks not found"})
		return
	}

	// Convert chunks to bytes
	var buffer bytes.Buffer
	for _, chunk := range chunks.([][]byte) {
		buffer.Write(chunk)
	}

	// Send the file to Discord channel
	dg.ChannelFileSend(channelID, filepath.Base(c.Param("filename")), &buffer)

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded and sent to Discord"})
}

