package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("This webhook denies all NetworkPolicies"))
}

func main() {
	certFile := flag.String("cert", "server.crt", "Server certificate file location")
	keyFile := flag.String("key", "server.key", "Server key file location")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/validate", validateHandler)

	server := &http.Server{
		Addr:    ":8443",
		Handler: mux,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	log.Fatal(server.ListenAndServeTLS(*certFile, *keyFile))
}
