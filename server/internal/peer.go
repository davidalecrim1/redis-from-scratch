package internal

import (
	"io"
	"log/slog"
	"net"
	"sync"
)

type Peer struct {
	conn net.Conn
	wg   sync.WaitGroup
}

func NewPeer(conn net.Conn) *Peer {
	return &Peer{
		conn: conn,
	}
}

func (p *Peer) Read(readCallback chan<- []byte) {
	buf := make([]byte, 1024)

	for {
		n, err := p.conn.Read(buf)
		if err != nil && err == io.EOF {
			slog.Info("reached the EOF of the current connection", "remoteAddr", p.conn.RemoteAddr())
			go p.Close(readCallback)
			return
		}

		if err != nil {
			slog.Error("peer failed to read the connection data", "err", err)
			return
		}

		slog.Debug("received data from peer", "message", string(buf[:n]), "bytesRead", n)
		p.wg.Add(1)
		readCallback <- buf[:n]
	}
}

func (p *Peer) Send(msg []byte) (int, error) {
	defer p.wg.Done()

	n, err := p.conn.Write(msg)
	if err != nil {
		return n, err
	}

	slog.Debug("sending data to the client", "message", string(msg))
	return n, err
}

func (p *Peer) Close(readCallback chan<- []byte) error {
	p.wg.Wait()

	close(readCallback)
	err := p.conn.Close()
	slog.Debug("closing the connection on peer", "remoteAddr", p.conn.RemoteAddr())
	return err
}
