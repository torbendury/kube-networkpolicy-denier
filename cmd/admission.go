package main

import (
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// createAdmissionReviewResponse creates an admission review response.
func createAdmissionReviewResponse(admissionReview *admissionv1.AdmissionReview) admissionv1.AdmissionReview {
	admissionResponse := admissionv1.AdmissionResponse{
		Allowed: false,
		Result: &metav1.Status{
			Message: *respMsg,
		},
		UID: admissionReview.Request.UID,
	}

	admissionReviewResponse := admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "admission.k8s.io/v1",
			Kind:       "AdmissionReview",
		},
		Response: &admissionResponse,
	}
	return admissionReviewResponse
}
