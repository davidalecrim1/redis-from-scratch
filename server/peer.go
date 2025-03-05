package main

import (
	"log/slog"
	"net"
)

type Peer struct {
	conn  net.Conn
	msgCh chan Message
}

func NewPeer(conn net.Conn, msgCh chan Message) *Peer {
	return &Peer{
		conn:  conn,
		msgCh: msgCh,
	}
}

// Reads until receives an `EOF` and returns nil
func (p *Peer) Read() error {
	buf := make([]byte, 1024)

	for {
		n, err := p.conn.Read(buf)
		if err != nil && err.Error() == "EOF" {
			slog.Info("reached the EOF of the current connection", "remoteAddr", p.conn.RemoteAddr())
			return nil
		}

		if err != nil {
			slog.Error("peer failed to read the connection data", "err", err)
			return err
		}

		slog.Debug("received data from peer", "bytesRead", n)

		msgBuf := make([]byte, n)
		copy(msgBuf, buf[:n])

		cmds, err := parseREPL(string(msgBuf))
		slog.Debug("received a message", "message", string(msgBuf))

		p.msgCh <- Message{
			cmds: cmds,
			peer: p,
		}
	}
}

func (p *Peer) Send(msg []byte) (int, error) {
	b, err := parseStringtoREPL(string(msg))
	if err != nil {
		return -1, err
	}

	return p.conn.Write(b)
}
