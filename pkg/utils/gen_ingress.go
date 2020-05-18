package utils

import (
	jianzhiuniquev1 "github.com/jianzhiunique/kafka-operator/pkg/apis/jianzhiunique/v1"
	v1beta12 "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func NewIngressForCRIfNotExists(cr *jianzhiuniquev1.Kafka) *v1beta12.Ingress {
	//通过比较资源的ns与ingress的配置ns来确定servicename要不要用external
	servicename := "kfk-m-svc-" + cr.Name
	if cr.Namespace != cr.Spec.IngressNamespace {
		servicename = "external-" + cr.Namespace + "-" + servicename
	}

	var annotations map[string]string
	annotations = make(map[string]string)
	annotations["nginx.ingress.kubernetes.io/rewrite-target"] = "/$1"

	//如果basepath不为空，rewrite kafka manager的path
	if cr.Spec.KafkaManagerBasePath != "" {
		annotations["nginx.ingress.kubernetes.io/configuration-snippet"] = "rewrite ^/(kfk-.*)$ " + cr.Spec.KafkaManagerBasePath + "/$1 break;"
	}

	paths := make([]v1beta12.HTTPIngressPath, 0)
	path := v1beta12.HTTPIngressPath{
		Path: "/(kfk-" + cr.Namespace + "-" + cr.Name + "/.*)",
		Backend: v1beta12.IngressBackend{
			ServicePort: intstr.IntOrString{
				IntVal: 9000,
			},
			ServiceName: servicename,
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
	//通过比较资源的ns与ingress的配置ns来确定servicename要不要用external
	servicename := "kfk-m-svc-" + cr.Name
	if cr.Namespace != cr.Spec.IngressNamespace {
		servicename = "external-" + cr.Namespace + "-" + servicename
	}

	if cr.Spec.KafkaManagerBasePath != "" {
		ingress.Annotations["nginx.ingress.kubernetes.io/configuration-snippet"] = "rewrite ^/(kfk-.*)$ " + cr.Spec.KafkaManagerBasePath + "/$1 break;"
	}

	path := v1beta12.HTTPIngressPath{
		Path: "/(kfk-" + cr.Namespace + "-" + cr.Name + "/.*)",
		Backend: v1beta12.IngressBackend{
			ServicePort: intstr.IntOrString{
				IntVal: 9000,
			},
			ServiceName: servicename,
		},
	}

	var found bool

	for _, path := range ingress.Spec.Rules[0].IngressRuleValue.HTTP.Paths {
		if path.Path == "/(kfk-"+cr.Namespace+"-"+cr.Name+"/.*)" {
			found = true
			break
		}
	}

	if !found {
		ingress.Spec.Rules[0].IngressRuleValue.HTTP.Paths = append(ingress.Spec.Rules[0].IngressRuleValue.HTTP.Paths, path)
	}
}

func DeleteKafkaManagerPathFromIngress(cr *jianzhiuniquev1.Kafka, ingress *v1beta12.Ingress) {

	paths := ingress.Spec.Rules[0].IngressRuleValue.HTTP.Paths

	for index, path := range ingress.Spec.Rules[0].IngressRuleValue.HTTP.Paths {
		if path.Path == "/(kfk-"+cr.Namespace+"-"+cr.Name+"/.*)" {
			ingress.Spec.Rules[0].IngressRuleValue.HTTP.Paths = append(paths[:index], paths[index+1:]...)
			break
		}
	}
}

func AppendKafkaToolsPathToIngress(cr *jianzhiuniquev1.Kafka, ingress *v1beta12.Ingress) {
	//通过比较资源的ns与ingress的配置ns来确定servicename要不要用external
	servicename := "kfk-tools-svc-" + cr.Name
	if cr.Namespace != cr.Spec.IngressNamespace {
		servicename = "external-" + cr.Namespace + "-" + servicename
	}

	path := v1beta12.HTTPIngressPath{
		Path: cr.Status.KafkaToolsPath + "(.*)",
		Backend: v1beta12.IngressBackend{
			ServicePort: intstr.IntOrString{
				IntVal: 8888,
			},
			ServiceName: servicename,
		},
	}

	var found bool

	for _, path := range ingress.Spec.Rules[0].IngressRuleValue.HTTP.Paths {
		if path.Path == cr.Status.KafkaToolsPath+"(.*)" {
			found = true
			break
		}
	}

	if !found {
		ingress.Spec.Rules[0].IngressRuleValue.HTTP.Paths = append(ingress.Spec.Rules[0].IngressRuleValue.HTTP.Paths, path)
	}
}

func DeleteKafkaToolsPathFromIngress(cr *jianzhiuniquev1.Kafka, ingress *v1beta12.Ingress) {

	paths := ingress.Spec.Rules[0].IngressRuleValue.HTTP.Paths

	for index, path := range ingress.Spec.Rules[0].IngressRuleValue.HTTP.Paths {
		if path.Path == cr.Status.KafkaToolsPath+"(.*)" {
			ingress.Spec.Rules[0].IngressRuleValue.HTTP.Paths = append(paths[:index], paths[index+1:]...)
			break
		}
	}
}
