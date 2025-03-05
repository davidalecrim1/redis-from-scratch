package main

import (
	"bytes"
	"context"
	"fmt"
	"net"

	"github.com/tidwall/resp"
)

type Client struct {
	addr string
}

func NewClient(addr string) *Client {
	return &Client{
		addr: addr,
	}
}

func (c *Client) Set(ctx context.Context, key string, value string) error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return err
	}

	var writeBuf bytes.Buffer
	wr := resp.NewWriter(&writeBuf)
	err = wr.WriteArray([]resp.Value{
		resp.StringValue("SET"),
		resp.StringValue(key),
		resp.StringValue(value),
	})
	if err != nil {
		return err
	}

	_, err = conn.Write(writeBuf.Bytes())
	return err
}

func (c *Client) Ping(ctx context.Context) (string, error) {
	conn, err := net.Dial("tcp", c.addr)

	err = c.writePing(ctx, conn)
	if err != nil {
		return "", err
	}

	return c.readPong(ctx, conn)
}

func (c *Client) writePing(_ context.Context, conn net.Conn) error {
	var writeBuf bytes.Buffer
	wr := resp.NewWriter(&writeBuf)
	err := wr.WriteString("PING")
	if err != nil {
		return err
	}

	_, err = conn.Write(writeBuf.Bytes())
	return err
}

func (c *Client) readPong(_ context.Context, conn net.Conn) (string, error) {
	readBuf := make([]byte, 1024)
	n, err := conn.Read(readBuf)
	if err != nil {
		return "", err
	}

	respBuf := bytes.NewBuffer(readBuf[:n])
	rr := resp.NewReader(respBuf)
	value, _, err := rr.ReadValue()
	if err != nil {
		panic(err)
	}

	if value.String() != "PONG" {
		return "", fmt.Errorf("unexpected response: %v", value)
	}

	return value.String(), nil
}
