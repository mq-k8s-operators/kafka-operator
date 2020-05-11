package utils

import (
	jianzhiuniquev1 "github.com/jianzhiunique/kafka-operator/pkg/apis/jianzhiunique/v1"
	v1beta12 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func NewIngressForCRIfNotExists(cr *jianzhiuniquev1.Kafka) *v1beta12.Ingress {
	var annotations map[string]string
	annotations = make(map[string]string)
	annotations["nginx.ingress.kubernetes.io/rewrite-target"] = "/$1"

	if cr.Spec.KafkaManagerBasePath != "" {
		annotations["nginx.ingress.kubernetes.io/configuration-snippet"] = "rewrite ^/(kfk-.*)$ " + cr.Spec.KafkaManagerBasePath + "/$1 break;"
	}

	paths := make([]v1beta12.HTTPIngressPath, 0)
	path := v1beta12.HTTPIngressPath{
		Path: cr.Status.KafkaManagerPath + "(.*)",
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
		Host: cr.Spec.KafkaManagerHost,
		IngressRuleValue: v1beta12.IngressRuleValue{
			HTTP: &v1beta12.HTTPIngressRuleValue{
				Paths: paths,
			},
		},
	}
	rules = append(rules, rule)

	return &v1beta12.Ingress{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "networking.k8s.io/v1beta1",
			Kind:       "Ingress",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        "mq-ingress",
			Namespace:   cr.Spec.IngressNamespace,
			Annotations: annotations,
		},
		Spec: v1beta12.IngressSpec{
			Rules: rules,
		},
	}
}

func AppendKafkaManagerPathToIngress(cr *jianzhiuniquev1.Kafka, ingress *v1beta12.Ingress) {
	if cr.Spec.KafkaManagerBasePath != "" {
		ingress.Annotations["nginx.ingress.kubernetes.io/configuration-snippet"] = "rewrite ^/(kfk-.*)$ " + cr.Spec.KafkaManagerBasePath + "/$1 break;"
	}

	path := v1beta12.HTTPIngressPath{
		Path: "/(kfk-" + cr.Namespace + "-" + cr.Name + "/.*)",
		Backend: v1beta12.IngressBackend{
			ServicePort: intstr.IntOrString{
				IntVal: 9000,
			},
			ServiceName: "kfk-m-svc-" + cr.Name,
		},
	}

	ingress.Spec.Rules[0].IngressRuleValue.HTTP.Paths = append(ingress.Spec.Rules[0].IngressRuleValue.HTTP.Paths, path)
}

func AppendKafkaToolsPathToIngress(cr *jianzhiuniquev1.Kafka, ingress *v1beta12.Ingress) {
	path := v1beta12.HTTPIngressPath{
		Path: cr.Status.KafkaToolsPath + "(.*)",
		Backend: v1beta12.IngressBackend{
			ServicePort: intstr.IntOrString{
				IntVal: 8888,
			},
			ServiceName: "kfk-tools-svc-" + cr.Name,
		},
	}

	ingress.Spec.Rules[0].IngressRuleValue.HTTP.Paths = append(ingress.Spec.Rules[0].IngressRuleValue.HTTP.Paths, path)
}
