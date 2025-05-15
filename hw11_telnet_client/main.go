package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var timeout time.Duration
	flag.DurationVar(&timeout, "timeout", time.Second*10, "timeout")
	flag.Parse()
	args := flag.Args()
	host, port := args[0], args[1]

	reader := &bytes.Buffer{}
	readerCloser := io.NopCloser(reader)
	writer := bytes.NewBuffer(nil)
	telnetClient := NewTelnetClient(fmt.Sprintf("%s:%s", host, port), timeout, readerCloser, writer)
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ctx, cancel := signal.NotifyContext(ctxWithTimeout, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	err := telnetClient.Connect()
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stderr, "...Connected to %s:%s\n", host, port)

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			message := scanner.Text()
			reader.WriteString(message)
			err := telnetClient.Send()
			if err != nil {
				telnetClient.Close()
				cancel()
				if errors.Is(err, ErrClosedConnection) {
					fmt.Fprintln(os.Stderr, "...Connection was closed")
				}
				return
			}
		}
		cancel()
		fmt.Fprintln(os.Stderr, "...EOF")
	}()

	go func() {
		for {
			err := telnetClient.Receive()
			if err != nil {
				telnetClient.Close()
				return
			}
		}
	}()

	<-ctx.Done()
	telnetClient.Close()
	// Place your code here,
	// P.S. Do not rush to throw context down, think think if it is useful with blocking operation?
}
