package utils

import (
	v1 "github.com/jianzhiunique/kafka-operator/pkg/apis/jianzhiunique/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func NewProxyForCR(cr *v1.Kafka) *appsv1.StatefulSet {
	var replica int32 = 2

	accessModes := make([]corev1.PersistentVolumeAccessMode, 0)
	accessModes = append(accessModes, corev1.ReadWriteOnce)
	pvc := make([]corev1.PersistentVolumeClaim, 0)
	pvc = append(pvc, corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kfk-mqp-data",
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: &cr.Spec.StorageClassName,
			Resources: corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(cr.Spec.ProxyDiskLimit),
				},
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(cr.Spec.ProxyDiskRequest),
				},
			},
			AccessModes: accessModes,
		},
	})

	containers := make([]corev1.Container, 0)
	ports := make([]corev1.ContainerPort, 0)
	ports = append(ports, corev1.ContainerPort{
		Name:          "kfk-mqp-port",
		ContainerPort: 8080,
		Protocol:      "TCP",
	})

	envs := make([]corev1.EnvVar, 0)
	envs = append(envs,
		corev1.EnvVar{
			Name:  "rabbit.address.host",
			Value: "kfk-svc-" + cr.Name,
		},
		corev1.EnvVar{
			Name:  "rabbit.address.port",
			Value: "5672",
		},
		corev1.EnvVar{
			Name:  "proxy.config.rabbitmq",
			Value: "true",
		},
		corev1.EnvVar{
			Name:  "proxy.config.kafka",
			Value: "false",
		},
		corev1.EnvVar{
			Name:  "logging.path",
			Value: "/data/mqp",
		},
	)
	vms := make([]corev1.VolumeMount, 0)
	vms = append(vms, corev1.VolumeMount{
		Name:      "kfk-mqp-data",
		MountPath: "/data/mqp",
	})
	healthCheck := corev1.Probe{
		Handler: corev1.Handler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					IntVal: 8080,
				},
			},
		},
		InitialDelaySeconds: 15,
		PeriodSeconds:       20,
	}
	c := corev1.Container{
		Name:           "kfk-mqp",
		Image:          "jianzhiunique/mqproxy:latest",
		Ports:          ports,
		Env:            envs,
		VolumeMounts:   vms,
		LivenessProbe:  &healthCheck,
		ReadinessProbe: &healthCheck,
		Resources: corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse(cr.Spec.ProxyMemoryLimit),
				corev1.ResourceCPU:    resource.MustParse(cr.Spec.ProxyCpuLimit),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse(cr.Spec.ProxyMemoryRequest),
				corev1.ResourceCPU:    resource.MustParse(cr.Spec.ProxyCpuRequest),
			},
		},
	}
	containers = append(containers, c)

	return &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kfk-mqp-sts-" + cr.Name,
			Namespace: cr.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:    &replica,
			ServiceName: "kfk-mqp-svc-" + cr.Name,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "kfk-mqp-" + cr.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "kfk-mqp-" + cr.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: containers,
				},
			},
			// for data store
			VolumeClaimTemplates: pvc,
		},
		Status: appsv1.StatefulSetStatus{},
	}
}
