A cloud drive that takes advantage of discord's infinite storage thats given to user servers. 
Personally, I use discord to transfer pdf files from my iPad to my computer and back and fourth for school.
But Discord has a restriction for the maximum size of files that you can upload and send which is 25 megabytes. 

This tool utilizes a discord bot library for Go and the Gin web framework to automatically break up large uploaded files into multiple small file chunks, which is then sent to discord channels within a discord server. Channels in a Discord server is the storage location of uploaded files. 

Files are uploaded and downloaded through HTTP requests handled by the GIN REST Api web server. 
An additional mongoDb cluster(can be replaced with any other KV stores like Redis) is used to keep a list of the files uploaded which is used for querying and searching uploaded files.
The KV store in mongoDb is used to keep track of the file names of the uploaded files, the order of individual file chunks associated with the file, and where the file chunks are distributed in the Discord Server(located by discord channelIDs and messageIDs)

Of course, it is not guaranteed that your file data will be safe as Discord is not meant to be used for this purpose. And messages may be cleaned up by Discord upon policy changes. 
Therefore please do not use this seriously, as it is merely a toy project for learning and experimenting with Go and REST Apis.

The download process is just the reverse of the upload process. The program searches for the specified filename on the mongodb collection and then gets the discord message IDs of the file chunks that belong to the file in interest. The file chunks are later retrieved into memory and recomposed to be a complete file. 
The composition and decomposition of the files were done through reading binary bits.

Go concurrency principles were used in this project to drastically improve the performance of upload and download through spawning multiple goroutines each taking care of one individual file chunk and joining them with a wait group. 

The overall procedure is somewhat inspired by the distributed big data processing framework MapReduce where a large file is broken into small pieces for different workers to process, but in a more simplified single machine manner.
