package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer stop()

	listener, err := net.Listen("tcp", ":25565")
	if err != nil {
		fmt.Printf("Failed to listen on tcp: %v\n", err)
		os.Exit(-1)
	}
	defer listener.Close()
	fmt.Printf("Listening on %s\n", listener.Addr())
	go listen(listener)

	fmt.Println("\nCtrl+C to stop")
	<-ctx.Done()
	fmt.Println("\nStopping...")
	fmt.Println("Bye!")
}

func listen(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if conn == nil {
			break
		}
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			conn.Close()
			continue
		}
		conn.Write([]byte("Hello World\n"))
		fmt.Printf("Got connection %s\n", conn.RemoteAddr())
		go connListen(conn)
	}
}

func connListen(conn net.Conn) {
	defer conn.Close()
	var err error
	for {
		var buf []byte
		buf, err = connRead(conn)
		if err != nil {
			break
		}
		if len(buf) == 0 {
			err = fmt.Errorf("client sent no data")
			break
		}

		fmt.Printf(">> %s: %d bytes\n   %s\n", conn.RemoteAddr(), len(buf), string(buf))
	}

	fmt.Printf("Closed connection to %s: %v\n", conn.RemoteAddr(), err)
}

func connRead(conn net.Conn) (buf []byte, err error) {
	buf = make([]byte, 0)
	b := make([]byte, 1024)
	for {
		var n int
		n, err = conn.Read(b)
		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			break
		}
		if n == 0 {
			continue
		}
		fmt.Printf("reading... got %d bytes: %v\n", n, b[:n])
		buf = append(buf, b[:n]...)
		if bytes.HasSuffix(buf, []byte{'\r', '\n'}) {
			break
		}
	}
	buf = bytes.TrimSuffix(buf, []byte{'\r', '\n'})
	return buf, err
}
