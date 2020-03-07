package utils

import (
	jianzhiuniquev1 "github.com/jianzhiunique/kafka-operator/pkg/apis/jianzhiunique/v1"
	"github.com/pravega/zookeeper-operator/pkg/apis/zookeeper/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewZkForCR(cr *jianzhiuniquev1.Kafka) *v1beta1.ZookeeperCluster {

	return &v1beta1.ZookeeperCluster{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ZookeeperCluster",
			APIVersion: "zookeeper.pravega.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kfk-zk-" + cr.Name,
			Namespace: cr.Namespace,
		},
		Spec: v1beta1.ZookeeperClusterSpec{
			Image: v1beta1.ContainerImage{
				Repository: cr.Spec.ZkImage,
				Tag:        cr.Spec.ZkVersion,
			},
			Replicas: cr.Spec.ZkSize,
			Persistence: &v1beta1.Persistence{
				VolumeReclaimPolicy: v1beta1.VolumeReclaimPolicyDelete,
				PersistentVolumeClaimSpec: corev1.PersistentVolumeClaimSpec{
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							corev1.ResourceStorage: resource.MustParse(cr.Spec.ZkDiskLimit),
						},
						Requests: corev1.ResourceList{
							corev1.ResourceStorage: resource.MustParse(cr.Spec.ZkDiskRequest),
						},
					},
					StorageClassName: &cr.Spec.StorageClassName,
				},
			},
			Pod: v1beta1.PodPolicy{
				Resources: corev1.ResourceRequirements{
					Limits: corev1.ResourceList{
						corev1.ResourceMemory: resource.MustParse(cr.Spec.ZkMemoryLimit),
						corev1.ResourceCPU:    resource.MustParse(cr.Spec.ZkCpuLimit),
					},
					Requests: corev1.ResourceList{
						corev1.ResourceMemory: resource.MustParse(cr.Spec.ZkMemoryRequest),
						corev1.ResourceCPU:    resource.MustParse(cr.Spec.ZkCpuRequest),
					},
				},
			},
		},
	}
}
