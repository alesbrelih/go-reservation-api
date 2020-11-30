package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alesbrelih/go-reservation-api/pkg/myutil"
	"github.com/alesbrelih/go-reservation-api/router"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	port := myutil.GetEnvOrDefault("APPLICATION_PORT", "8080")

	mux := router.InitializeRouter()

	l := log.New(os.Stdout, "reservations", log.LstdFlags)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ErrorLog:     l,
		ReadTimeout:  2 * time.Second,   // max time to read request
		WriteTimeout: 3 * time.Second,   // max time to write response
		IdleTimeout:  120 * time.Second, // max time for TPC keepalive conns
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
	}

	// open server in  nonblocking way
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			l.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
	)

	<-signalChan
	log.Print("Gracefull shutting down...\n")

	gracefullCtx, cancelShutdown := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancelShutdown()

	err := server.Shutdown(gracefullCtx)
	if err != nil {
		log.Printf("Shutdown error: %v\n", err)
		os.Exit(1)
	} else {
		log.Printf("Gracefully stopped\n")
	}

	cancel()

	defer os.Exit(0)
	return
}
