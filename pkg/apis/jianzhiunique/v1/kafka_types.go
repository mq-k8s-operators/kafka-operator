package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KafkaSpec defines the desired state of Kafka
// +k8s:openapi-gen=true
type KafkaSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Size int32 `json:"size"`
	Image string `json:"image"`
	DiskLimit string `json:"disk_limit"`
	DiskRequest string `json:"disk_request"`
	StorageClassName string `json:"storage_class_name"`
	KafkaManagerHost string `json:"kafka_manager_host"`
	ZkSize int32 `json:"zk_size"`
	ZkDiskLimit string `json:"zk_disk_limit"`
	ZkDiskRequest string `json:"zk_disk_request"`
	KafkaNumPartitions int32 `json:"default_partitions"`
	KafkaLogRetentionHours int32 `json:"log_hours"`
	KafkaLogRetentionBytes int32 `json:"log_bytes"`
	KafkaDefaultReplicationFactor int32 `json:"replication_factor"`
	KafkaMessageMaxBytes int64 `json:"message_max_bytes"`
	KafkaCompressionType string `json:"compression_type"`
	KafkaUncleanLeaderElectionEnable bool `json:"unclean_election"`
	KafkaLogCleanupPolicy string `json:"cleanup_policy"`
	KafkaLogMessageTimestampType string `json:"message_timestamp_type"`
}

// KafkaStatus defines the observed state of Kafka
// +k8s:openapi-gen=true
type KafkaStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Kafka is the Schema for the kafkas API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=kafkas,scope=Namespaced
type Kafka struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KafkaSpec   `json:"spec,omitempty"`
	Status KafkaStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KafkaList contains a list of Kafka
type KafkaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Kafka `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Kafka{}, &KafkaList{})
}
