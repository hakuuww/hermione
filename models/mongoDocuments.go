package models

import (
	"bytes"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Document represents the structure of your MongoDB document
type Document struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	FileName      string             `bson:"fileName"`
	FileType      string             `bson:"fileType, omitempty"`
	NumOfChunks   int64              `bson:"numOfChunks"`
	FileSizeBytes int64              `bson:"fileSizeBytes"`
	FileChunks    []FileChunk        `bson:"fileChunks"`
	NomralChunkSize int64 `bson:"normralChunkSize"`
	LastChunkSize int64 `bson:"lastChunkSize"`
}

// FileChunk represents a chunk of a file in your MongoDB document
type FileChunk struct {
	MessageID      string `bson:"messageId"`
	ChannelID      string `bson:"channelId"`
	SequenceNumber int    `bson:"sequenceNumber"`
	ChunkSizeBytes int64  `bson:"chunkSizeBytes"`
	StartingByteIndexofFile int64 `bson:"startingByteIndexofFile"`
	EndingByteIndexofFile int64 `bson:"endingByteIndexofFile"`
}

// ResponseDocument represents the response structure without certain fields
type ResponseDocument struct {
	FileName      string             `json:"fileName"`
	FileType      string             `json:"fileType,omitempty"`
	FileSizeBytes int64              `json:"fileSizeBytes"`
	// Omitting the fields that are not needed by the frontend
}

type ChunkData struct {
	Buf *bytes.Buffer
	Seq int 
	Size int64
}

// DocumentService provides methods to interact with Document entities.
type DocumentService interface {
	recombineFileChunks() (bytes.Buffer, error)
}


