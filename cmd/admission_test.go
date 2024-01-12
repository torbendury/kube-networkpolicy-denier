package main

import (
	"testing"

	admissionv1 "k8s.io/api/admission/v1"
)

func TestCreateAdmissionReviewResponse(t *testing.T) {
	admissionReview := admissionv1.AdmissionReview{
		Request: &admissionv1.AdmissionRequest{
			UID: "1234",
		},
	}
	admissionReviewResponse := createAdmissionReviewResponse(&admissionReview)
	if admissionReviewResponse.Response.UID != "1234" {
		t.Errorf("Expected UID to be %s, got %s", "1234", admissionReviewResponse.Response.UID)
	}
	if admissionReviewResponse.Response.Allowed != false {
		t.Errorf("Expected Allowed to be %t, got %t", false, admissionReviewResponse.Response.Allowed)
	}
	if admissionReviewResponse.Response.Result.Message != *respMsg {
		t.Errorf("Expected Result.Message to be %s, got %s", *respMsg, admissionReviewResponse.Response.Result.Message)
	}
}
