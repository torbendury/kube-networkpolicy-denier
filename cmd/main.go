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
	"time"

	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var infoLogger *log.Logger
var errorLogger *log.Logger
var respMsg *string

// init initializes the infoLogger variable with a new instance of log.Logger.
// It sets the output to os.Stdout and prefixes log messages with "[INFO] ".
// The infoLogger is configured to include the date, time, and short file name in log messages.
func init() {
	infoLogger = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}

// healthHandler handles the health check request.
// It returns a 200 OK response.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		done <- nil
	}()

	select {
	case <-done:
		return
	case <-ctx.Done():
		errorLogger.Println("health check timed out")
		http.Error(w, "request timed out", http.StatusRequestTimeout)
	}
}

// validateHandler is a function that handles the validation of admission review requests.
// It reads the admission review request, denies all NetworkPolicies, and sends the admission review response.
// Parameters:
// - w: http.ResponseWriter - the response writer used to write the admission review response.
// - r: *http.Request - the HTTP request containing the admission review request.
// Returns: None
func validateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	done := make(chan error, 1)

	admissionReview := admissionv1.AdmissionReview{}
	err := json.NewDecoder(r.Body).Decode(&admissionReview)
	if err != nil {
		errorLogger.Println("Failed to decode admission review request:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	go func() {
		// Read the admission review request

		infoLogger.Println(r.URL.Path, " name: ", admissionReview.Request.Name, " namespace: ", admissionReview.Request.Namespace, " operation: ", admissionReview.Request.Operation, " uid: ", admissionReview.Request.UID)

		// Get the UID from the AdmissionRequest
		uid := admissionReview.Request.UID

		// Create the admission review response
		admissionResponse := admissionv1.AdmissionResponse{
			Allowed: false,
			Result: &metav1.Status{
				Message: *respMsg,
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
			errorLogger.Println("Failed to encode admission review response:", err)
			w.WriteHeader(http.StatusInternalServerError)
			done <- err
		}

		done <- nil
	}()

	select {
	case err := <-done:
		if err != nil {
			errorLogger.Println("Failed to process admission review request:", err)
		}
		return
	case <-ctx.Done():
		errorLogger.Println("validation request timed out")
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
