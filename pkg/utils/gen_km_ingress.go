package utils

import (
	jianzhiuniquev1 "github.com/jianzhiunique/kafka-operator/pkg/apis/jianzhiunique/v1"
	corev1 "k8s.io/api/core/v1"
	v1beta12 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func NewKafkaManagerIngressForCR(cr *jianzhiuniquev1.Kafka) *v1beta12.Ingress {
	cr.Status.KafkaManagerUrl = cr.Name + cr.Spec.KafkaManagerHost
	paths := make([]v1beta12.HTTPIngressPath, 0)
	path := v1beta12.HTTPIngressPath{
		Path: "/",
		Backend: v1beta12.IngressBackend{
			ServicePort: intstr.IntOrString{
				IntVal: 9000,
			},
			ServiceName: "kfk-m-svc-" + cr.Name,
		},
	}
	paths = append(paths, path)

	rules := make([]v1beta12.IngressRule, 0)
	rule := v1beta12.IngressRule{
		Host: cr.Status.KafkaManagerUrl,
		IngressRuleValue: v1beta12.IngressRuleValue{
			HTTP: &v1beta12.HTTPIngressRuleValue{
				Paths: paths,
			},
		},
	}
	rules = append(rules, rule)

	port := corev1.ServicePort{Port: 9092}
	ports := make([]corev1.ServicePort, 0)
	ports = append(ports, port)
	return &v1beta12.Ingress{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "networking.k8s.io/v1beta1",
			Kind:       "Ingress",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "km-ingress-" + cr.Name,
			Namespace: cr.Namespace,
		},
		Spec: v1beta12.IngressSpec{
			Rules: rules,
		},
	}
}