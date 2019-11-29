package kafka

import (
	"context"
	"fmt"
	v1beta12 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
	"math/rand"
	"time"

	jianzhiuniquev1 "github.com/jianzhiunique/kafka-operator/pkg/apis/jianzhiunique/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	appsv1 "k8s.io/api/apps/v1"
	"github.com/pravega/zookeeper-operator/pkg/apis/zookeeper/v1beta1"
	_ "github.com/pravega/zookeeper-operator/pkg/apis"
)

var log = logf.Log.WithName("controller_kafka")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Kafka Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileKafka{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("kafka-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Kafka
	err = c.Watch(&source.Kind{Type: &jianzhiuniquev1.Kafka{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Kafka
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &jianzhiuniquev1.Kafka{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileKafka implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileKafka{}

// ReconcileKafka reconciles a Kafka object
type ReconcileKafka struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Kafka object and makes changes based on the state read
// and what is in the Kafka.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileKafka) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Kafka")

	// Fetch the Kafka instance
	// 获取Kafka CR
	instance := &jianzhiuniquev1.Kafka{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, fmt.Errorf("GET Kafka CR fail : %s",err)
	}

	//check and handle default value
	checkCR(instance)

	// Define a new Pod object
	//pod := newPodForCR(instance)

	zk := newZkForCR(instance)

	// Set zk instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, zk, r.scheme); err != nil {
		return reconcile.Result{}, fmt.Errorf("SET ZK Owner fail : %s",err)
	}

	// Check if this Pod already exists
	/*found := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}*/

	//检查zk是否存在
	foundzk := &v1beta1.ZookeeperCluster{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: zk.Name, Namespace:zk.Namespace}, foundzk)
	//如果sts不存在,新建
	if err != nil && errors.IsNotFound(err){
		reqLogger.Info("Creating a new Zk", "Namespace", zk.Namespace, "Name", zk.Name)

		err = r.client.Create(context.TODO(), zk)
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("Create ZK Fail : %s",err)
		}
		//创建zk成功,继续进行
	}else if err != nil {
		//有异常
		return  reconcile.Result{}, fmt.Errorf("GET ZK Fail : %s",err)
	}

	//等待
	//time.Sleep(time.Second * 1)

	//检查zk是否就绪
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: zk.Name, Namespace:zk.Namespace}, foundzk)
	if err != nil {
		return  reconcile.Result{}, fmt.Errorf("CHECK ZK Status Fail : %s",err)
	}
	if foundzk.Status.ReadyReplicas != zk.Spec.Replicas {
		reqLogger.Info("Zk Not Ready", "Namespace", zk.Namespace, "Name", zk.Name)
		return reconcile.Result{Requeue: true, RequeueAfter: time.Second*5}, nil
	}
	reqLogger.Info("Zk Ready", "Namespace", zk.Namespace, "Name", zk.Name, "found", foundzk)

	//获取zk服务地址

	zkUrl := zk.Name + "-client:2181"
	sts := newStsForCR(instance, zkUrl)
	// Set Kafka instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, sts, r.scheme); err != nil {
		return reconcile.Result{}, fmt.Errorf("SET Kafka Owner fail : %s",err)
	}

	//检查sts是否存在
	found := &appsv1.StatefulSet{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: sts.Name, Namespace:sts.Namespace}, found)
	//如果sts不存在,新建
	if err != nil && errors.IsNotFound(err){
		reqLogger.Info("Creating a new Sts", "Sts.Namespace", sts.Namespace, "Sts.Name", sts.Name)
		err = r.client.Create(context.TODO(), sts)
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("Create sts fail : %s",err)
		}

		//创建成功

	}else if err != nil {
		//有异常
		return  reconcile.Result{}, fmt.Errorf("GET sts fail : %s",err)
	}

	//检查kfk是否可用

	//创建kafka service
	svc := newSvcForCR(instance)
	if err := controllerutil.SetControllerReference(instance, svc, r.scheme); err != nil {
		return reconcile.Result{}, fmt.Errorf("SET Kafka Owner fail : %s",err)
	}
	foundSvc := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name:svc.Name, Namespace:svc.Namespace}, foundSvc)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new kafka headless svc", "Svc.Namespace", svc.Namespace, "Svc.Name", svc.Name)
		err = r.client.Create(context.TODO(), svc)
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("Create kafka headless svc fail : %s",err)
		}
	}else if err != nil {
		return  reconcile.Result{}, fmt.Errorf("GET svc fail : %s",err)
	}

	//创建kafka nodeport,对外提供服务

	//部署kmonitor等周边资源
	km := newKafkaManagerForCR(instance, zkUrl)
	if err := controllerutil.SetControllerReference(instance, km, r.scheme); err != nil {
		return reconcile.Result{}, fmt.Errorf("SET Kafka Owner fail : %s",err)
	}
	foundKm := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name:km.Name, Namespace:km.Namespace}, foundKm)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new kafka manager deployment", "Namespace", km.Namespace, "Name", km.Name)
		err = r.client.Create(context.TODO(), km)
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("Create kafka manager deployment fail : %s",err)
		}
	}else if err != nil {
		return  reconcile.Result{}, fmt.Errorf("GET kafka manager deployment fail : %s",err)
	}

	kmsvc := newKmSvcForCR(instance)
	if err := controllerutil.SetControllerReference(instance, kmsvc, r.scheme); err != nil {
		return reconcile.Result{}, fmt.Errorf("SET Kafka Owner fail : %s",err)
	}
	foundKmSvc := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name:kmsvc.Name, Namespace:kmsvc.Namespace}, foundKmSvc)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new kafka manager svc", "Namespace", kmsvc.Namespace, "Name", kmsvc.Name)
		err = r.client.Create(context.TODO(), kmsvc)
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("Create kafka manager svc fail : %s",err)
		}
	}else if err != nil {
		return  reconcile.Result{}, fmt.Errorf("GET kafka manager svc fail : %s",err)
	}

	kmi := newKafkaManagerIngressForCR(instance)
	if err := controllerutil.SetControllerReference(instance, kmi, r.scheme); err != nil {
		return reconcile.Result{}, fmt.Errorf("SET Kafka Owner fail : %s",err)
	}
	foundKmi := &v1beta12.Ingress{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name:kmi.Name, Namespace:kmi.Namespace}, foundKmi)

	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new kafka manager ingress", "Namespace", kmi.Namespace, "Name", kmi.Name)
		err = r.client.Create(context.TODO(), kmi)
		if err != nil {
			return reconcile.Result{}, fmt.Errorf("Create kafka manager ingress fail : %s",err)
		}
	}else if err != nil {
		return  reconcile.Result{}, fmt.Errorf("GET kafka manager ingress fail : %s",err)
	}



	// Pod already exists - don't requeue
	// sts已经存在，对比
	//reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	return reconcile.Result{}, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *jianzhiuniquev1.Kafka) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}

