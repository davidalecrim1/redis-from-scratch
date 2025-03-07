# Redis from Scratch

A Go Implementation of an In-Memory Key-Value Store.

This project is a Go-based implementation of a simplified in-memory key-value store, inspired by the Redis database. It provides a foundational understanding of building a network server and handling the RESP protocol. This project is ideal for learning about networking, concurrency, and data structures in Go.


## Technology Stack

* **Programming Language:** Go (Go 1.23.1 or later recommended)
* **Networking:** TCP sockets
* **Serialization:** RESP (using the `github.com/tidwall/resp` library)


## Project Structure

The project is divided into two main components:
`/server`: This is responsible for acting as a Redis server com in memory key value store.
`example`: This has a custom made client and the `go-redis` client library showing how to interact with the server.


## Getting Started

1. Clone the repository
2. Navigate to the directory:
3. Run the client and server in two terminals:

First Terminal:
```bash
make run-server
```

Second Terminal:
```bash
make run-client
```


## Architecture

The server uses a simple in-memory map to store key-value pairs. The client interacts with the server via TCP sockets, sending commands and receiving responses using the RESP protocol. Error handling is implemented to ensure robustness.


## Further Development
- Improve code readability
- Improve the in and out of the GET and SET methods to use []byte instead of string
- See if there are other core features to add using Code Crafters