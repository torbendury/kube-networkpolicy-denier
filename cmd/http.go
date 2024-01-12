package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	admissionv1 "k8s.io/api/admission/v1"
)

// writeResponse writes the provided body to the response writer and sets the provided header.
func writeResponse(w http.ResponseWriter, header int, body []byte) error {
	w.WriteHeader(header)
	_, err := w.Write(body)
	return err
}

// healthHandler handles the health check request.
// It returns a 200 OK response.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		err := writeResponse(w, http.StatusOK, []byte("OK"))
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			errorLogger.Println("Failed to write health check response:", err)
		}
		return
	case <-ctx.Done():
		errorLogger.Println("health check timed out")
		http.Error(w, "request timed out", http.StatusRequestTimeout)
	}
}

// validateHandler is a function that handles the validation of admission review requests.
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
		respErr := writeResponse(w, http.StatusBadRequest, []byte("Failed to decode admission review request"))
		if respErr != nil {
			errorLogger.Println("Failed to return error response:", err)
		}
		return
	}

	go func() {
		infoLogger.Println(r.URL.Path, " name: ", admissionReview.Request.Name, " namespace: ", admissionReview.Request.Namespace, " operation: ", admissionReview.Request.Operation, " uid: ", admissionReview.Request.UID)

		// Create the admission review response
		admissionReviewResponse := createAdmissionReviewResponse(&admissionReview)

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&admissionReviewResponse)
		if err != nil {
			errorLogger.Println("Failed to encode admission review response:", err)
			respErr := writeResponse(w, http.StatusInternalServerError, []byte("Failed to encode admission review response"))
			if respErr != nil {
				errorLogger.Println("Failed to return error response:", err)
			}
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
