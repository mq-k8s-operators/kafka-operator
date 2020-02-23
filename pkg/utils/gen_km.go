package utils

import (
	jianzhiuniquev1 "github.com/jianzhiunique/kafka-operator/pkg/apis/jianzhiunique/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"math/rand"
	"time"
)

func NewKafkaManagerForCR(cr *jianzhiuniquev1.Kafka) *appsv1.Deployment {
	zkUrl := cr.Status.ZkUrl
	if cr.Status.KafkaManagerPassword == "" {
		cr.Status.KafkaManagerPassword = GetRandomString(16)
	}
	cr.Status.KafkaManagerUsername = "admin"

	cport := corev1.ContainerPort{ContainerPort: 9000}
	cports := make([]corev1.ContainerPort, 0)
	cports = append(cports, cport)

	envs := make([]corev1.EnvVar, 0)
	envs = append(envs,
		corev1.EnvVar{
			Name:  "ZK_HOSTS",
			Value: zkUrl,
		},
		corev1.EnvVar{
			Name:  "KAFKA_MANAGER_AUTH_ENABLED",
			Value: "true",
		},
		corev1.EnvVar{
			Name:  "KAFKA_MANAGER_USERNAME",
			Value: cr.Status.KafkaManagerUsername,
		},
		corev1.EnvVar{
			Name:  "KAFKA_MANAGER_PASSWORD",
			Value: cr.Status.KafkaManagerPassword,
		},
	)

	containers := make([]corev1.Container, 0)
	container := corev1.Container{
		Name:  "kfk-m-c-" + cr.Name,
		Image: cr.Spec.ManagerImage,
		Ports: cports,
		Env:   envs,
	}
	containers = append(containers, container)

	port := corev1.ServicePort{Port: 9092}
	ports := make([]corev1.ServicePort, 0)
	ports = append(ports, port)

	var replica int32
	replica = 1
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kfk-manager-" + cr.Name,
			Namespace: cr.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replica,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "kfk-m-" + cr.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "kfk-m-" + cr.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: containers,
				},
			},
		},
	}
}

func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
