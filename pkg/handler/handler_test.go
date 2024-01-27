package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// TestHealthHandler is a unit test function that tests the healthHandler function.
// It creates a new HTTP request to the "/health" endpoint and checks if the handler returns the expected status code and body.
func TestHealthHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `OK`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

// TestValidateHandler is a unit test function that tests the validateHandler function.
// It creates an AdmissionReview object with a NetworkPolicy object and sends a POST request to the /validate endpoint.
// The function checks if the handler returns the expected status code, admission response, kind, and API version.
func TestValidateHandler(t *testing.T) {
	admissionReview := admissionv1.AdmissionReview{
		Request: &admissionv1.AdmissionRequest{
			Kind: metav1.GroupVersionKind{
				Group:   "networking.k8s.io",
				Version: "v1",
				Kind:    "NetworkPolicy",
			},
			Operation: admissionv1.Create,
			Object: runtime.RawExtension{
				Raw: []byte(`{"apiVersion":"networking.k8s.io/v1","kind":"NetworkPolicy","metadata":{"name":"test-policy"},"spec":{"podSelector":{"matchLabels":{"app":"test"}},"ingress":[{"from":[{"namespaceSelector":{"matchLabels":{"name":"test-ns"}}}]}]}}`),
			},
		},
	}

	admissionReviewBytes, err := json.Marshal(admissionReview)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/validate", bytes.NewBuffer(admissionReviewBytes))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ValidateHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var admissionResponse admissionv1.AdmissionReview
	err = json.Unmarshal(rr.Body.Bytes(), &admissionResponse)
	if err != nil {
		t.Fatal(err)
	}

	expectedAllowed := admissionResponse.Response.Allowed
	if expectedAllowed {
		t.Errorf("handler returned unexpected admission response: got allowed, want denied")
	}

	expectedKind := "AdmissionReview"
	if admissionResponse.Kind != expectedKind {
		t.Errorf("handler returned unexpected admission response: got %v want %v",
			admissionResponse.Kind, expectedKind)
	}

	expectedAPIVersion := "admission.k8s.io/v1"
	if admissionResponse.APIVersion != expectedAPIVersion {
		t.Errorf("handler returned unexpected admission response: got %v want %v",
			admissionResponse.APIVersion, expectedAPIVersion)
	}
}

func BenchmarkValidateHandler(b *testing.B) {
	admissionReview := admissionv1.AdmissionReview{
		Request: &admissionv1.AdmissionRequest{
			Kind: metav1.GroupVersionKind{
				Group:   "networking.k8s.io",
				Version: "v1",
				Kind:    "NetworkPolicy",
			},
			Operation: admissionv1.Create,
			Object: runtime.RawExtension{
				Raw: []byte(`{"apiVersion":"networking.k8s.io/v1","kind":"NetworkPolicy","metadata":{"name":"test-policy"},"spec":{"podSelector":{"matchLabels":{"app":"test"}},"ingress":[{"from":[{"namespaceSelector":{"matchLabels":{"name":"test-ns"}}}]}]}}`),
			},
		},
	}

	admissionReviewBytes, err := json.Marshal(admissionReview)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		req, err := http.NewRequest("POST", "/validate", bytes.NewBuffer(admissionReviewBytes))
		if err != nil {
			b.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(ValidateHandler)
		handler.ServeHTTP(rr, req)
	}
}

func BenchmarkHealthHandler(b *testing.B) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		b.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthHandler)

	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(rr, req)
	}
}
