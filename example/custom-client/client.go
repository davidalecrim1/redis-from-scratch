// Example using a custom made client
package main

import (
	"bytes"
	"context"
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

func (c *Client) Set(ctx context.Context, key string, value string) (string, error) {
	c.MustCreateConn()

	var writeBuf bytes.Buffer
	wr := resp.NewWriter(&writeBuf)
	err := wr.WriteArray([]resp.Value{
		resp.StringValue("SET"),
		resp.StringValue(key),
		resp.StringValue(value),
	})
	if err != nil {
		return "", err
	}

	_, err = c.conn.Write(writeBuf.Bytes())
	if err != nil {
		return "", err
	}

	msg, err := c.readString(ctx)
	if err != nil {
		return "", err
	}

	return msg, nil
}

func (c *Client) Ping(ctx context.Context) (string, error) {
	c.MustCreateConn()

	err := c.writeString(ctx, "PING")
	if err != nil {
		return "", err
	}

	return c.readString(ctx)
}

func (c *Client) writeString(_ context.Context, value string) error {
	var writeBuf bytes.Buffer
	wr := resp.NewWriter(&writeBuf)
	err := wr.WriteString(value)
	if err != nil {
		return err
	}

	_, err = io.Copy(c.conn, &writeBuf)
	return err
}

func (c *Client) readString(ctx context.Context) (string, error) {
	value, err := c.read(ctx)
	if err != nil {
		return "", err
	}

	resp := value.String()
	return resp, nil
}

func (c *Client) read(_ context.Context) (resp.Value, error) {
	readBuf := make([]byte, 1024)
	n, err := c.conn.Read(readBuf)
	if err != nil {
		return resp.NullValue(), err
	}

	respBuf := bytes.NewBuffer(readBuf[:n])
	rd := resp.NewReader(respBuf)
	value, _, err := rd.ReadValue()
	if err != nil {
		return resp.NullValue(), err
	}

	return value, nil
}

func (c *Client) Get(ctx context.Context, key string) (string, error) {
	c.MustCreateConn()

	var writeBuf bytes.Buffer
	wr := resp.NewWriter(&writeBuf)
	err := wr.WriteArray([]resp.Value{
		resp.StringValue("GET"),
		resp.StringValue(key),
	})
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(c.conn, &writeBuf); err != nil {
		return "", err
	}

	return c.readString(ctx)
}

func (c *Client) Close() error {
	if c.conn != nil {
		err := c.conn.Close()
		c.conn = nil
		return err
	}
	return nil
}
