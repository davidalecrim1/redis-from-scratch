package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"

	"github.com/tidwall/resp"
)

type Client struct {
	addr string
	conn net.Conn
}

func NewClient(addr string) *Client {
	return &Client{
		addr: addr,
	}
}

func (c *Client) MustCreateConn() {
	if c.conn == nil {
		conn, err := net.Dial("tcp", c.addr)
		if err != nil {
			panic(err)
		}
		c.conn = conn
	}
}

func (c *Client) Set(ctx context.Context, key string, value string) error {
	c.MustCreateConn()

	var writeBuf bytes.Buffer
	wr := resp.NewWriter(&writeBuf)
	err := wr.WriteArray([]resp.Value{
		resp.StringValue("SET"),
		resp.StringValue(key),
		resp.StringValue(value),
	})
	if err != nil {
		return err
	}

	_, err = c.conn.Write(writeBuf.Bytes())
	return err
}

func (c *Client) Ping(ctx context.Context) (string, error) {
	c.MustCreateConn()

	err := c.writePing(ctx)
	if err != nil {
		return "", err
	}

	return c.readPong(ctx)
}

func (c *Client) writePing(_ context.Context) error {
	var writeBuf bytes.Buffer
	wr := resp.NewWriter(&writeBuf)
	err := wr.WriteString("PING")
	if err != nil {
		return err
	}

	_, err = io.Copy(c.conn, &writeBuf)
	return err
}

func (c *Client) readPong(_ context.Context) (string, error) {
	readBuf := make([]byte, 1024)
	n, err := c.conn.Read(readBuf)
	if err != nil {
		return "", err
	}

	respBuf := bytes.NewBuffer(readBuf[:n])
	rr := resp.NewReader(respBuf)
	value, _, err := rr.ReadValue()
	if err != nil {
		return "", err
	}

	if value.String() != "PONG" {
		return "", fmt.Errorf("unexpected response: %v", value)
	}

	return value.String(), nil
}