func newStsForCR(cr *jianzhiuniquev1.Kafka, zk string) *appsv1.StatefulSet{
	accessModes := make([]corev1.PersistentVolumeAccessMode, 0)
	accessModes = append(accessModes, corev1.ReadWriteOnce)
	pvc := make([]corev1.PersistentVolumeClaim, 0)
	pvc = append(pvc, corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name: "kfk-data",
		},
		Spec:       corev1.PersistentVolumeClaimSpec{
			StorageClassName: &cr.Spec.StorageClassName,
			Resources: corev1.ResourceRequirements{
				Limits:   corev1.ResourceList{
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
			Name: "KAFKA_ZOOKEEPER_CONNECT",
			Value: zk,
		},
		corev1.EnvVar{
			Name: "BROKER_ID_COMMAND",
			Value: "hostname | awk -F'-' '{print $$4}'",
		},
		corev1.EnvVar{
			Name: "KAFKA_AUTO_CREATE_TOPICS_ENABLE",
			Value: "false",
		},
		corev1.EnvVar{
			Name: "KAFKA_DELETE_TOPIC_ENABLE",
			Value: "true",
		},
		corev1.EnvVar{
			Name: "KAFKA_LISTENERS",
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
			Name: "KAFKA_ADVERTISED_LISTENERS",
			Value: "PLAINTEXT://$(KAFKA_ADVERTISED_HOST_NAME).kafka.$(MY_POD_NAMESPACE).svc.cluster.local:9092",
		},
		corev1.EnvVar{
			Name: "KAFKA_LOG_DIRS",
			Value: "/data/kafka",
		},
		corev1.EnvVar{
			Name: "KAFKA_ZOOKEEPER_SESSION_TIMEOUT_MS",
			Value: "120000",
		},
		corev1.EnvVar{
			Name: "KAFKA_AUTO_LEADER_REBANLANCE_ENABLE",
			Value: "false",
		},corev1.EnvVar{
			Name: "KAFKA_OFFSETS_RETENTION_MINUTES",
			Value: "14400",
		},
		corev1.EnvVar{
			Name: "JMX_PORT",
			Value: "9999",
		},
		corev1.EnvVar{
			Name: "KAFKA_NUM_PARTITIONS",
			Value: "3",
		},
		corev1.EnvVar{
			Name: "KAFKA_LOG_RETENTION_HOURS",
			Value: "168",
		},
		corev1.EnvVar{
			Name: "KAFKA_LOG_RETENTION_BYTES",
			Value: "-1",
		},
		corev1.EnvVar{
			Name: "KAFKA_DEFAULT_REPLICATION_FACTOR",
			Value: "2",
		},
		corev1.EnvVar{
			Name: "KAFKA_MESSAGE_MAX_BYTES",
			Value: "1073741824",
		},
		corev1.EnvVar{
			Name: "KAFKA_COMPRESSION_TYPE",
			Value: "producer",
		},
		corev1.EnvVar{
			Name: "KAFKA_UNCLEAN_LEADER_ELECTION_ENABLE",
			Value: "false",
		},
		corev1.EnvVar{
			Name: "KAFKA_LOG_CLEANUP_POLICY",
			Value: "delete",
		},
		corev1.EnvVar{
			Name: "KAFKA_LOG_MESSAGE_TIMESTAMP_TYPE",
			Value: "CreateTime",
		},
		corev1.EnvVar{
			Name: "KAFKA_NUM_NETWORK_THREADS",
			Value: "3",
		},
		corev1.EnvVar{
			Name: "KAFKA_NUM_IO_THREADS",
			Value: "8",
		},
		corev1.EnvVar{
			Name: "KAFKA_NUM_RECOVERY_THREADS_PER_DATA_DIR",
			Value: "2",
		},
		corev1.EnvVar{
			Name: "KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR",
			Value: "3",
		},
		corev1.EnvVar{
			Name: "KAFKA_NUM_REPLICA_FETCHERS",
			Value: "2",
		},
		corev1.EnvVar{
			Name: "KAFKA_MIN_INSYNC_REPLICAS",
			Value: "2",
		},
		corev1.EnvVar{
			Name: "KAFKA_GROUP_INITIAL_REBANLANCE_DELAY_MS",
			Value: "3000",
		},
		corev1.EnvVar{
			Name: "KAFKA_TRANSACTION_STATE_LOG_REPLICATION_FACTOR",
			Value: "3",
		},
		corev1.EnvVar{
			Name: "KAFKA_HEAP_OPTS",
			Value: "-Xms1g -Xmx1g -XX:MetaspaceSize=96m -XX:+UseG1GC -XX:MaxGCPauseMillis=20 -XX:InitiatingHeapOccupancyPercent=35 -XX:G1HeapRegionSize=16M -XX:MinMetaspaceFreeRatio=50 -XX:MaxMetaspaceFreeRatio=80 -XX:ParallelGCThreads=16",
		},
	)
	vms := make([]corev1.VolumeMount, 0)
	vms = append(vms, corev1.VolumeMount{
		Name:             "kfk-data",
		MountPath:        "/data/kafka",
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
		Name:                     "kafka",
		Image:                    cr.Spec.Image,
		Ports:                    ports,
		Env:                      envs,
		VolumeMounts:             vms,
		LivenessProbe:            &healthCheck,
		ReadinessProbe:           &healthCheck,
	}
	containers = append(containers, kfk)

	return &appsv1.StatefulSet{
		TypeMeta:   metav1.TypeMeta{
			Kind: "StatefulSet",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "kfk-sts-" + cr.Name,
			Namespace: cr.Namespace,
		},
		Spec:       appsv1.StatefulSetSpec{
			Replicas: &cr.Spec.Size,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "kfk-pod-" + cr.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "kfk-pod-" + cr.Name,
					},
				},
				Spec:       corev1.PodSpec{
					Containers: containers,
				},
			},
			VolumeClaimTemplates: pvc,
		},
		Status:     appsv1.StatefulSetStatus{},
	}
}

