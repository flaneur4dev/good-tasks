package main

import (
	"bufio"
	"errors"
	"io"
	"net"
	"time"
)

type telnetClient struct {
	addr       string
	dialer     net.Dialer
	conn       net.Conn
	in         io.ReadCloser
	inScanner  *bufio.Scanner
	out        io.Writer
	outScanner *bufio.Scanner
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) *telnetClient {
	d := net.Dialer{Timeout: timeout}
	return &telnetClient{dialer: d, addr: address, out: out, in: in}
}

func (tc *telnetClient) Connect() error {
	conn, err := tc.dialer.Dial("tcp", tc.addr)
	if err != nil {
		return err
	}

	tc.conn = conn
	tc.outScanner = bufio.NewScanner(conn)
	tc.inScanner = bufio.NewScanner(tc.in)
	return nil
}

func (tc *telnetClient) Send() error {
	if !tc.inScanner.Scan() {
		return errors.New("...EOF")
	}

	tc.conn.Write([]byte(tc.inScanner.Text() + "\n"))
	return nil
}

func (tc *telnetClient) Receive() error {
	if !tc.outScanner.Scan() {
		return errors.New("...Connection was closed by peer")
	}

	tc.out.Write([]byte(tc.outScanner.Text() + "\n"))
	return nil
}

func (tc *telnetClient) Close() error {
	err := tc.conn.Close()

	err2 := tc.in.Close()
	if err2 != nil {
		err = err2
	}

	return err
}
