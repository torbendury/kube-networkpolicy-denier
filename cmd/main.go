package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var infoLogger *log.Logger
var errorLogger *log.Logger
var respMsg *string

// init initializes the infoLogger variable with a new instance of log.Logger.
// It sets the output to os.Stdout and prefixes log messages with "[INFO] ".
// The infoLogger is configured to include the date, time, and short file name in log messages.
func init() {
	defaultResp := "This webhook denies all NetworkPolicies"
	respMsg = &defaultResp
	infoLogger = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}

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
		infoLogger.Println("Received sigterm. Stopping server...")
		if err := server.Shutdown(context.Background()); err != nil {
			errorLogger.Printf("HTTP Server Shutdown Error: %v\n", err)
		}
		close(idleConnectionsClosed)
	}()

	certFile := flag.String("cert", "server.crt", "Server certificate file location")
	keyFile := flag.String("key", "server.key", "Server key file location")
	respMsg = flag.String("response", "This webhook denies all NetworkPolicies", "The response message to send back to the client")
	flag.Parse()

	mux := http.NewServeMux()

	mux.Handle("/health", http.TimeoutHandler(http.HandlerFunc(healthHandler), 4*time.Second, "request timed out"))
	mux.Handle("/validate", http.TimeoutHandler(http.HandlerFunc(validateHandler), 4*time.Second, "request timed out"))

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

	infoLogger.Println("Server started on port 8443")
	if err := server.ListenAndServeTLS(*certFile, *keyFile); err != http.ErrServerClosed {
		errorLogger.Println("HTTP server ListenAndServeTLS:", err)
	}

	<-idleConnectionsClosed

	infoLogger.Println("Server stopped. Shutting down...")
}
