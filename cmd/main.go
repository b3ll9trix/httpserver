package main

import (
	"context"
	"httpserver/internal/app/handlers"
	"httpserver/internal/pkg/httpserver"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	wg.Add(1)
	go startServer(ctx, &wg)

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM)

	<-signalChannel
	log.Print("\n Program Interrupted. Gracefully Shutting down...")
	cancel()
	wg.Wait()
	log.Print("Shutdown Complete")
}

func startServer(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	addr := os.Getenv("SERVER_PORT")
	server := httpserver.NewServer(addr)
	server.Handle("/joke", true, httpserver.HandlerFunc(handlers.GetRandomJoke))
	log.Printf("Server listening on %s\n", addr)

	go func() {
		<-ctx.Done()
		log.Print("Saving State...")
		server.Cleanup()
		log.Print("Shutting down the httpserver...")
		server.Shutdown(ctx)

	}()
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Error Starting Server:", err)
	}

}
