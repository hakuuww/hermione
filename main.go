package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/hakuuww/hermione/database"
	"github.com/hakuuww/hermione/routes"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"os/signal"
	"github.com/joho/godotenv"
)

var (
	server      *gin.Engine
	ctx         context.Context
	userc       *mongo.Collection
	mongoClient *mongo.Client
	err         error
)

const GuildID = "1181344672010997820"
const channelID = "1182079083249680404"

func main() {
	// Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	mongoClient, err = database.InitDB()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	dbName := "discordFileStorageServer"
	db := mongoClient.Database(dbName)
	collectionName := "fileList"
	fileList := db.Collection(collectionName)

	//init discord
	dg := discordGoInit()

	go func() { // Start the server
		server = routes.SetupRouter(dg, fileList)
		server.MaxMultipartMemory = 1000 << 20  
		port := 8081
		err = server.Run(fmt.Sprintf("localhost:%d", port))
		if err != nil {
			panic(err)
		}
	}()

	// Wait for a signal to exit (e.g., Ctrl+C)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Cleanly close down the Discord session.
	defer dg.Close()
}

func discordGoInit() *discordgo.Session {
	// Create a new Discord session using the provided bot token.
	//token := os.Getenv("DISCORD_BOT_TOKEN")
	dg, err := discordgo.New(os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return nil
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return nil
	}

	return dg
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
