package web

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	jianzhiuniquev1 "github.com/jianzhiunique/kafka-operator/pkg/apis/jianzhiunique/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func StartWeb(manager manager.Manager) {
	http := gin.Default()

	http.GET("/ping", func(c *gin.Context) {
		err := manager.GetClient().Create(context.TODO(), &jianzhiuniquev1.Kafka{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Kafka",
				APIVersion: "jianzhiunique.github.io/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test",
				Namespace: "default",
			},
			Spec: jianzhiuniquev1.KafkaSpec{
				StorageClassName: "standard",
				KafkaManagerHost: ".km.com",
			},
		})

		if err != nil {
			fmt.Println(err)
		}

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	go http.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
