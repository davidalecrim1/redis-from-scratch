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

		slog.Debug("received data from peer", "bytesRead", n, "dataReceived", string(buf[:n]))

		msgBuf := make([]byte, n)
		copy(msgBuf, buf[:n])

		cmd, err := parseREPLtoCommand(string(msgBuf))
		slog.Debug("received a message", "message", string(msgBuf))

		p.msgCh <- Message{
			cmd:  cmd,
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
