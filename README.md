# Hermione: Your own free unlimited cloud drive using Discord

![image](https://github.com/user-attachments/assets/dfb29f36-55ce-427c-9ae0-5b36f1ebe345)

![Go](https://img.shields.io/badge/Language-Go-00ADD8)
![MongoDB](https://img.shields.io/badge/Database-MongoDB-47A248)
![Gin](https://img.shields.io/badge/Framework-Gin-00b140)
![Discord API](https://img.shields.io/badge/API-Discord-7289da)
![Discord Go](https://img.shields.io/badge/Library-DiscordGo-7289da)

Hermione is a toy project to make use of Discord’s free unlimited storage provided to user servers.
Named after Hermione Granger from *Harry Potter*.
Just like Hermione who can remember everything, this project allows users to upload, store, and retrieve large files through Discord channels for free. 
Files transfer are through HTTP, a large file is physically split into chunks and split around different destinations on Discord to work around Discord’s 25 MB file upload limitation. 

This project demonstrates the use of a Discord bot library for Go, Gin web framework, MongoDB (or any KV store like Redis), and Go concurrency principles to facilitate a very basic Go web server with a database backend.

> **Disclaimer**: This project is meant purely for learning and experimenting with Go, REST APIs, and exploiting Discord’s infrastructure in creative ways. Files may not be safe long-term as Discord’s policies may change, and messages may be cleaned up.
> **This is only a toy project, I only uploaded a few files for testing purposes, please don't sue me Discord.**
> **Do not use this for anything critical!**

## Implementation
- **File Splitting**: Large files are automatically split into 25 MB chunks for uploading.
- **HTTP API**: Use the Gin web framework to handle file upload and download requests.
- **Discord API**: Implemented a Discord bot that can read, write and delete files in specified channels(directories) within a Discord server.
- **MongoDB**: Tracks uploaded files, their names, type, meta data, their chunk distributions/locations in specific Discord channels within a server.
- **Slow but working file upload and retrieval**: Once uploaded through HTTP Rest to the Web Server, Files are decomposed into and recomposed from under 25MB file chunks.
- **Concurrency/Optimization**: Go’s concurrency model (goroutines and wait groups) drastically improves upload and download performance. A single Goroutine is responsible of its own chunk and the Discord API logic for writing and retrieving it.

## Technologies

- **Go**: The core language used to implement the bot, API, and file chunking mechanisms.
- **Gin**: A fast web framework for building RESTful APIs in Go.
- **MongoDB**: (You can replace this with any KV database, or relational database as well).Key-value store to track file metadata, chunking info, and Discord message IDs.
- **Discord API**: Used to upload and manage file chunks in Discord channels as a Discord Bot.

## Setup

### Prerequisites
- MongoDB (or any other KV store by implementing your own interface).
- A Discord bot token (you need to create a bot on the Discord Developer portal).
- Discord server with appropriate permissions to upload files.

### Setup

1. Configure the Discord bot token and channel IDs by editing the `config.json` file:
    ```json
    {
      "discord_token": "YOUR_DISCORD_BOT_TOKEN",
      "mongo_uri": "mongodb://localhost:27017",
      "mongo_db_name": "hermione",
      "discord_channel_ids": [
        "YOUR_CHANNEL_ID_1",
        "YOUR_CHANNEL_ID_2"
      ]
    }
    ```

2. Start the API server:
    ```bash
    go run main.go
    ```

## Usage

### Uploading a File

To upload a file, send an HTTP `POST` request to the `/upload` endpoint. The file will be automatically split into chunks and sent to the Discord server.

Example (using `curl`):

```bash
curl -X POST -F "file=@path_to_your_large_file" http://localhost:8080/upload
```

### Downloading a File

To download a previously uploaded file, send an HTTP `GET` request to the `/download/{filename}` endpoint.

Example:

```bash
curl http://localhost:8080/download/filename.pdf --output filename.pdf
```

### Querying Files

You can search for uploaded files using the `/search/{filename}` endpoint. This will return metadata about the file, including where each chunk is stored.

Example:

```bash
curl http://localhost:8080/search/filename.pdf
```

### License

This project is licensed under the MIT License - see the LICENSE file for details.

### Disclaimer

This tool is a proof of concept for learning purposes. Using Discord to store files in this manner may violate Discord’s Terms of Service, and there is a risk that your files may be lost or deleted without warning. Use it at your own risk.








