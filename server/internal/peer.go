package internal

import (
	"io"
	"log/slog"
	"net"
	"sync"
)

type Peer struct {
	conn  net.Conn
	wg    sync.WaitGroup
	msgCh chan<- Message
	delCh chan<- *Peer
}

func NewPeer(conn net.Conn, msgCh chan<- Message, delCh chan<- *Peer) *Peer {
	return &Peer{
		conn:  conn,
		msgCh: msgCh,
		delCh: delCh,
	}
}

// Reads until receives an `EOF` and returns nil
func (p *Peer) Read() error {
	buf := make([]byte, 1024)

	for {
		n, err := p.conn.Read(buf)
		if err != nil && err == io.EOF {
			slog.Info("reached the EOF of the current connection", "remoteAddr", p.conn.RemoteAddr())
			go p.Close()
			return nil
		}

		if err != nil {
			slog.Error("peer failed to read the connection data", "err", err)
			return err
		}

		slog.Debug("received data from peer", "bytesRead", n)

		msgBuf := make([]byte, n)
		copy(msgBuf, buf[:n])

		// TODO: maybe i can increase performance using a channel to dispatch the commands to a channel
		// one by one to be handled instead of using only a message with multiple commands
		// Think this later
		cmds, err := ParseReplToCommand(string(msgBuf))
		if err != nil {
			slog.Error("received an error when parsing the REPL to command", "error", err)
			return err
		}

		slog.Debug("received a message", "message", string(msgBuf))

		p.addReadToWaitForWrite(len(cmds))

		p.msgCh <- Message{
			Cmds: cmds,
			Peer: p,
		}
	}
}

// Waits for each command to respond in write to close the connection
func (p *Peer) addReadToWaitForWrite(len int) {
	p.wg.Add(len)
}

func (p *Peer) Send(msg []byte) (int, error) {
	defer p.wg.Done()

	n, err := p.conn.Write(msg)
	if err != nil {
		return n, err
	}

	slog.Debug("sending data to the client", "msg", string(msg))
	return n, err
}

func (p *Peer) Close() error {
	p.wg.Wait()
	err := p.conn.Close()
	slog.Debug("closing the connection on peer", "remoteAddr", p.conn.RemoteAddr())
	p.delCh <- p
	return err
}
