# The Joker
The Joker is an HTTP server that serves lame jokes.

# Not your average Joker
The Joker is built with Go's `http.Server` but wrapped. The wrapping mainly extends the STL server to have a middleware that tracks the number of hits to the server in the last `n` seconds. The number of hits is fault-tolerant by persisting the values in files - `.hits` and `.windowmetadata` in the root directory. This protects continuity in case Joker kills/terminates itself. `n` is configurable with environmental variable - `SERVER_WINDOW_SIZE_IN_SECONDS`. It is set to `60`. If there is no `SERVER_WINDOW_SIZE_IN_SECONDS`, it defaults to `10`.

# How to start the server
## Requirements
1. `go version 1.22`
2. `GNU Make version 4.3`


