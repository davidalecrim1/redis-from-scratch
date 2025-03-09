package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"

	"redis-from-scratch/internal"
)

const defaultPort = 6379

type Config struct {
	ListenAddr string
}

type Server struct {
	Config
	ln    net.Listener
	mu    sync.RWMutex
	peers map[*internal.Peer]bool
	delCh chan *internal.Peer
	msgCh chan internal.Message
	kvs   *internal.KeyValueStorage
}

func NewServer(cfg Config) *Server {
	if cfg.ListenAddr == "" {
		cfg.ListenAddr = fmt.Sprintf(":%d", defaultPort)
	}

	return &Server{
		Config: cfg,
		peers:  make(map[*internal.Peer]bool),
		msgCh:  make(chan internal.Message), // TODO: make this buffered to improve performance
		delCh:  make(chan *internal.Peer),
		kvs:    internal.NewKeyValueStorage(),
	}
}

func (s *Server) Start(ctx context.Context) error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		slog.Error("failed to listen for incoming tcp connections", "error", err)
		return err
	}
	s.ln = ln
	slog.Info("listening tcp connection", "port:", s.ListenAddr)

	go s.watchMessages(ctx)
	go s.watchClose(ctx)

	return s.acceptLoop()
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("failed to accept incoming connection", "error", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := internal.NewPeer(conn, s.msgCh, s.delCh)
	s.AddPeer(peer)

	slog.Info("new peer connected", "remoteAddr", conn.RemoteAddr(), "localAddr", conn.LocalAddr())

	// TODO: Should I close the connecton one I reach the EOF of the connection?
	// Think this later
	if err := peer.Read(); err != nil {
		slog.Error("failed to read from peer", "error", err, "remoteAddr", conn.RemoteAddr(), "localAddr", conn.LocalAddr())
	}
}

func (s *Server) AddPeer(p *internal.Peer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.peers[p] = true
}

func (s *Server) DeletePeer(p *internal.Peer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.peers, p)
}

func (s *Server) watchMessages(ctx context.Context) {
	for {
		select {
		case msg := <-s.msgCh:
			if err := s.handleMessage(msg); err != nil {
				slog.Error("failed to handle message", "error", err)
			}
		case <-ctx.Done():
			slog.Debug("context was canceled, closing watchMessages")
			return
		}
	}
}

func (s *Server) watchClose(ctx context.Context) {
	for {
		select {
		case peer, ok := <-s.delCh:
			if !ok {
				slog.Debug("the channel is closed")
			}

			s.DeletePeer(peer)
			return
		case <-ctx.Done():
			slog.Debug("context was canceled, closing watchClose")
			return
		}
	}
}

func (s *Server) handleMessage(msg internal.Message) error {
	// TODO: This msg.peer is bothering me. Doesnt make sense to access the peer using the msg.
	// rethink this later

	for _, cmd := range msg.Cmds {
		switch receivedCmd := cmd.(type) {
		case internal.SetCommand:
			err := s.kvs.Set(receivedCmd.Key, receivedCmd.Val)
			if err != nil {
				slog.Error("received an error while 'setting' a value from KVS", "error", err)
				return err
			}

			resp, err := internal.ParseStringToREPL("OK")
			if err != nil {
				return err
			}

			_, err = msg.Peer.Send(resp)
			return err

		case internal.GetCommand:
			val, err := s.kvs.Get(receivedCmd.Key)
			// TODO: Do I actually need to return an error here if the key is invalid?
			if err != nil && errors.Is(err, internal.ErrKeyDoesntExist) {
				nilMsg, er := internal.ParseNilToREPL()
				if er != nil {
					return er
				}

				_, er = msg.Peer.Send(nilMsg)
				return er
			}

			if err != nil {
				slog.Error("received an error while 'getting' a value from KVS", "error", err)
				return err
			}
			resp, err := internal.ParseStringToREPL(string(val))
			if err != nil {
				return err
			}
			_, err = msg.Peer.Send(resp)
			return err

		case internal.PingCommand:
			resp, err := internal.ParseStringToREPL("PONG")
			if err != nil {
				return err
			}
			_, err = msg.Peer.Send(resp)
			return err

		case internal.HelloCommand:
			resp := map[string]string{
				"server": "redis",
			}

			_, err := msg.Peer.Send(internal.ParseMaptoREPL(resp))
			return err
		case internal.ClientCommand:
			resp, err := internal.ParseStringToREPL("OK")
			if err != nil {
				return err
			}

			_, err = msg.Peer.Send(resp)
			return err
		default:
			return fmt.Errorf("unknown command type '%v'", cmd)
		}
	}

	return nil
}

func (s *Server) Close() error {
	for p := range s.peers {
		s.DeletePeer(p)
	}

	return s.ln.Close()
}
