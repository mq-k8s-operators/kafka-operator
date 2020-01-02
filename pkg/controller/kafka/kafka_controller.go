package kafka

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	jianzhiuniquev1 "github.com/jianzhiunique/kafka-operator/pkg/apis/jianzhiunique/v1"
	"github.com/jianzhiunique/kafka-operator/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	v1beta12 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
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
	"time"

	_ "github.com/jianzhiunique/kafka-operator/pkg/utils"
	_ "github.com/pravega/zookeeper-operator/pkg/apis"
	"github.com/pravega/zookeeper-operator/pkg/apis/zookeeper/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
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
	// Watch for changes to secondary resource ZookeeperCluster and requeue the owner Kafka
	err = c.Watch(&source.Kind{Type: &v1beta1.ZookeeperCluster{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &jianzhiuniquev1.Kafka{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource StatefulSet and requeue the owner Kafka
	err = c.Watch(&source.Kind{Type: &appsv1.StatefulSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &jianzhiuniquev1.Kafka{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Service and requeue the owner Kafka
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &jianzhiuniquev1.Kafka{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Deployment and requeue the owner Kafka
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &jianzhiuniquev1.Kafka{},
	})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Ingress and requeue the owner Kafka
	err = c.Watch(&source.Kind{Type: &v1beta12.Ingress{}}, &handler.EnqueueRequestForOwner{
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
	log    logr.Logger
}

const ReconcileTime = 30 * time.Second

type reconcileFun func(cluster *jianzhiuniquev1.Kafka) error

// Reconcile reads that state of the cluster for a Kafka object and makes changes based on the state read
// and what is in the Kafka.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileKafka) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	r.log = log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	r.log.Info("Reconciling Kafka")

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
		return reconcile.Result{}, fmt.Errorf("GET Kafka CR fail : %s", err)
	}

	//check if default values will be used
	changed := utils.CheckCR(instance)

	if changed {
		r.log.Info("Setting default settings for kafka", instance)
		if err := r.client.Update(context.TODO(), instance); err != nil {
			return reconcile.Result{}, fmt.Errorf("Setting default fail : %s", err)
		}
		//retry reconcile
		return reconcile.Result{Requeue: true}, nil
	}

	// reconcile
	for _, fun := range []reconcileFun{
		r.reconcileZooKeeper,
		r.reconcileKafka,
		r.reconcileKafkaManager,
		r.reconcileClusterStatus,
		r.reconcileKafkaProxy,
	} {
		if err = fun(instance); err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

func (r *ReconcileKafka) reconcileZooKeeper(instance *jianzhiuniquev1.Kafka) (err error) {
	zk := utils.NewZkForCR(instance)
	// Set zk instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, zk, r.scheme); err != nil {
		return fmt.Errorf("SET ZK Owner fail : %s", err)
	}

	//检查zk是否存在
	foundzk := &v1beta1.ZookeeperCluster{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: zk.Name, Namespace: zk.Namespace}, foundzk)
	//如果sts不存在,新建
	if err != nil && errors.IsNotFound(err) {
		r.log.Info("Creating a new Zk", "Namespace", zk.Namespace, "Name", zk.Name)

		err = r.client.Create(context.TODO(), zk)
		if err != nil {
			return fmt.Errorf("Create ZK Fail : %s", err)
		}
		//创建zk成功,继续进行
	} else if err != nil {
		//有异常
		return fmt.Errorf("GET ZK Fail : %s", err)
	} else {
		//检查是否有变化，如果有变化，则Update
		//对于zk，目前只更新节点数
		if zk.Spec.Replicas != foundzk.Spec.Replicas {
			foundzk.Spec.Replicas = zk.Spec.Replicas
			err = r.client.Update(context.TODO(), foundzk)
			if err != nil {
				return fmt.Errorf("Update ZK Fail : %s", err)
			}
		}
	}

	//检查zk是否就绪
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: zk.Name, Namespace: zk.Namespace}, foundzk)
	if err != nil {
		return fmt.Errorf("CHECK ZK Status Fail : %s", err)
	}
	if foundzk.Status.ReadyReplicas != zk.Spec.Replicas {
		r.log.Info("Zk Not Ready", "Namespace", zk.Namespace, "Name", zk.Name)
		return fmt.Errorf("Zk Not Ready")
	}
	r.log.Info("Zk Ready", "Namespace", zk.Namespace, "Name", zk.Name, "found", foundzk)

	return nil
}

func (r *ReconcileKafka) reconcileKafka(instance *jianzhiuniquev1.Kafka) (err error) {
	sts := utils.NewStsForCR(instance)
	// Set Kafka instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, sts, r.scheme); err != nil {
		return fmt.Errorf("SET Kafka Owner fail : %s", err)
	}

	//检查sts是否存在
	found := &appsv1.StatefulSet{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: sts.Name, Namespace: sts.Namespace}, found)
	//如果sts不存在,新建
	if err != nil && errors.IsNotFound(err) {
		r.log.Info("Creating a new Sts", "Sts.Namespace", sts.Namespace, "Sts.Name", sts.Name)
		err = r.client.Create(context.TODO(), sts)
		if err != nil {
			return fmt.Errorf("Create sts fail : %s", err)
		}
	} else if err != nil {
		//有异常
		return fmt.Errorf("GET sts fail : %s", err)
	} else {
		utils.SyncKafkaSts(found, sts)
		err = r.client.Update(context.TODO(), found)
		if err != nil {
			return fmt.Errorf("Update ZK Fail : %s", err)
		}
	}

	//检查kfk是否可用
	//创建kafka service
	svc := utils.NewSvcForCR(instance)
	if err := controllerutil.SetControllerReference(instance, svc, r.scheme); err != nil {
		return fmt.Errorf("SET Kafka Owner fail : %s", err)
	}
	foundSvc := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: svc.Name, Namespace: svc.Namespace}, foundSvc)

	if err != nil && errors.IsNotFound(err) {
		r.log.Info("Creating a new kafka headless svc", "Svc.Namespace", svc.Namespace, "Svc.Name", svc.Name)
		err = r.client.Create(context.TODO(), svc)
		if err != nil {
			return fmt.Errorf("Create kafka headless svc fail : %s", err)
		}
	} else if err != nil {
		return fmt.Errorf("GET svc fail : %s", err)
	}

	return nil
}

