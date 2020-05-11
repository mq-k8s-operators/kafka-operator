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

	if cr.Spec.ZkImage == "" {
		cr.Spec.ZkImage = "pravega/zookeeper"
		changed = true
	}

	if cr.Spec.ZkVersion == "" {
		cr.Spec.ZkVersion = "0.2.4"
		changed = true
	}

	if cr.Spec.ManagerImage == "" {
		cr.Spec.ManagerImage = "registry.cn-hangzhou.aliyuncs.com/jianzhiunique/kafka-manager:latest"
		changed = true
	}

	if cr.Spec.ProxyImage == "" {
		cr.Spec.ProxyImage = "jianzhiunique/mqproxy:latest"
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
		cr.Spec.KafkaManagerHost = "kfk.cloudmq.com"
		changed = true
	}

	/*
		if cr.Spec.KafkaManagerHostAlias == "" {
			cr.Spec.KafkaManagerHostAlias = "kfk.cloudmq.com/cloudmq"
			changed = true
		}

		if cr.Spec.KafkaManagerBasePath == "" {
			cr.Spec.KafkaManagerBasePath = "/cloudmq"
			changed = true
		}
	*/

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

	if cr.Spec.KafkaUncleanLeaderElectionEnable == "" {
		cr.Spec.KafkaUncleanLeaderElectionEnable = "false"
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

	if cr.Spec.ProxyDiskLimit == "" {
		cr.Spec.ProxyDiskLimit = "10Gi"
		changed = true
	}

	if cr.Spec.ProxyDiskRequest == "" {
		cr.Spec.ProxyDiskRequest = "1Gi"
		changed = true
	}

	if cr.Spec.ProxyMemoryLimit == "" {
		cr.Spec.ProxyMemoryLimit = "2Gi"
		changed = true
	}

	if cr.Spec.ProxyMemoryRequest == "" {
		cr.Spec.ProxyMemoryRequest = "1Gi"
		changed = true
	}

	if cr.Spec.ProxyCpuLimit == "" {
		cr.Spec.ProxyCpuLimit = "2000m"
		changed = true
	}

	if cr.Spec.ProxyCpuRequest == "" {
		cr.Spec.ProxyCpuRequest = "500m"
		changed = true
	}

	if cr.Spec.ZkMemoryLimit == "" {
		cr.Spec.ZkMemoryLimit = "2Gi"
		changed = true
	}

	if cr.Spec.ZkMemoryRequest == "" {
		cr.Spec.ZkMemoryRequest = "1Gi"
		changed = true
	}

	if cr.Spec.ZkCpuLimit == "" {
		cr.Spec.ZkCpuLimit = "2000m"
		changed = true
	}

	if cr.Spec.ZkCpuRequest == "" {
		cr.Spec.ZkCpuRequest = "500m"
		changed = true
	}

	if cr.Spec.KafkaManagerMemoryLimit == "" {
		cr.Spec.KafkaManagerMemoryLimit = "2Gi"
		changed = true
	}

	if cr.Spec.KafkaManagerMemoryRequest == "" {
		cr.Spec.KafkaManagerMemoryRequest = "1Gi"
		changed = true
	}

	if cr.Spec.KafkaManagerCpuLimit == "" {
		cr.Spec.KafkaManagerCpuLimit = "2000m"
		changed = true
	}

	if cr.Spec.KafkaManagerCpuRequest == "" {
		cr.Spec.KafkaManagerCpuRequest = "500m"
		changed = true
	}

	if cr.Spec.ToolsDiskLimit == "" {
		cr.Spec.ToolsDiskLimit = "10Gi"
		changed = true
	}

	if cr.Spec.ToolsDiskRequest == "" {
		cr.Spec.ToolsDiskRequest = "1Gi"
		changed = true
	}

	if cr.Spec.ToolsMemoryLimit == "" {
		cr.Spec.ToolsMemoryLimit = "2Gi"
		changed = true
	}

	if cr.Spec.ToolsMemoryRequest == "" {
		cr.Spec.ToolsMemoryRequest = "1Gi"
		changed = true
	}

	if cr.Spec.ToolsCpuLimit == "" {
		cr.Spec.ToolsCpuLimit = "2000m"
		changed = true
	}

	if cr.Spec.ToolsCpuRequest == "" {
		cr.Spec.ToolsCpuRequest = "500m"
		changed = true
	}

	if cr.Spec.ToolsAdminDingUrl == "" {
		cr.Spec.ToolsAdminDingUrl = "https://oapi.dingtalk.com/robot/send?access_token="
		changed = true
	}

	if cr.Spec.ToolsImage == "" {
		cr.Spec.ToolsImage = "registry.cn-hangzhou.aliyuncs.com/jianzhiunique/kafka-management:latest"
		changed = true
	}

	if cr.Spec.ExporterImage == "" {
		cr.Spec.ExporterImage = "danielqsj/kafka-exporter:latest"
		changed = true
	}

	if cr.Spec.ExporterDiskLimit == "" {
		cr.Spec.ExporterDiskLimit = "10Gi"
		changed = true
	}

	if cr.Spec.ExporterDiskRequest == "" {
		cr.Spec.ExporterDiskRequest = "1Gi"
		changed = true
	}

	if cr.Spec.ExporterMemoryLimit == "" {
		cr.Spec.ExporterMemoryLimit = "2Gi"
		changed = true
	}

	if cr.Spec.ExporterMemoryRequest == "" {
		cr.Spec.ExporterMemoryRequest = "1Gi"
		changed = true
	}

	if cr.Spec.ExporterCpuLimit == "" {
		cr.Spec.ExporterCpuLimit = "2000m"
		changed = true
	}

	if cr.Spec.ExporterCpuRequest == "" {
		cr.Spec.ExporterCpuRequest = "500m"
		changed = true
	}

	if cr.Spec.IngressNamespace == "" {
		cr.Spec.IngressNamespace = "default"
		changed = true
	}

	return changed
}
