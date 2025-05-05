package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

var ErrClosedConnection = errors.New("closed connection")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (c *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return fmt.Errorf("tcp connect %w", err)
	}

	c.conn = conn
	return nil
}

func (c *telnetClient) Close() error {
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("tcp close %w", err)
	}
	return nil
}

func (c *telnetClient) Send() error {
	scanner := bufio.NewScanner(c.in)
	if scan := scanner.Scan(); !scan {
		if scanner.Err() != nil {
			return fmt.Errorf("read buffer in")
		}
		return nil
	}

	message := scanner.Bytes()
	message = append(message, '\n')

	_, err := c.conn.Write(message)
	if err != nil {
		c.Close()
		return ErrClosedConnection
	}

	return nil
}

func (c *telnetClient) Receive() error {
	scanner := bufio.NewScanner(c.conn)
	if scan := scanner.Scan(); !scan {
		if err := scanner.Err(); err != nil {
			return fmt.Errorf("read from server %w", err)
		}
		return ErrClosedConnection
	}

	message := scanner.Bytes()
	message = append(message, '\n')
	_, err := c.out.Write(message)
	if err != nil {
		return fmt.Errorf("receive message %w", err)
	}

	return nil
}

// Place your code here.
// P.S. Author's solution takes no more than 50 lines.
