package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer stop()

	server, err := StartServer(25565)
	if err != nil {
		fmt.Printf("Failed to start server: %v", err)
		os.Exit(-1)
	}
	defer server.Stop()

	fmt.Println("\nCtrl+C to stop")
	<-ctx.Done()
	fmt.Println("\nStopping...")
	server.Stop()
	fmt.Println("Bye!")
}
