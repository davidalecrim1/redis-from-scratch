package main

import (
	"bytes"
	"fmt"
	"log/slog"

	"github.com/tidwall/resp"
)

const (
	CommandSet  = "SET"
	CommandGet  = "GET"
	CommandPing = "PING"
)

type Command interface{}

type SetCommand struct {
	key, val []byte
}

type GetCommand struct {
	key []byte
}

type PingCommand struct{}

var ErrUnknownCommand = fmt.Errorf("unknown command received")

func parseREPL(raw string) ([]Command, error) {
	cmds := make([]Command, 0, 1) // at least one command should be received
	rd := resp.NewReader(bytes.NewBufferString(raw))

	for {
		value, _, err := rd.ReadValue()

		if err != nil && err.Error() == "EOF" {
			return cmds, nil
		}

		if err != nil {
			slog.Error("received an unexpected error", "error", err)
			panic(err)
		}

		if value.Type() == resp.Array {
			if len(value.Array()) != 2 && len(value.Array()) != 3 {
				return nil, ErrUnknownCommand
			}

			for _, val := range value.Array() {
				switch val.String() {
				case CommandSet:
					cmd := SetCommand{
						key: value.Array()[1].Bytes(), // key
						val: value.Array()[2].Bytes(), // value
					}
					cmds = append(cmds, cmd)
				case CommandGet:
					cmd := GetCommand{
						key: value.Array()[1].Bytes(), // key
					}
					cmds = append(cmds, cmd)
				}
			}
		}

		if value.Type() == resp.BulkString {
			switch value.String() {
			case CommandPing:
				cmds = append(cmds, PingCommand{})
			}
		}
	}
}

func parseNilToREPL() ([]byte, error) {
	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)
	err := wr.WriteNull()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func parseStringToREPL(msg string) ([]byte, error) {
	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)
	err := wr.WriteString(msg)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
