package internal

type Message struct {
	Cmds []Command
	Peer *Peer
}
