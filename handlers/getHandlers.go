package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"regexp"

	//"io/ioutil"
	//"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/hakuuww/hermione/database"
	"github.com/hakuuww/hermione/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/time/rate"
)

var downloadLimiter = rate.NewLimiter(rate.Every(time.Second)/2, 1)

// DownloadMiddleware1 is a middleware that retrieves the file chunks from the database based on the provided filename.
func DownloadMiddleware1(fileList *mongo.Collection) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get filename from path parameter
		filename := c.Param("filename")
		filename = strings.TrimLeft(filename, "/")

		fmt.Println(filename)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Query the database to get the document that matches the filename
		documentObject, err := database.GetDocumentByFilename(ctx, fileList, filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving file chunks from the database"})
			return
		}

		docID := documentObject.ID
		fileChunks := documentObject.FileChunks
		numOfChunks := documentObject.NumOfChunks

		// Set document and fileChunks and docID in Gin context for later use in the handler
		c.Set("document", documentObject)
		c.Set("fileChunks", fileChunks)
		c.Set("docID", docID)
		c.Set("numOfChunks", numOfChunks)

		// Call the next middleware or handler
		c.Next()
	}
}

// DownloadHandler is the handler function that downloads individual files from Discord and combines them into a single file.
func DownloadHandler(c *gin.Context, dg *discordgo.Session) {
	var wg sync.WaitGroup
	var mtx sync.Mutex

	// ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	// defer cancel()

	// Retrieve fileChunks and docID from Gin context
	fileChunks, exists := c.Get("fileChunks")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File chunks not found"})
		return
	}

	// Retrieve fileChunks and docID from Gin context
	numOfChunks, exists := c.Get("numOfChunks")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "numOfChunks not found"})
		return
	}

	fileChunksInfo := fileChunks.([]models.FileChunk)
	numOfChunksCasted := numOfChunks.(int64)

	// Create an array of empty ChunkData
	chunkDataArray := make([]*models.ChunkData, numOfChunksCasted)

	// Download individual files from Discord and combine them
	for _, chunk := range fileChunksInfo {
		wg.Add(1)

		go func(chunk models.FileChunk, wg *sync.WaitGroup, mtx *sync.Mutex, dg *discordgo.Session, chunkDataArray []*models.ChunkData) error {
			defer wg.Done()

			readChunkData, err := ReadAttachment(dg, chunk.ChannelID, chunk.MessageID, chunk.SequenceNumber)
			if err != nil {
				return err
			}

			mtx.Lock()
			chunkDataArray[readChunkData.Seq-1] = readChunkData
			mtx.Unlock()
			return nil
		}(chunk, &wg, &mtx, dg, chunkDataArray[:])

	}

	wg.Wait()

	combinedFileBuf, err := combineChunks(chunkDataArray)
	if err != nil {
		fmt.Println("Error combining chunks:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot combine buffer"})
		return
	}

	// Set response headers
	// Set the content disposition header to trigger a download
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(c.Param("filename"))))
	c.Data(http.StatusOK, "application/octet-stream", combinedFileBuf.Bytes())
}

// ReadAttachment reads the first attachment of a Discord message and returns it as a bytes.Buffer.
func ReadAttachment(session *discordgo.Session, channelID string, messageID string, seq_num int) (*models.ChunkData, error) {
	// Retrieve the message
	message, err := session.ChannelMessage(channelID, messageID)
	if err != nil {
		return nil, err
	}

	// Check if there are any attachments
	if len(message.Attachments) == 0 {
		return nil, fmt.Errorf("no attachments found in the message")
	}

	// Get the first attachment
	attachment := message.Attachments[0]
	url := attachment.URL

	buffer, err := downloadAttachmentToBuffer(url)
	if err != nil {
		return nil, err
	}

	// Create and populate ChunkData
	chunkData := &models.ChunkData{
		Buf:  buffer,
		Seq:  seq_num,             // Set the appropriate sequence number
		Size: int64(buffer.Len()), // Set the size of the buffer
	}

	return chunkData, nil
}

func downloadAttachmentToBuffer(url string) (*bytes.Buffer, error) {
	// Create a new HTTP client with a timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make a GET request to the attachment URL
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body (attachment content) into a bytes buffer
	buffer, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(buffer), nil
}

func combineChunks(chunkDataArray []*models.ChunkData) (*bytes.Buffer, error) {
	var combinedBuffer bytes.Buffer

	for _, chunk := range chunkDataArray {
		_, err := combinedBuffer.Write(chunk.Buf.Bytes())
		if err != nil {
			return nil, err
		}
	}
	return &combinedBuffer, nil
}

func DisplayFiles(c *gin.Context, fileList *mongo.Collection) {
	ctx := context.TODO()

	// Fetch all documents from the collection
	cursor, err := fileList.Find(ctx, bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving documents"})
		return
	}
	defer cursor.Close(ctx)

	// Iterate through the cursor and store documents in a slice
	var documents []models.ResponseDocument
	for cursor.Next(ctx) {
		var doc models.Document
		if err := cursor.Decode(&doc); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding document"})
			return
		}

		// Create a response document without certain fields
		responseDoc := models.ResponseDocument{
			FileName:      doc.FileName,
			FileType:      doc.FileType,
			FileSizeBytes: doc.FileSizeBytes,
		}

		documents = append(documents, responseDoc)
	}

	// Check if there are any documents
	if len(documents) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No documents found"})
		return
	}

	// Return the documents as JSON
	c.JSON(http.StatusOK, documents)
}

// SearchFile queries the MongoDB collection and returns all entries that resemble the provided filename.
func SearchFile(c *gin.Context, fileList *mongo.Collection) {
	// Get the filename parameter from the path
	filenameParam := c.Param("filename")

	// Create a case-insensitive regular expression pattern for the filename
	pattern := "(?i)" + regexp.QuoteMeta(filenameParam)

	// Define the filter to find documents with filenames resembling the provided filename
	filter := bson.M{"fileName": bson.M{"$regex": pattern}}

	// Perform the query to find matching documents
	cursor, err := fileList.Find(context.TODO(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying documents"})
		return
	}
	defer cursor.Close(context.TODO())

	// Iterate through the cursor and collect matching documents
	var matchingDocuments []models.ResponseDocument
	for cursor.Next(context.TODO()) {
		var document models.Document
		if err := cursor.Decode(&document); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error decoding documents"})
			return
		}

		// Create a response document without certain fields
		responseDoc := models.ResponseDocument{
			FileName:      document.FileName,
			FileType:      document.FileType,
			FileSizeBytes: document.FileSizeBytes,
		}

		matchingDocuments = append(matchingDocuments, responseDoc)
	}

	// Check if any matching documents were found
	if len(matchingDocuments) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No matching documents found"})
		return
	}

	// Documents successfully found
	c.JSON(http.StatusOK, matchingDocuments)
}
