package utils

import (
	v1 "github.com/jianzhiunique/kafka-operator/pkg/apis/jianzhiunique/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func NewToolsForCR(cr *v1.Kafka) *appsv1.Deployment {
	var replica int32 = 1

	limit := resource.MustParse(cr.Spec.ToolsDiskLimit)
	//this section should use persistent pv, such as ceph, to do
	pv := make([]corev1.Volume, 0)
	pv = append(pv, corev1.Volume{
		Name: "kfk-tools-data",
		VolumeSource: corev1.VolumeSource{
			EmptyDir: &corev1.EmptyDirVolumeSource{
				Medium:    "",
				SizeLimit: &limit,
			},
		},
	})

	containers := make([]corev1.Container, 0)
	ports := make([]corev1.ContainerPort, 0)
	ports = append(ports, corev1.ContainerPort{
		Name:          "kfk-tools-port",
		ContainerPort: 8888,
		Protocol:      "TCP",
	})

	envs := make([]corev1.EnvVar, 0)
	envs = append(envs,
		corev1.EnvVar{
			Name:  "kafka.config.bootstrap.servers",
			Value: cr.Status.KafkaUrlAll,
		},
		corev1.EnvVar{
			Name:  "kafka.config.zookeeper",
			Value: cr.Status.ZkUrlAll,
		},
		corev1.EnvVar{
			Name:  "server.port",
			Value: "8888",
		},
		corev1.EnvVar{
			Name:  "kafka.config.name",
			Value: cr.Namespace + "-" + cr.Name,
		},
		corev1.EnvVar{
			Name:  "kafka.config.admin_ding_url",
			Value: cr.Spec.ToolsAdminDingUrl,
		},
		corev1.EnvVar{
			Name:  "kafka.config.localdb",
			Value: "/data/kafka-localdb-",
		},
	)
	vms := make([]corev1.VolumeMount, 0)
	vms = append(vms, corev1.VolumeMount{
		Name:      "kfk-tools-data",
		MountPath: "/data",
	})
	healthCheck := corev1.Probe{
		Handler: corev1.Handler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					IntVal: 8888,
				},
			},
		},
		InitialDelaySeconds: 15,
		PeriodSeconds:       20,
	}
	c := corev1.Container{
		Name:           "kfk-tools",
		Image:          cr.Spec.ToolsImage,
		Ports:          ports,
		Env:            envs,
		VolumeMounts:   vms,
		LivenessProbe:  &healthCheck,
		ReadinessProbe: &healthCheck,
		Resources: corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse(cr.Spec.ToolsMemoryLimit),
				corev1.ResourceCPU:    resource.MustParse(cr.Spec.ToolsCpuLimit),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse(cr.Spec.ToolsMemoryRequest),
				corev1.ResourceCPU:    resource.MustParse(cr.Spec.ToolsCpuRequest),
			},
		},
	}
	containers = append(containers, c)

	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kfk-tools-sts-" + cr.Name,
			Namespace: cr.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replica,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "kfk-tools-" + cr.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "kfk-tools-" + cr.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: containers,
					Volumes:    pv,
					//ServiceAccountName: "kafka-operator",
				},
			},
			Strategy:                appsv1.DeploymentStrategy{},
			MinReadySeconds:         0,
			RevisionHistoryLimit:    nil,
			Paused:                  false,
			ProgressDeadlineSeconds: nil,
		},
	}
}
