package utils

import corev1 "k8s.io/api/core/v1"

var ImagePullSecret []corev1.LocalObjectReference

func init() {
	ImagePullSecret = make([]corev1.LocalObjectReference, 0)
	ImagePullSecret = append(ImagePullSecret, corev1.LocalObjectReference{Name: ""})
}
