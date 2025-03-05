package main

import (
	"bytes"
	"fmt"

	"github.com/tidwall/resp"
)

const (
	CommandSet  = "SET"
	CommandGet  = "GET"
	CommandPing = "PING"
)

type Command interface{}

type SetCommand struct {
	key, val string
}

type GetCommand struct {
	key string
}

type PingCommand struct{}

var ErrUnknownCommand = fmt.Errorf("unknown command received")

func parseREPLtoCommand(raw string) (Command, error) {
	rd := resp.NewReader(bytes.NewBufferString(raw))

	value, _, err := rd.ReadValue()
	if err != nil {
		panic(err)
	}

	if value.Type() == resp.Array {
		if len(value.Array()) != 3 {
			return nil, ErrUnknownCommand
		}

		for _, val := range value.Array() {
			switch val.String() {
			case CommandSet:
				cmd := SetCommand{
					key: value.Array()[1].String(), // key
					val: value.Array()[2].String(), // value
				}
				return cmd, nil
			case CommandGet:
				cmd := GetCommand{
					key: value.Array()[1].String(), // key
				}
				return cmd, nil
			}
		}
	}

	if value.Type() == resp.BulkString {
		switch value.String() {
		case CommandPing:
			return PingCommand{}, nil
		}
	}

	return nil, ErrUnknownCommand
}

func parseStringtoREPL(msg string) ([]byte, error) {
	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)
	err := wr.WriteString(msg)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
