# Redis from Scratch

A Go Implementation of an In-Memory Key-Value Store.

This project is a Go-based implementation of a simplified in-memory key-value store, inspired by the Redis database. It provides a foundational understanding of building a network server and handling the RESP protocol. This project is ideal for learning about networking, concurrency, and data structures in Go.

## Technology Stack

* **Programming Language:** Go (Go 1.23.1 or later recommended)
* **Networking:** TCP sockets
* **Serialization:** RESP (using the `github.com/tidwall/resp` library)

## Project Structure

The project is divided into two main components:
`/server`: This is responsible for acting as a Redis server with in memory key value storage.
`/example`: This has a custom made client and the `go-redis` client library showing how to interact with the server.


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

## Documentation
Redis communicates over raw TCP, and to maximize performance, our implementation will also use raw TCP.

### TCP
TCP is the underlying protocol used by protocols like HTTP, SSH and others you're probably familiar with. Redis clients & servers use TCP to communicate with each other.

Redis clients communicate with Redis servers by sending "commands". For each command, a Redis server sends a response back to the client. Commands and responses are both encoded using the [Redis protocol](https://redis.io/docs/latest/develop/reference/protocol-spec/).

**PING** is one of the simplest Redis commands. It's used to check whether a Redis server is healthy.

The response for the PING command is **+PONG\r\n**. This is the string "PONG" encoded using the Redis protocol.


### Peers in TCP Connections
In Go, a peer refers to any entity (client or server) that participates in a TCP connection. Each peer has an associated IP address and port.


#### Peers in a Client-Server Model
- A client peer initiates a connection to a server peer using a TCP socket.
- The server peer listens for incoming connections and establishes a TCP session with clients.
- Both entities are peers since they communicate directly over a bi-directional TCP stream.
- Identifying Peers in Golang

In Go, you can retrieve peer addresses using:
- net.Conn.RemoteAddr() – gets the remote peer's IP and port.
- net.Conn.LocalAddr() – gets the local peer's IP and port.

Once connected, both peers exchange data over the same TCP stream, regardless of their role as a client or server.


## Single Threaded Event Loop
Redis operates using a single-threaded event loop, meaning it processes one command at a time in a sequential manner. This design choice simplifies concurrency by avoiding the complexities of locks and race conditions, making it highly efficient for the I/O-bound tasks it handles. However, this approach isn’t well-suited to Go, which is built around lightweight goroutines and channels that allow concurrent execution without the need to restrict tasks to a single thread. In Go, leveraging these native concurrency features leads to more idiomatic and efficient code, rather than forcing a single-threaded event loop pattern that could undermine the advantages of the language. Therefore the goal is to leverage goroutines and channels for concorrency in this project.


## Further Development
- Improve code readability
- Improve the in and out of the GET and SET methods to use []byte instead of string
- See if there are other core features to add using Code Crafters