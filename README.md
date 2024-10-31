A cloud drive that takes advantage of discord's infinite storage thats given to user servers. 
Personally, I use discord to transfer pdf files from my iPad to my computer and back and fourth for school.
But Discord has a restriction for the maximum size of files that you can upload and send which is 25mb. 

This tool utilizes a discord bot library for Go and the Gin web framework to automatically break up uploaded files into file chunks that are small, which is then sent over to a discord channel in a discord server. 

An additional mongoDb cluster(can be replaced with any other KV stores like Redis) is used to keep a list of the files uploaded which is used for querying and searching uploaded files.

The download process is just the reverse of the upload process. The program searches for the specified filename on the mongodb collection and then gets the discord message IDs of the file chunks that belong to the file in interest. The file chunks are later retrieved into memory and recomposed to be a complete file. 
The composition and decomposition of the files were done through reading binary bits.

Go concurrency principles were used in this project to drastically improve the performance of upload and download through spawning multiple goroutines each taking care of one individual file chunk and joining them with a wait group. 

The overall procedure is somewhat inspired by the distributed big data processing framework MapReduce where a large file is broken into small pieces for different workers to process, but in a more simplified single machine manner.
