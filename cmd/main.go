package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"

	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var logger *log.Logger

func init() {
	logger = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	logger.Println("Health check requested")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	logger.Println("Validation requested")

	// Read the admission review request
	admissionReview := admissionv1.AdmissionReview{}
	err := json.NewDecoder(r.Body).Decode(&admissionReview)
	if err != nil {
		logger.Println("Failed to decode admission review request:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

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

	logger.Println("Server started on port 8443")
	log.Fatal(server.ListenAndServeTLS(*certFile, *keyFile))
}
