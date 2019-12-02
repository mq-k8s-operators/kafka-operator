package web

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	jianzhiuniquev1 "github.com/jianzhiunique/kafka-operator/pkg/apis/jianzhiunique/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type KafkaStruct struct {
	Name                             string `json:"name,omitempty" form:"name" binding:"required"`
	Namespace                        string `json:"namespace,omitempty" form:"namespace" binding:"required"`
	Image                            string `json:"image,omitempty" form:"image"`
	Size                             int32  `json:"size,omitempty" form:"size"`
	DiskLimit                        string `json:"disk_limit,omitempty" form:"disk_limit"`
	DiskRequest                      string `json:"disk_request,omitempty" form:"disk_request"`
	StorageClassName                 string `json:"storage_class_name" form:"storage_class_name"`
	KafkaManagerHost                 string `json:"kafka_manager_host,omitempty" form:"kafka_manager_host"`
	ZkSize                           int32  `json:"zk_size,omitempty" form:"zk_size"`
	ZkDiskLimit                      string `json:"zk_disk_limit,omitempty" form:"zk_disk_limit"`
	ZkDiskRequest                    string `json:"zk_disk_request,omitempty" form:"zk_disk_request"`
	KafkaNumPartitions               int32  `json:"default_partitions,omitempty" form:"default_partitions"`
	KafkaLogRetentionHours           int32  `json:"log_hours,omitempty" form:"log_hours"`
	KafkaLogRetentionBytes           int64  `json:"log_bytes,omitempty" form:"log_bytes"`
	KafkaDefaultReplicationFactor    int32  `json:"replication_factor,omitempty" form:"replication_factor"`
	KafkaMessageMaxBytes             int64  `json:"message_max_bytes,omitempty" form:"message_max_bytes"`
	KafkaCompressionType             string `json:"compression_type,omitempty" form:"compression_type"`
	KafkaUncleanLeaderElectionEnable bool   `json:"unclean_election,omitempty" form:"unclean_election"`
	KafkaLogCleanupPolicy            string `json:"cleanup_policy,omitempty" form:"cleanup_policy"`
	KafkaLogMessageTimestampType     string `json:"message_timestamp_type,omitempty" form:"message_timestamp_type"`
}

func StartWeb(manager manager.Manager) {
	http := gin.Default()

	http.POST("/new_kafka", func(c *gin.Context) {

		var data KafkaStruct
		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(400, gin.H{
				"status": 0,
				"data":   "",
			})

			return
		}
		err := manager.GetClient().Create(context.TODO(), &jianzhiuniquev1.Kafka{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Kafka",
				APIVersion: "jianzhiunique.github.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      data.Name,
				Namespace: data.Namespace,
			},
			Spec: jianzhiuniquev1.KafkaSpec{
				Image:                            data.Image,
				Size:                             data.Size,
				DiskLimit:                        data.DiskLimit,
				DiskRequest:                      data.DiskRequest,
				StorageClassName:                 data.StorageClassName,
				KafkaManagerHost:                 data.KafkaManagerHost,
				ZkSize:                           data.ZkSize,
				ZkDiskLimit:                      data.ZkDiskLimit,
				ZkDiskRequest:                    data.ZkDiskRequest,
				KafkaNumPartitions:               data.KafkaNumPartitions,
				KafkaLogRetentionHours:           data.KafkaLogRetentionHours,
				KafkaLogRetentionBytes:           data.KafkaLogRetentionBytes,
				KafkaDefaultReplicationFactor:    data.KafkaDefaultReplicationFactor,
				KafkaMessageMaxBytes:             data.KafkaMessageMaxBytes,
				KafkaCompressionType:             data.KafkaCompressionType,
				KafkaUncleanLeaderElectionEnable: data.KafkaUncleanLeaderElectionEnable,
				KafkaLogCleanupPolicy:            data.KafkaLogCleanupPolicy,
				KafkaLogMessageTimestampType:     data.KafkaLogMessageTimestampType,
			},
		})

		if err != nil {
			fmt.Println(err)
		}

		c.JSON(200, gin.H{
			"status": 1,
			"data":   "ok",
		})
	})

	http.GET("/list_kafka", func(c *gin.Context) {
		list := jianzhiuniquev1.KafkaList{}
		err := manager.GetClient().List(context.TODO(), &list)

		if err != nil {
			fmt.Print(err)
		}

		c.JSON(200, gin.H{
			"status": 1,
			"data":   list,
		})
	})

	go http.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
