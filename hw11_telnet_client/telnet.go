package main

import (
	"errors"
	"io"
	"net"
	"time"
)

type TelnetClient struct {
	addr   string
	dialer net.Dialer
	conn   net.Conn
	in     io.ReadCloser
	out    io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) *TelnetClient {
	d := net.Dialer{Timeout: timeout}
	return &TelnetClient{dialer: d, addr: address, out: out, in: in}
}

func (tc *TelnetClient) Connect() error {
	conn, err := tc.dialer.Dial("tcp", tc.addr)
	if err != nil {
		return err
	}

	tc.conn = conn
	return nil
}

func (tc *TelnetClient) Send() error {
	if tc.conn == nil || tc.in == nil {
		return errors.New("invalid send connection")
	}

	_, err := io.Copy(tc.conn, tc.in)
	return err
}

func (tc *TelnetClient) Receive() error {
	if tc.conn == nil || tc.out == nil {
		return errors.New("invalid receive connection")
	}

	_, err := io.Copy(tc.out, tc.conn)
	return err
}

func (tc *TelnetClient) Close() error {
	if tc.conn == nil {
		return errors.New("no connection")
	}

	return tc.conn.Close()
}
