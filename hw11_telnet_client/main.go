package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var timeout = flag.Duration("timeout", 10*time.Second, "timeout description")

func main() {
	l := len(os.Args)
	if l < 3 {
		log.Fatalln("input error: not enough arguments")
	}

	flag.Parse()

	host, port := os.Args[l-2], os.Args[l-1]
	addr := net.JoinHostPort(host, port)

	client := NewTelnetClient(addr, *timeout, os.Stdin, os.Stdout)

	err := client.Connect()
	if err != nil {
		log.Fatalf("connect error: connect to %s is failed", addr)
	}

	msg := fmt.Sprintf("...Connected to %s", addr)
	os.Stderr.Write([]byte(msg + "\n"))

	ctx, stop := signal.NotifyContext(context.TODO(), syscall.SIGHUP, syscall.SIGINT)
	defer stop()

	ch := make(chan error, 2)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for {
			if err := client.Receive(); err != nil {
				ch <- err
				break
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			if err := client.Send(); err != nil {
				ch <- err
				break
			}
		}
	}()

	select {
	case <-ctx.Done():
		msg = "...Quit"
		break
	case err = <-ch:
		msg = err.Error()
		break

	}

	err = client.Close()
	if err != nil {
		fmt.Println(err)
	}

	wg.Wait()
	os.Stderr.Write([]byte(msg + "\n"))
}
