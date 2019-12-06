package utils

import (
	jianzhiuniquev1 "github.com/jianzhiunique/kafka-operator/pkg/apis/jianzhiunique/v1"
)

func CheckCR(cr *jianzhiuniquev1.Kafka) bool {
	/*
	  size: 3
	  image: wurstmeister/kafka:2.11-0.11.0.3
	  disk_limit: 10Gi
	  disk_request: 1Gi
	  storage_class_name: standard
	  kafka_manager_host: ".km.com"
	  zk_size: 3
	  zk_disk_limit: 10Gi
	  zk_disk_request: 1Gi
	*/
	var changed bool
	changed = false

	if cr.Spec.Size == 0 || cr.Spec.Size < 3 {
		cr.Spec.Size = 3
		changed = true
	}

	if cr.Spec.Image == "" {
		cr.Spec.Image = "wurstmeister/kafka:2.11-0.11.0.3"
		changed = true
	}

	if cr.Spec.DiskLimit == "" {
		cr.Spec.DiskLimit = "500Gi"
		changed = true
	}

	if cr.Spec.DiskRequest == "" {
		cr.Spec.DiskRequest = "1Gi"
		changed = true
	}

	if cr.Spec.MemoryLimit == "" {
		cr.Spec.MemoryLimit = "32Gi"
		changed = true
	}

	if cr.Spec.MemoryRequest == "" {
		cr.Spec.MemoryRequest = "1Gi"
		changed = true
	}

	if cr.Spec.CpuLimit == "" {
		cr.Spec.CpuLimit = "4000m"
		changed = true
	}

	if cr.Spec.CpuRequest == "" {
		cr.Spec.CpuRequest = "500m"
		changed = true
	}

	if cr.Spec.KafkaManagerHost == "" {
		cr.Spec.KafkaManagerHost = ".kfk.cloudmq.com"
		changed = true
	}

	if cr.Spec.ZkSize == 0 {
		cr.Spec.ZkSize = 3
		changed = true
	}

	if cr.Spec.ZkDiskLimit == "" {
		cr.Spec.ZkDiskLimit = "500Gi"
		changed = true
	}

	if cr.Spec.ZkDiskRequest == "" {
		cr.Spec.ZkDiskRequest = "1Gi"
		changed = true
	}

	/*
	  default_partitions: 3
	  log_hours: 168
	  log_bytes: -1
	  replication_factor: 2
	  message_max_bytes: 1073741824
	  compression_type: producer
	  unclean_election: false
	  cleanup_policy: delete
	  message_timestamp_type: CreateTime
	*/

	if cr.Spec.KafkaNumPartitions == 0 {
		cr.Spec.KafkaNumPartitions = 3
		changed = true
	}

	if cr.Spec.KafkaLogRetentionHours == 0 {
		cr.Spec.KafkaLogRetentionHours = 168
		changed = true
	}

	if cr.Spec.KafkaLogRetentionBytes == 0 {
		cr.Spec.KafkaLogRetentionBytes = -1
		changed = true
	}

	if cr.Spec.KafkaDefaultReplicationFactor == 0 {
		cr.Spec.KafkaDefaultReplicationFactor = 2
		changed = true
	}

	if cr.Spec.KafkaMessageMaxBytes == 0 {
		cr.Spec.KafkaMessageMaxBytes = 1073741824
		changed = true
	}

	if cr.Spec.KafkaCompressionType == "" {
		cr.Spec.KafkaCompressionType = "producer"
		changed = true
	}

	if cr.Spec.KafkaLogCleanupPolicy == "" {
		cr.Spec.KafkaLogCleanupPolicy = "delete"
		changed = true
	}

	if cr.Spec.KafkaLogMessageTimestampType == "" {
		cr.Spec.KafkaLogMessageTimestampType = "CreateTime"
		changed = true
	}

	if cr.Spec.KafkaJvmXms == 0 {
		cr.Spec.KafkaJvmXms = 1
	}

	if cr.Spec.KafkaJvmXmx == 0 {
		cr.Spec.KafkaJvmXmx = 1
	}

	return changed
}
