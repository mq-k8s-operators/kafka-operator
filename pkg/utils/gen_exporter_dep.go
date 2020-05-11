package utils

import (
	v1 "github.com/jianzhiunique/kafka-operator/pkg/apis/jianzhiunique/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strings"
)

func NewExporterForCR(cr *v1.Kafka) *appsv1.Deployment {
	var replica int32 = 1

	limit := resource.MustParse(cr.Spec.ExporterDiskLimit)
	//this section should use persistent pv, such as ceph, to do
	pv := make([]corev1.Volume, 0)
	pv = append(pv, corev1.Volume{
		Name: "kfk-exporter-data",
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
		Name:          "kfk-e-port",
		ContainerPort: 9308,
		Protocol:      "TCP",
	})

	args := make([]string, 0)
	servers := strings.Split(cr.Status.KafkaUrlAll, ",")
	for i := 0; i < len(servers); i++ {
		args = append(args, "--kafka.server="+servers[i])
	}

	//TODO
	kafkaAndVersion := strings.Split(cr.Spec.Image, "-")
	args = append(args, "--kafka.version="+kafkaAndVersion[len(kafkaAndVersion)-1])

	vms := make([]corev1.VolumeMount, 0)
	vms = append(vms, corev1.VolumeMount{
		Name:      "kfk-exporter-data",
		MountPath: "/data",
	})
	healthCheck := corev1.Probe{
		Handler: corev1.Handler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					IntVal: 9308,
				},
			},
		},
		InitialDelaySeconds: 15,
		PeriodSeconds:       20,
	}
	c := corev1.Container{
		Name:           "kfk-exporter",
		Image:          cr.Spec.ExporterImage,
		Ports:          ports,
		Args:           args,
		VolumeMounts:   vms,
		LivenessProbe:  &healthCheck,
		ReadinessProbe: &healthCheck,
		Resources: corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse(cr.Spec.ExporterMemoryLimit),
				corev1.ResourceCPU:    resource.MustParse(cr.Spec.ExporterCpuLimit),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse(cr.Spec.ExporterMemoryRequest),
				corev1.ResourceCPU:    resource.MustParse(cr.Spec.ExporterCpuRequest),
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
			Name:      "kfk-exporter-dep-" + cr.Name,
			Namespace: cr.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replica,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "kfk-exporter-" + cr.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":     "kfk-exporter-" + cr.Name,
						"cluster": "kfk-" + cr.Namespace + "-" + cr.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: containers,
					Volumes:    pv,
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