func newZkForCR(cr *jianzhiuniquev1.Kafka) *v1beta1.ZookeeperCluster{
	return &v1beta1.ZookeeperCluster{
		TypeMeta:   metav1.TypeMeta{
			Kind:       "ZookeeperCluster",
			APIVersion: "zookeeper.pravega.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "kfk-zk-" + cr.Name,
			Namespace: cr.Namespace,
		},
		Spec:       v1beta1.ZookeeperClusterSpec{
			//Image: v1beta1.ContainerImage{
			//	Tag: "0.2.4",
			//},
			Replicas: cr.Spec.ZkSize,
			Persistence: &v1beta1.Persistence{
				VolumeReclaimPolicy:       v1beta1.VolumeReclaimPolicyDelete,
				PersistentVolumeClaimSpec: corev1.PersistentVolumeClaimSpec{
					Resources:        corev1.ResourceRequirements{
						Limits:   corev1.ResourceList{
							corev1.ResourceStorage: resource.MustParse(cr.Spec.ZkDiskLimit),
						},
						Requests: corev1.ResourceList{
							corev1.ResourceStorage: resource.MustParse(cr.Spec.ZkDiskRequest),
						},
					},
					StorageClassName: &cr.Spec.StorageClassName,
				},
			},
		},
	}
}