func (r *ReconcileKafka) reconcileKafkaManager(instance *jianzhiuniquev1.Kafka) (err error) {
	km := utils.NewKafkaManagerForCR(instance)

	if err := controllerutil.SetControllerReference(instance, km, r.scheme); err != nil {
		return fmt.Errorf("SET Kafka Owner fail : %s", err)
	}
	foundKm := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: km.Name, Namespace: km.Namespace}, foundKm)

	if err != nil && errors.IsNotFound(err) {
		r.log.Info("Creating a new kafka manager deployment", "Namespace", km.Namespace, "Name", km.Name)
		err = r.client.Create(context.TODO(), km)
		if err != nil {
			return fmt.Errorf("Create kafka manager deployment fail : %s", err)
		}
	} else if err != nil {
		return fmt.Errorf("GET kafka manager deployment fail : %s", err)
	}

	kmsvc := utils.NewKmSvcForCR(instance)
	if err := controllerutil.SetControllerReference(instance, kmsvc, r.scheme); err != nil {
		return fmt.Errorf("SET Kafka Owner fail : %s", err)
	}
	foundKmSvc := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: kmsvc.Name, Namespace: kmsvc.Namespace}, foundKmSvc)

	if err != nil && errors.IsNotFound(err) {
		r.log.Info("Creating a new kafka manager svc", "Namespace", kmsvc.Namespace, "Name", kmsvc.Name)
		err = r.client.Create(context.TODO(), kmsvc)
		if err != nil {
			return fmt.Errorf("Create kafka manager svc fail : %s", err)
		}
	} else if err != nil {
		return fmt.Errorf("GET kafka manager svc fail : %s", err)
	}

	kmi := utils.NewKafkaManagerIngressForCR(instance)
	if err := controllerutil.SetControllerReference(instance, kmi, r.scheme); err != nil {
		return fmt.Errorf("SET Kafka Owner fail : %s", err)
	}
	foundKmi := &v1beta12.Ingress{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: kmi.Name, Namespace: kmi.Namespace}, foundKmi)

	if err != nil && errors.IsNotFound(err) {
		r.log.Info("Creating a new kafka manager ingress", "Namespace", kmi.Namespace, "Name", kmi.Name)
		err = r.client.Create(context.TODO(), kmi)
		if err != nil {
			return fmt.Errorf("Create kafka manager ingress fail : %s", err)
		}
	} else if err != nil {
		return fmt.Errorf("GET kafka manager ingress fail : %s", err)
	}

	return nil
}

func (r *ReconcileKafka) reconcileKafkaProxy(instance *jianzhiuniquev1.Kafka) (err error) {
	//check
	dep := utils.NewProxyForCR(instance)
	if err := controllerutil.SetControllerReference(instance, dep, r.scheme); err != nil {
		return fmt.Errorf("SET proxy Owner fail : %s", err)
	}
	found := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: dep.Name, Namespace: dep.Namespace}, found)

	if err != nil && errors.IsNotFound(err) {
		r.log.Info("Creating a new Proxy", "Namespace", dep.Namespace, "Name", dep.Name)
		err = r.client.Create(context.TODO(), dep)
		if err != nil {
			return fmt.Errorf("Create proxy fail : %s", err)
		}
	} else if err != nil {
		return fmt.Errorf("GET proxy fail : %s", err)
	}

	//check svc
	svc := utils.NewMqpSvcForCR(instance)
	if err := controllerutil.SetControllerReference(instance, svc, r.scheme); err != nil {
		return fmt.Errorf("SET mqp svc Owner fail : %s", err)
	}
	foundSvc := &corev1.Service{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: svc.Name, Namespace: svc.Namespace}, foundSvc)

	if err != nil && errors.IsNotFound(err) {
		r.log.Info("Creating proxy svc", "Namespace", svc.Namespace, "Name", svc.Name)
		err = r.client.Create(context.TODO(), svc)
		if err != nil {
			return fmt.Errorf("Create proxy svc fail : %s", err)
		}
	} else if err != nil {
		return fmt.Errorf("GET proxy svc fail : %s", err)
	}

	return nil
}

func (r *ReconcileKafka) reconcileClusterStatus(instance *jianzhiuniquev1.Kafka) (err error) {
	return r.client.Status().Update(context.TODO(), instance)
}
