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

var timeout = flag.Duration("timeout", 10*time.Second, "timeout for connection")

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) != 2 {
		log.Fatalln("input error: not enough arguments")
	}

	host, port := args[0], args[1]
	addr := net.JoinHostPort(host, port)

	client := NewTelnetClient(addr, *timeout, os.Stdin, os.Stdout)
	err := client.Connect()
	if err != nil {
		log.Fatalf("connect error: connect to %s is failed", addr)
	}

	ctx, stop := signal.NotifyContext(context.TODO(), syscall.SIGHUP, syscall.SIGINT)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer stop()

		if err := client.Receive(); err != nil {
			os.Stderr.Write([]byte(err.Error() + "\n"))

			return
		}
	}()

	go func() {
		defer wg.Done()
		defer stop()

		if err := client.Send(); err != nil {
			os.Stderr.Write([]byte(err.Error() + "\n"))
			return
		}
	}()

	<-ctx.Done()

	err = client.Close()
	if err != nil {
		fmt.Println(err)
	}

	wg.Wait()
}