func newKmSvcForCR(cr *jianzhiuniquev1.Kafka) *corev1.Service{
	port := corev1.ServicePort{Port:9000}
	ports := make([]corev1.ServicePort, 0)
	ports = append(ports, port)
	return &corev1.Service{
		TypeMeta:   metav1.TypeMeta{
			APIVersion: "v1",
			Kind: "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "kfk-m-svc-" + cr.Name,
			Namespace: cr.Namespace,
		},
		Spec:       corev1.ServiceSpec{
			Ports: ports,
			Selector: map[string]string{
				"app": "kfk-m-" + cr.Name,
			},
		},
	}
}

func newSvcForCR(cr *jianzhiuniquev1.Kafka) *corev1.Service{
	port := corev1.ServicePort{Port:9092}
	ports := make([]corev1.ServicePort, 0)
	ports = append(ports, port)
	return &corev1.Service{
		TypeMeta:   metav1.TypeMeta{
			APIVersion: "v1",
			Kind: "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "kfk-svc-" + cr.Name,
			Namespace: cr.Namespace,
		},
		Spec:       corev1.ServiceSpec{
			Ports: ports,
			Selector: map[string]string{
				"app": "kfk-pod-" + cr.Name,
			},
		},
	}
}

func newKafkaManagerForCR(cr *jianzhiuniquev1.Kafka, zk string) *appsv1.Deployment{
	cport := corev1.ContainerPort{ContainerPort:9000}
	cports := make([]corev1.ContainerPort, 0)
	cports = append(cports, cport)

	envs := make([]corev1.EnvVar, 0)
	envs = append(envs,
		corev1.EnvVar{
			Name: "ZK_HOSTS",
			Value: zk,
		},
		corev1.EnvVar{
			Name: "KAFKA_MANAGER_AUTH_ENABLED",
			Value: "true",
		},
		corev1.EnvVar{
			Name: "KAFKA_MANAGER_USERNAME",
			Value: "admin",
		},
		corev1.EnvVar{
			Name: "KAFKA_MANAGER_PASSWORD",
			Value: GetRandomString(16),
		},
	)

	containers := make([]corev1.Container, 0)
	container := corev1.Container{
		Name:                     "kfk-m-c-" + cr.Name,
		Image:                    "kafkamanager/kafka-manager:latest",
		Ports:                    cports,
		Env:                      envs,
	}
	containers = append(containers, container)

	port := corev1.ServicePort{Port:9092}
	ports := make([]corev1.ServicePort, 0)
	ports = append(ports, port)

	var replica int32
	replica = 1
	return &appsv1.Deployment{
		TypeMeta:   metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind: "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "kfk-manager-" + cr.Name,
			Namespace: cr.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas:                &replica,
			Selector:                &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "kfk-m-" + cr.Name,
				},
			},
			Template:                corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "kfk-m-" + cr.Name,
					},
				},
				Spec:       corev1.PodSpec{
					Containers:                    containers,
				},
			},
		},
	}
}

