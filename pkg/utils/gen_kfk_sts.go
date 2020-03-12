package utils

import (
	jianzhiuniquev1 "github.com/jianzhiunique/kafka-operator/pkg/apis/jianzhiunique/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strconv"
)

func NewStsForCR(cr *jianzhiuniquev1.Kafka) *appsv1.StatefulSet {
	svcName := "kfk-svc-" + cr.Name
	zkUrl := cr.Status.ZkUrl

	if cr.DeletionTimestamp.IsZero() {

	}

	accessModes := make([]corev1.PersistentVolumeAccessMode, 0)
	accessModes = append(accessModes, corev1.ReadWriteOnce)
	pvc := make([]corev1.PersistentVolumeClaim, 0)
	pvc = append(pvc, corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kfk-data",
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: &cr.Spec.StorageClassName,
			Resources: corev1.ResourceRequirements{
				Limits: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(cr.Spec.DiskLimit),
				},
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(cr.Spec.DiskRequest),
				},
			},
			AccessModes: accessModes,
		},
	})

	containers := make([]corev1.Container, 0)
	ports := make([]corev1.ContainerPort, 0)
	ports = append(ports, corev1.ContainerPort{
		Name:          "kfk-port",
		ContainerPort: 9092,
		Protocol:      "TCP",
	})
	ports = append(ports, corev1.ContainerPort{
		Name:          "jmx-port",
		ContainerPort: 9999,
		Protocol:      "TCP",
	})

	envs := make([]corev1.EnvVar, 0)
	envs = append(envs,
		corev1.EnvVar{
			Name:  "KAFKA_ZOOKEEPER_CONNECT",
			Value: zkUrl,
		},
		corev1.EnvVar{
			Name:  "BROKER_ID_COMMAND",
			Value: "hostname | awk -F'-' '{print $$4}'",
		},
		corev1.EnvVar{
			Name:  "KAFKA_AUTO_CREATE_TOPICS_ENABLE",
			Value: "false",
		},
		corev1.EnvVar{
			Name:  "KAFKA_DELETE_TOPIC_ENABLE",
			Value: "true",
		},
		corev1.EnvVar{
			Name:  "KAFKA_LISTENERS",
			Value: "PLAINTEXT://0.0.0.0:9092",
		},
		corev1.EnvVar{
			Name: "KAFKA_ADVERTISED_HOST_NAME",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.name",
				},
			},
		},
		corev1.EnvVar{
			Name: "MY_POD_NAMESPACE",
			ValueFrom: &corev1.EnvVarSource{
				FieldRef: &corev1.ObjectFieldSelector{
					FieldPath: "metadata.namespace",
				},
			},
		},
		corev1.EnvVar{
			Name:  "KAFKA_ADVERTISED_LISTENERS",
			Value: "PLAINTEXT://$(KAFKA_ADVERTISED_HOST_NAME)." + svcName + ".$(MY_POD_NAMESPACE).svc.cluster.local:9092",
		},
		corev1.EnvVar{
			Name:  "KAFKA_LOG_DIRS",
			Value: "/data/kafka",
		},
		corev1.EnvVar{
			Name:  "KAFKA_ZOOKEEPER_SESSION_TIMEOUT_MS",
			Value: "120000",
		},
		corev1.EnvVar{
			Name:  "KAFKA_AUTO_LEADER_REBANLANCE_ENABLE",
			Value: "false",
		}, corev1.EnvVar{
			Name:  "KAFKA_OFFSETS_RETENTION_MINUTES",
			Value: "14400",
		},
		corev1.EnvVar{
			Name:  "JMX_PORT",
			Value: "9999",
		},
		corev1.EnvVar{
			Name:  "KAFKA_NUM_PARTITIONS",
			Value: strconv.Itoa(int(cr.Spec.KafkaNumPartitions)),
		},
		corev1.EnvVar{
			Name:  "KAFKA_LOG_RETENTION_HOURS",
			Value: strconv.Itoa(int(cr.Spec.KafkaLogRetentionHours)),
		},
		corev1.EnvVar{
			Name:  "KAFKA_LOG_RETENTION_BYTES",
			Value: strconv.FormatInt(cr.Spec.KafkaLogRetentionBytes, 10),
		},
		corev1.EnvVar{
			Name:  "KAFKA_DEFAULT_REPLICATION_FACTOR",
			Value: strconv.Itoa(int(cr.Spec.KafkaDefaultReplicationFactor)),
		},
		corev1.EnvVar{
			Name:  "KAFKA_MESSAGE_MAX_BYTES",
			Value: strconv.FormatInt(cr.Spec.KafkaMessageMaxBytes, 10),
		},
		corev1.EnvVar{
			Name:  "KAFKA_COMPRESSION_TYPE",
			Value: cr.Spec.KafkaCompressionType,
		},
		corev1.EnvVar{
			Name:  "KAFKA_UNCLEAN_LEADER_ELECTION_ENABLE",
			Value: cr.Spec.KafkaUncleanLeaderElectionEnable,
		},
		corev1.EnvVar{
			Name:  "KAFKA_LOG_CLEANUP_POLICY",
			Value: cr.Spec.KafkaLogCleanupPolicy,
		},
		corev1.EnvVar{
			Name:  "KAFKA_LOG_MESSAGE_TIMESTAMP_TYPE",
			Value: cr.Spec.KafkaLogMessageTimestampType,
		},
		corev1.EnvVar{
			Name:  "KAFKA_NUM_NETWORK_THREADS",
			Value: getNumThreads(cr.Spec.Size),
		},
		corev1.EnvVar{
			Name:  "KAFKA_NUM_IO_THREADS",
			Value: getNumThreads(cr.Spec.Size),
		},
		corev1.EnvVar{
			Name:  "KAFKA_NUM_RECOVERY_THREADS_PER_DATA_DIR",
			Value: "2",
		},
		corev1.EnvVar{
			Name:  "KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR",
			Value: "3",
		},
		corev1.EnvVar{
			Name:  "KAFKA_NUM_REPLICA_FETCHERS",
			Value: "2",
		},
		corev1.EnvVar{
			Name:  "KAFKA_MIN_INSYNC_REPLICAS",
			Value: "2",
		},
		corev1.EnvVar{
			Name:  "KAFKA_GROUP_INITIAL_REBANLANCE_DELAY_MS",
			Value: "3000",
		},
		corev1.EnvVar{
			Name:  "KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR",
			Value: "3",
		},
		corev1.EnvVar{
			Name:  "KAFKA_HEAP_OPTS",
			Value: "-Xms" + strconv.Itoa(cr.Spec.KafkaJvmXms) + "g -Xmx" + strconv.Itoa(cr.Spec.KafkaJvmXmx) + "g -XX:MetaspaceSize=96m -XX:+UseG1GC -XX:MaxGCPauseMillis=20 -XX:InitiatingHeapOccupancyPercent=35 -XX:G1HeapRegionSize=16M -XX:MinMetaspaceFreeRatio=50 -XX:MaxMetaspaceFreeRatio=80 -XX:ParallelGCThreads=16",
		},
	)
	vms := make([]corev1.VolumeMount, 0)
	vms = append(vms, corev1.VolumeMount{
		Name:      "kfk-data",
		MountPath: "/data/kafka",
	})
	healthCheck := corev1.Probe{
		Handler: corev1.Handler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.IntOrString{
					IntVal: 9092,
				},
			},
		},
		InitialDelaySeconds: 15,
		PeriodSeconds:       20,
	}
	kfk := corev1.Container{
		Name:           "kafka",
		Image:          cr.Spec.Image,
		Ports:          ports,
		Env:            envs,
		VolumeMounts:   vms,
		LivenessProbe:  &healthCheck,
		ReadinessProbe: &healthCheck,
		Resources: corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse(cr.Spec.MemoryLimit),
				corev1.ResourceCPU:    resource.MustParse(cr.Spec.CpuLimit),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceMemory: resource.MustParse(cr.Spec.MemoryRequest),
				corev1.ResourceCPU:    resource.MustParse(cr.Spec.CpuRequest),
			},
		},
	}
	containers = append(containers, kfk)

	return &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "kfk-sts-" + cr.Name,
			Namespace: cr.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:    &cr.Spec.Size,
			ServiceName: svcName,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "kfk-pod-" + cr.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "kfk-pod-" + cr.Name,
					},
				},
				Spec: corev1.PodSpec{
					Containers: containers,
					//ServiceAccountName: "kafka-operator",
					/*
						Affinity: &corev1.Affinity{
							NodeAffinity:    &corev1.NodeAffinity{
								RequiredDuringSchedulingIgnoredDuringExecution:  &corev1.NodeSelector{
									NodeSelectorTerms: []corev1.NodeSelectorTerm{
										{
											MatchExpressions: []corev1.NodeSelectorRequirement{
												{
													Key: "nodegroup/kafka",
													Operator: corev1.NodeSelectorOpExists,
												},
											},
										},
									},
								},
							},
						},
					*/
				},
			},
			VolumeClaimTemplates: pvc,
		},
		Status: appsv1.StatefulSetStatus{},
	}
}

func getNumThreads(size int32) string {
	if size > 3 {
		return "64"
	} else {
		return "32"
	}
}
