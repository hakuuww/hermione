package handlers

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/hakuuww/hermione/database"
	"github.com/hakuuww/hermione/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"time"
)

// ChunkSize is the size of each chunk in bytes (25MB)
const ChunkSize = int64(5 * 1024 * 1024)

// Define a rate limiter with a limit of 5 requests per second.
var uploadLimiter = rate.NewLimiter(rate.Every(time.Second)/2, 1)

const channelID string = "1182079083249680404"

// FileUploadMiddleware is a middleware that breaks the uploaded file into chunks and stores them in Gin context.
func FileUploadMiddleware(fileList *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {

		var w http.ResponseWriter = c.Writer
		c.Request.Body = http.MaxBytesReader(w, c.Request.Body, 1000<<40)

		file, header, err := c.Request.FormFile("file")
		log.Println(file)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			log.Println("cannot get formfile")
			return
		}
		defer file.Close()

		// Get the file size
		fileSize := header.Size
		fmt.Println("file size: %d", fileSize)
		fileName := header.Filename
		log.Println("file name: %s", fileName)

		// Extract file extension (postfix)
		fileExt := filepath.Ext(fileName)
		fmt.Printf("file extension: %s\n", fileExt)

		// Calculate the number of chunks
		numChunks := (fileSize + ChunkSize - 1) / ChunkSize

		ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
		defer cancel()

		newId := primitive.NewObjectID()

		// Example document with empty fileChunks
		docWithEmptyChunks := models.Document{
			ID:            newId,
			FileName:      fileName,
			FileType:      fileExt,
			NumOfChunks:   numChunks,
			FileSizeBytes: fileSize,
			FileChunks:    []models.FileChunk{}, // Create an empty slice of the struct type,
		}

		if err = database.InsertEmptyFileChunksDocument(ctx, fileList, docWithEmptyChunks) ; err!=nil {
			log.Println("cannot insert empty chunks document", err)
			c.Next()
			return
		}

		/*
			If you later use chunks in your program or return it from your function,
			the memory won't be freed until there are no more references to chunks.
			If, however, chunks goes out of scope and there are no other references to it,
			the memory occupied by the underlying array and the individual slices should be eligible for garbage collection.

			If your variable is stored in the Gin context, it becomes part of the heap and will be considered reachable during the mark phase if there are references to it.
			If there are no references to the variable outside the Gin context, and the context itself becomes unreachable (e.g., the HTTP request is completed),
			the variable should eventually become eligible for garbage collection.
		*/
		// Create a buffer to store each chunk
		chunks := make([][]byte, numChunks)

		// Read the file into chunks
		for i := int64(0); i < numChunks; i++ {
			chunkSize := ChunkSize
			if i == numChunks-1 {
				// Last chunk might be smaller
				chunkSize = fileSize - i*ChunkSize
			}

			chunk := make([]byte, chunkSize)
			_, err := file.Read(chunk)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			chunks[i] = chunk
		}

		// Store chunks in Gin context
		c.Set("fileChunks", chunks)
		//store number of chunks in Gin context
		c.Set("numChunks", numChunks)
		c.Set("fileName", fileName)
		c.Set("docID", newId)
		c.Next()
	}
}

// UploadHandler is the handler function that sends each chunk to a specific channel in Discord.
func UploadHandler(c *gin.Context, dg *discordgo.Session, fileList *mongo.Collection) {
	var wg sync.WaitGroup
	var mtx sync.Mutex

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	// Retrieve chunks from Gin context
	docID, exists := c.Get("docID")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Doc id not found"})
		return
	}

	// Retrieve chunks from Gin context
	chunks, exists := c.Get("fileChunks")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File chunks not found"})
		return
	}

	numChunks_fromCTX, exists := c.Get("numChunks")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File chunks not found"})
		return
	}

	var numChunks_sent int = 0
	var chunk_seq_num int = 0

	// Convert chunks to bytes
	for _, chunk := range chunks.([][]byte) {
		chunkSize := len(chunk)
		fmt.Printf("Chunk size: %d bytes\n", chunkSize)
		// Rate limit the ChannelFileSend operation
		if err := uploadLimiter.Wait(c.Request.Context()); err != nil {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			return
		}

		chunk_seq_num++

		wg.Add(1)

		go func(chunk_copy []byte, wg *sync.WaitGroup, mtx *sync.Mutex, chunk_seq_num int) {
			defer wg.Done()
			// Send the file to Discord channel
			message, err := dg.ChannelFileSend(channelID, filepath.Base(c.Param("filename")), bytes.NewReader(chunk_copy))

			//loop until successfully sending the file chunk
			for err != nil {
				// Rate limit the ChannelFileSend operation
				if err = uploadLimiter.Wait(c.Request.Context()); err != nil {
					c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
					return
				}

				message, err = dg.ChannelFileSend(channelID, filepath.Base(c.Param("filename")), bytes.NewReader(chunk_copy))
				if err != nil {
					fmt.Println(err.Error())
					fmt.Println(message)
				}

			}

			// Create a FileChunk struct
			fileChunk := models.FileChunk{
				MessageID:      message.ID,
				ChannelID:      channelID,
				SequenceNumber: chunk_seq_num,          // Provide the actual sequence number
				ChunkSizeBytes: int64(len(chunk_copy)), // Use the appropriate chunk size
			}

			database.AddFileChunkToDocument(ctx, fileList, docID.(primitive.ObjectID), fileChunk)
			mtx.Lock()
			defer mtx.Unlock()
			fmt.Println("-------------------------------")
			numChunks_sent++
			fmt.Println("numChunks_sent%d", numChunks_sent)
			fmt.Println("chunk_seq_num:%d", chunk_seq_num)
			fmt.Println("-------------------------------")

		}(chunk, &wg, &mtx, chunk_seq_num)

	}

	wg.Wait()

	fmt.Println("Number of chunks: %d", numChunks_sent)
	fmt.Println("Number of chunks from ctx: %d", numChunks_fromCTX)

	if int(numChunks_sent) != int(numChunks_fromCTX.(int64)) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in breaking up file"})
		return
	} else {
		fmt.Println("Number of chunks sent: %d", numChunks_sent)
	}

	c.JSON(http.StatusOK, gin.H{"message": "File uploaded and sent to Discord"})
}
