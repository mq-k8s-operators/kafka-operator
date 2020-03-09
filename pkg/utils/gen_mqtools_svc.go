package utils

import (
	v1 "github.com/jianzhiunique/kafka-operator/pkg/apis/jianzhiunique/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewToolsSvcForCR(cr *v1.Kafka) *corev1.Service {
	port := corev1.ServicePort{Port: 8888, Name: "mqtools"}
	ports := make([]corev1.ServicePort, 0)
	ports = append(ports, port)

	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kfk-tools-svc-" + cr.Name,
			Namespace: cr.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: ports,
			Selector: map[string]string{
				"app": "kfk-tools-" + cr.Name,
			},
		},
	}
}
