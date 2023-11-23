package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var logger *log.Logger

// init initializes the logger variable with a new instance of log.Logger.
// It sets the output to os.Stdout and prefixes log messages with "[INFO] ".
// The logger is configured to include the date, time, and short file name in log messages.
func init() {
	logger = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
}

// healthHandler handles the health check request.
// It logs the request and returns a 200 OK response.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	logger.Println(r.URL.Path)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// validateHandler is a function that handles the validation of admission review requests.
// It reads the admission review request, denies all NetworkPolicies, and sends the admission review response.
// Parameters:
// - w: http.ResponseWriter - the response writer used to write the admission review response.
// - r: *http.Request - the HTTP request containing the admission review request.
// Returns: None
func validateHandler(w http.ResponseWriter, r *http.Request) {
	// Read the admission review request
	admissionReview := admissionv1.AdmissionReview{}
	err := json.NewDecoder(r.Body).Decode(&admissionReview)
	if err != nil {
		logger.Println("Failed to decode admission review request:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.Println(r.URL.Path, " name: ", admissionReview.Request.Name, " namespace: ", admissionReview.Request.Namespace, " operation: ", admissionReview.Request.Operation, " uid: ", admissionReview.Request.UID)

	// Get the UID from the AdmissionRequest
	uid := admissionReview.Request.UID

	// Create the admission review response
	admissionResponse := admissionv1.AdmissionResponse{
		Allowed: false,
		Result: &metav1.Status{
			Message: "This webhook denies all NetworkPolicies",
		},
		UID: uid,
	}

	// Create the admission review response review
	admissionReviewResponse := admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "admission.k8s.io/v1",
			Kind:       "AdmissionReview",
		},
		Response: &admissionResponse,
	}

	// Write the admission review response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&admissionReviewResponse)
	if err != nil {
		logger.Println("Failed to encode admission review response:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
		logger.Println("Received sigterm. Stopping server...")
		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatalf("HTTP Server Shutdown Error: %v", err)
		}
		close(idleConnectionsClosed)
	}()

	certFile := flag.String("cert", "server.crt", "Server certificate file location")
	keyFile := flag.String("key", "server.key", "Server key file location")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/validate", validateHandler)

	server = &http.Server{
		Addr:    ":8443",
		Handler: mux,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	logger.Println("Server started on port 8443")
	if err := server.ListenAndServeTLS(*certFile, *keyFile); err != http.ErrServerClosed {
		logger.Println("HTTP server ListenAndServeTLS:", err)
	}

	<-idleConnectionsClosed

	logger.Println("Server stopped. Shutting down...")
}
