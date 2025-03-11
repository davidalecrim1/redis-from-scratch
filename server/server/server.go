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
	peers map[*internal.Peer]struct{}
	kvs   *internal.KeyValueStorage
}

func NewServer(cfg Config) *Server {
	if cfg.ListenAddr == "" {
		cfg.ListenAddr = fmt.Sprintf(":%d", defaultPort)
	}

	return &Server{
		Config: cfg,
		peers:  make(map[*internal.Peer]struct{}),
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

	return s.acceptConnectionsLoop()
}

func (s *Server) acceptConnectionsLoop() error {
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
	peer := internal.NewPeer(conn)
	s.AddPeer(peer)
	defer s.DeletePeer(peer)

	slog.Info("new peer connected", "remoteAddr", conn.RemoteAddr(), "localAddr", conn.LocalAddr())

	readCallback := make(chan []byte, 1)
	go peer.Read(readCallback)

outer:
	for {
		select {
		case message, ok := <-readCallback:
			if !ok {
				slog.Debug("stoping the handle new messages")
				break outer
			}

			response, err := s.handleMessage(message)
			if err != nil {
				slog.Error("received an error while handling message", "error", err)
				return
			}

			if _, err := peer.Send(response); err != nil {
				slog.Error("failed to write message to connection", "error", err)
				return
			}
		}
	}
}

func (s *Server) AddPeer(p *internal.Peer) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.peers[p] = struct{}{}
}

func (s *Server) DeletePeer(p *internal.Peer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.peers, p)
}

func (s *Server) Close() error {
	for p := range s.peers {
		s.DeletePeer(p)
	}

	return s.ln.Close()
}

func (s *Server) handleMessage(message []byte) (response []byte, err error) {
	parsedMessage, err := internal.ParseReplToCommand(string(message))
	if err != nil {
		return nil, err
	}
	switch cmd := parsedMessage.(type) {
	case internal.SetCommand:
		if err := s.kvs.Set(cmd.Key, cmd.Val); err != nil {
			return nil, err
		}

		response, err := internal.ParseStringToREPL("OK")
		if err != nil {
			return nil, err
		}

		return response, nil

	case internal.GetCommand:
		// TODO: Do I actually need to return an error here if the key is invalid?
		val, err := s.kvs.Get(cmd.Key)

		if err != nil && errors.Is(err, internal.ErrKeyDoesntExist) {
			slog.Error("received an error while 'getting' a value from KVS", "error", err)
			return internal.ParseNilToREPL()
		}

		return internal.ParseStringToREPL(string(val))

	case internal.PingCommand:
		return internal.ParseStringToREPL("PONG")

	case internal.HelloCommand:
		resp := map[string]string{
			"server": "redis",
		}
		return internal.ParseMaptoREPL(resp), nil

	case internal.ClientCommand:
		return internal.ParseStringToREPL("OK")

	case internal.EchoCommand:
		return internal.ParseStringToREPL(cmd.Value)

	default:
		return nil, fmt.Errorf("unknown command type '%v'", cmd)

	}
}
