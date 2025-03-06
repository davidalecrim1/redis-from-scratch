# Redis from Scratch

// TODO: Improve this docs, now it's just random important notes.

The goal is the create a redis like in-memory database. The idea is to explore how things are built from scratch to enhance the engineering skills.


rename this to redis from scratch on my github

## Documentation

Redis works with raw tcp, so we are going to work with raw tcp too to have a great performance.

## Peers

In the context of TCP connections using Golang, a peer refers to any entity (client or server) that participates in a TCP connection. Each peer in a TCP connection has an associated IP address and port.

Peers in a Client-Server Model:

A client (peer) initiates a connection to a server (peer) using a TCP socket.
The server listens for incoming connections and establishes a TCP session with the client.
Both entities are considered peers because they are communicating directly over a bi-directional TCP stream.

Each peer is identified by its IP address and port.

Peers can be either clients or servers, but once connected, they both exchange data over the same TCP stream.

In Golang, net.Conn.RemoteAddr() and net.Conn.LocalAddr() can be used to get peer addresses.

## Folders
- Client folder is for sample redis client like
- Server is the redis server like memory database


## RESP Protocol

using the same from redis. using the resp library that is archived no be focus too much on it.


## Inspirations
https://github.com/tidwall/redcon
https://app.codecrafters.io/
https://www.youtube.com/watch?v=LMrxfWB6sbQ


## Future Features
- Use the actual redis client to test the server
