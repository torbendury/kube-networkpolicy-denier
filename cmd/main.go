package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	handler "github.com/torbendury/kube-networkpolicy-denier/pkg/handler"
	log "github.com/torbendury/kube-networkpolicy-denier/pkg/logging"
)

// main is the entry point of the program.
// It initializes the server with the provided certificate and key files,
// sets up the request handlers, and starts the server.
func main() {
	var server *http.Server

	idleConnectionsClosed := make(chan struct{})
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGTERM)
		<-sig
		log.Info("Received sigterm. Stopping server...")
		if err := server.Shutdown(context.Background()); err != nil {
			log.Error(fmt.Sprintf("HTTP Server Shutdown Error: %v", err))
		}
		close(idleConnectionsClosed)
	}()

	certFile := flag.String("cert", "server.crt", "Server certificate file location")
	keyFile := flag.String("key", "server.key", "Server key file location")
	respMsg := flag.String("response", "This webhook denies all NetworkPolicies", "The response message to send back to the client")
	flag.Parse()

	mux := http.NewServeMux()

	handler.RespMsg = *respMsg

	mux.Handle("/health", http.TimeoutHandler(http.HandlerFunc(handler.HealthHandler), 4*time.Second, "request timed out"))
	mux.Handle("/validate", http.TimeoutHandler(http.HandlerFunc(handler.ValidateHandler), 4*time.Second, "request timed out"))

	server = &http.Server{
		Addr:    ":8443",
		Handler: mux,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
		IdleTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
		ReadTimeout:       10 * time.Second,
	}

	log.Info("Server started on port 8443")
	if err := server.ListenAndServeTLS(*certFile, *keyFile); err != http.ErrServerClosed {
		log.Error(fmt.Sprintf("HTTP server ListenAndServeTLS: %v", err))
	}

	<-idleConnectionsClosed

	log.Info("Server stopped. Shutting down...")
}
