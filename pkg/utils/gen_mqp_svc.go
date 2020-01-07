package utils

import (
	v1 "github.com/jianzhiunique/kafka-operator/pkg/apis/jianzhiunique/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewMqpSvcForCR(cr *v1.Kafka) *corev1.Service {
	//A headless service must be used to control network identity of the pods (their hostnames), which in turn affect RabbitMQ node names.
	//see: https://www.rabbitmq.com/cluster-formation.html#peer-discovery-k8s
	//In fact, if service is not headless, the network identity of the pods can be resolve by specify ServiceName field on statefulset
	port := corev1.ServicePort{Port: 8080, Name: "mqp"}
	ports := make([]corev1.ServicePort, 0)
	ports = append(ports, port)
	cr.Status.KafkaProxyUrl = "kfk-mqp-svc-" + cr.Name + ":8080"

	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kfk-mqp-svc-" + cr.Name,
			Namespace: cr.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Ports: ports,
			Selector: map[string]string{
				"app": "kfk-mqp-" + cr.Name,
			},
		},
	}
}
