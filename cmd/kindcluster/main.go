package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// The main function is the entry point for the application
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure cancel is called on exit

	// Create a channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	// Create a new application
	app := newApplication("zpool")

	op, err := NewPGOperator("persistence.db")
	if err != nil {
		app.logger.Error("Cannot create PG Operator")
	}

	go op.reconcile(ctx)

	app.logger.Info("Starting server on :8080")

	srv := &http.Server{Addr: ":8080", Handler: app.routes()}

	// Start HTTP server in a goroutine
	go func() {
		fmt.Println("HTTP server listening on port 8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.logger.Error("Cannot Start http listening")
		}
	}()

	// Wait for an OS signal
	sig := <-sigChan
	fmt.Println("\nReceived OS signal:", sig)

	// Gracefully shut down HTTP server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("HTTP server shutdown error: %v\n", err)
	} else {
		fmt.Println("HTTP server shut down gracefully.")
	}
	// Notify workers to stop
	cancel()

	// Wait a bit for cleanup before exiting
	time.Sleep(2 * time.Second)
	fmt.Println("Main process exiting.")

}
