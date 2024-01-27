// Handler package contains the HTTP handlers for the admission webhook.
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	admissionv1 "k8s.io/api/admission/v1"

	"github.com/torbendury/kube-networkpolicy-denier/pkg/admission"
	log "github.com/torbendury/kube-networkpolicy-denier/pkg/logging"
)

var (
	RespMsg string
)

// writeResponse writes the provided body to the response writer and sets the provided header.
func writeResponse(w http.ResponseWriter, header int, body []byte) error {
	w.WriteHeader(header)
	_, err := w.Write(body)
	return err
}

// HealthHandler handles the health check request.
// It returns a 200 OK response.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
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
			log.Error(fmt.Sprintf("Failed to write health check response: %v", err))
		}
		return
	case <-ctx.Done():
		log.Error("health check timed out")
		http.Error(w, "request timed out", http.StatusRequestTimeout)
	}
}

// ValidateHandler is a function that handles the validation of admission review requests.
// Parameters:
// - w: http.ResponseWriter - the response writer used to write the admission review response.
// - r: *http.Request - the HTTP request containing the admission review request.
// Returns: None
func ValidateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	done := make(chan error, 1)

	admissionReview := admissionv1.AdmissionReview{}
	err := json.NewDecoder(r.Body).Decode(&admissionReview)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to decode admission review request: %v", err))
		respErr := writeResponse(w, http.StatusBadRequest, []byte("Failed to decode admission review request"))
		if respErr != nil {
			log.Error(fmt.Sprintf("Failed to return error response: %v", err))
		}
		return
	}

	go func() {
		log.Info(fmt.Sprintf("%v name: %v namespace: %v operation: %v uid: %v", r.URL.Path, admissionReview.Request.Name, admissionReview.Request.Namespace, admissionReview.Request.Operation, admissionReview.Request.UID))

		// Create the admission review response
		admissionReviewResponse := admission.CreateAdmissionReviewResponse(&admissionReview, &RespMsg)

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&admissionReviewResponse)
		if err != nil {
			log.Error(fmt.Sprintf("Failed to encode admission review response: %v", err))
			respErr := writeResponse(w, http.StatusInternalServerError, []byte("Failed to encode admission review response"))
			if respErr != nil {
				log.Error(fmt.Sprintf("Failed to return error response: %v", err))
			}
			done <- err
		}

		done <- nil
	}()

	select {
	case err := <-done:
		if err != nil {
			log.Error(fmt.Sprintf("Failed to process admission review request: %v", err))
		}
		return
	case <-ctx.Done():
		log.Error("validation request timed out")
	}
}
