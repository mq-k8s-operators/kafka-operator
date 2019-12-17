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

	//optional: for mark user who register this cluster
	Username string `json:"username,omitempty"`
	Image    string `json:"image,omitempty"`
	// +kubebuilder:validation:Minimum=3
	Size int32 `json:"size,omitempty"`
	// resource requests and limits
	DiskLimit     string `json:"disk_limit,omitempty"`
	DiskRequest   string `json:"disk_request,omitempty"`
	MemoryRequest string `json:"memory_request,omitempty"`
	MemoryLimit   string `json:"memory_limit,omitempty"`
	CpuLimit      string `json:"cpu_limit,omitempty"`
	CpuRequest    string `json:"cpu_request,omitempty"`
	// we suggest to use local pv, so the storage class name must be set
	StorageClassName string `json:"storage_class_name"`
	// specify the hostname suffix for kafka manager,
	// for example, when ObjectMeta.Name = test and this field is .km.com,
	// we will generate a ingress whose rule host is test.km.com for kafka manager
	// then you can bind hosts test.km.com to access it
	// default value is .kfk.cloudmq.com
	KafkaManagerHost string `json:"kafka_manager_host,omitempty"`
	// +kubebuilder:validation:Enum=1,3,5,7
	ZkSize        int32  `json:"zk_size,omitempty"`
	ZkDiskLimit   string `json:"zk_disk_limit,omitempty"`
	ZkDiskRequest string `json:"zk_disk_request,omitempty"`

	//kafka jvm
	// +kubebuilder:validation:Minimum=1
	KafkaJvmXms int `json:"kafka_jvm_xms,omitempty"`
	// +kubebuilder:validation:Minimum=1
	KafkaJvmXmx int `json:"kafka_jvm_xmx,omitempty"`

	//kafka's config items
	// +kubebuilder:validation:Minimum=1
	KafkaNumPartitions            int32 `json:"default_partitions,omitempty"`
	KafkaLogRetentionHours        int32 `json:"log_hours,omitempty"`
	KafkaLogRetentionBytes        int64 `json:"log_bytes,omitempty"`
	KafkaDefaultReplicationFactor int32 `json:"replication_factor,omitempty"`
	KafkaMessageMaxBytes          int64 `json:"message_max_bytes,omitempty"`
	// +kubebuilder:validation:Enum=gzip,snappy,lz4,uncompressed,producer
	KafkaCompressionType             string `json:"compression_type,omitempty"`
	KafkaUncleanLeaderElectionEnable bool   `json:"unclean_election,omitempty"`
	// +kubebuilder:validation:Enum=delete,compact
	KafkaLogCleanupPolicy string `json:"cleanup_policy,omitempty"`
	// +kubebuilder:validation:Enum=CreateTime,LogAppendTime
	KafkaLogMessageTimestampType string `json:"message_timestamp_type,omitempty"`
}

// KafkaStatus defines the observed state of Kafka
// +k8s:openapi-gen=true
type KafkaStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	ZkUrl                string `json:zk_url`
	KafkaUrl             string `json:kafka_url`
	KafkaManagerUrl      string `json:kafka_manager_url`
	KafkaManagerPassword string `json:kafka_manager_password`
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
