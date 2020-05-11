package utils

import (
	v1 "github.com/jianzhiunique/kafka-operator/pkg/apis/jianzhiunique/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewExporterSvcForCR(cr *v1.Kafka) *corev1.Service {
	metrics := cr.Status.KafkaUrl + "-metrics"
	port := corev1.ServicePort{Port: 9308, Name: metrics + "-port"}
	ports := make([]corev1.ServicePort, 0)
	ports = append(ports, port)

	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kfk-e-svc-" + cr.Name,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"kfk-metrics": metrics,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: ports,
			Selector: map[string]string{
				"app": "kfk-exporter-" + cr.Name,
			},
		},
	}
}