func  GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func newKafkaManagerIngressForCR(cr *jianzhiuniquev1.Kafka) *v1beta12.Ingress{
	paths := make([]v1beta12.HTTPIngressPath, 0)
	path := v1beta12.HTTPIngressPath{
		Path:    "/",
		Backend: v1beta12.IngressBackend{
			ServicePort: intstr.IntOrString{
				IntVal: 9000,
			},
			ServiceName: "kfk-m-svc-" + cr.Name,
		},
	}
	paths = append(paths, path)

	rules := make([]v1beta12.IngressRule, 0)
	rule := v1beta12.IngressRule{
		Host: cr.Name + cr.Spec.KafkaManagerHost,
		IngressRuleValue: v1beta12.IngressRuleValue{
			HTTP: &v1beta12.HTTPIngressRuleValue{
				Paths: paths,
			},
		},
	}
	rules = append(rules, rule)

	port := corev1.ServicePort{Port:9092}
	ports := make([]corev1.ServicePort, 0)
	ports = append(ports, port)
	return &v1beta12.Ingress{
		TypeMeta:   metav1.TypeMeta{
			APIVersion: "networking.k8s.io/v1beta1",
			Kind: "Ingress",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "km-ingress-" + cr.Name,
			Namespace: cr.Namespace,
		},
		Spec: v1beta12.IngressSpec{
			Rules:   rules,
		},
	}
}

func checkCR(cr *jianzhiuniquev1.Kafka){
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

	if cr.Spec.Size == 0 {
		cr.Spec.Size = 3
	}

	if cr.Spec.Image == "" {
		cr.Spec.Image = "wurstmeister/kafka:2.11-0.11.0.3"
	}

	if cr.Spec.DiskLimit == "" {
		cr.Spec.DiskLimit = "500Gi"
	}

	if cr.Spec.DiskRequest == "" {
		cr.Spec.DiskRequest = "100Gi"
	}

	if cr.Spec.KafkaManagerHost == "" {
		cr.Spec.KafkaManagerHost = ".km.com"
	}

	if cr.Spec.ZkSize == 0 {
		cr.Spec.ZkSize = 3
	}

	if cr.Spec.ZkDiskLimit == "" {
		cr.Spec.ZkDiskLimit = "100Gi"
	}

	if cr.Spec.ZkDiskRequest == "" {
		cr.Spec.ZkDiskRequest = "20Gi"
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
	}

	if cr.Spec.KafkaLogRetentionHours == 0 {
		cr.Spec.KafkaLogRetentionHours = 168
	}

	if cr.Spec.KafkaLogRetentionBytes == 0 {
		cr.Spec.KafkaLogRetentionHours = -1
	}

	if cr.Spec.KafkaDefaultReplicationFactor == 0 {
		cr.Spec.KafkaDefaultReplicationFactor = 2
	}

	if cr.Spec.KafkaMessageMaxBytes == 0 {
		cr.Spec.KafkaMessageMaxBytes = 1073741824
	}

	if cr.Spec.KafkaCompressionType == "" {
		cr.Spec.KafkaCompressionType = "producer"
	}

	if cr.Spec.KafkaLogCleanupPolicy == "" {
		cr.Spec.KafkaLogCleanupPolicy = "delete"
	}

	if cr.Spec.KafkaLogMessageTimestampType == "" {
		cr.Spec.KafkaLogMessageTimestampType = "CreateTime"
	}
}
