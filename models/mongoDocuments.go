package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Document represents the structure of your MongoDB document
type Document struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	FileName      string             `bson:"fileName"`
	FileType      string             `bson:"fileType, omitempty"`
	NumOfChunks   int64              `bson:"numOfChunks"`
	FileSizeBytes int64              `bson:"fileSizeBytes"`
	FileChunks    []FileChunk        `bson:"fileChunks"`
}

// FileChunk represents a chunk of a file in your MongoDB document
type FileChunk struct {
	MessageID      string `bson:"messageId"`
	ChannelID      string `bson:"channelId"`
	SequenceNumber int    `bson:"sequenceNumber"`
	ChunkSizeBytes int64  `bson:"chunkSizeBytes"`
}
