# The Joker
The Joker is an HTTP server that serves lame jokes.

# Not your average Joker
The Joker is built with Go's `http.Server` but wrapped. The wrapping mainly extends the STL server to have a middleware that tracks the number of hits to the server in the last `n` seconds. The number of hits is fault-tolerant by persisting the values in files - `.hits` and `.windowmetadata` in the root directory. This protects continuity in case Joker kills/terminates itself. `n` is configurable with environmental variable - `SERVER_WINDOW_SIZE_IN_SECONDS`. It is set to `60`. If there is no `SERVER_WINDOW_SIZE_IN_SECONDS`, it defaults to `10`.

# Instructions
Instructions to start the server
## Requirements
1. `go version 1.22`
2. `GNU Make version 4.3`

## Environment Variables
There are 3 environment variables,
1. `SERVER_WINDOW_SIZE_IN_SECONDS` - Number of Hits in the last `n` seconds, `n` being `SERVER_WINDOW_SIZE_IN_SECONDS`
2. `PROJECT_ROOT` - Project Root
3. `SERVER_PORT`  - Port the Server should listen on

The variables are set to their defaults in the `Makefile`

## How to start the server
`make start-server` from the root directory will start the server.

## Cleanup
`make clean` deletes the `.hits` and  `.windowmetadata` files in the root directory.

## Usage
The server can be used in `Postman`/similar by using the API  
 GET `http://localhost:8080/joke`  
or using `curl`  
  `curl -X GET http://localhost:8080/joke`


