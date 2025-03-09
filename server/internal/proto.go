package internal

import (
	"bytes"
	"fmt"
	"log/slog"
	"strings"

	"github.com/tidwall/resp"
)

const (
	CommandSet    = "set"
	CommandGet    = "get"
	CommandPing   = "ping"
	CommandHello  = "hello"
	CommandClient = "client"
)

func ParseREPL(raw string) ([]Command, error) {
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
			for _, val := range value.Array() {
				switch strings.ToLower(val.String()) {
				case CommandSet:
					cmd := SetCommand{
						Key: value.Array()[1].Bytes(), // key
						Val: value.Array()[2].Bytes(), // value
					}
					cmds = append(cmds, cmd)
				case CommandGet:
					cmd := GetCommand{
						Key: value.Array()[1].Bytes(), // key
					}
					cmds = append(cmds, cmd)
				case CommandHello:
					cmd := HelloCommand{
						Value: value.Array()[1].String(), // value
					}
					cmds = append(cmds, cmd)
				case CommandClient:
					cmd := ClientCommand{
						Value: value.Array()[1].String(), // ?
					}
					cmds = append(cmds, cmd)
				}
			}
		}

		if value.Type() == resp.BulkString {
			switch strings.ToLower(value.String()) {
			case CommandPing:
				cmds = append(cmds, PingCommand{})
			}
		}
	}
}

func ParseNilToREPL() ([]byte, error) {
	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)
	err := wr.WriteNull()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func ParseStringToREPL(msg string) ([]byte, error) {
	var buf bytes.Buffer
	wr := resp.NewWriter(&buf)
	err := wr.WriteString(msg)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func ParseMaptoREPL(msg map[string]string) []byte {
	var buf bytes.Buffer
	buf.WriteString("%" + fmt.Sprintf("%d\r\n", len(msg)))
	wr := resp.NewWriter(&buf)

	for k, v := range msg {
		wr.WriteString(k)
		wr.WriteString(":" + v)
	}

	return buf.Bytes()
}
