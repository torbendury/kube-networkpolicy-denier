package admission

import (
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateAdmissionReviewResponse creates an admission review response.
func CreateAdmissionReviewResponse(admissionReview *admissionv1.AdmissionReview, respMsg *string) admissionv1.AdmissionReview {
	return admissionv1.AdmissionReview{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "admission.k8s.io/v1",
			Kind:       "AdmissionReview",
		},
		Response: &admissionv1.AdmissionResponse{
			Allowed: false,
			Result: &metav1.Status{
				Message: *respMsg,
			},
			UID: admissionReview.Request.UID,
		},
	}
}
