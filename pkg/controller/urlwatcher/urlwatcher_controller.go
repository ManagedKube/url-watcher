package urlwatcher

import (
	"context"
	"encoding/json"
	"fmt"
	"k8s.io/api/extensions/v1beta1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	urlwatcherv1alpha1 "managedkube.com/url-watcher/pkg/apis/urlwatcher/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_urlwatcher")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new UrlWatcher Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileUrlWatcher{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("urlwatcher-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource UrlWatcher
	err = c.Watch(&source.Kind{Type: &urlwatcherv1alpha1.UrlWatcher{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner UrlWatcher
	// Doc: https://github.com/operator-framework/operator-sdk-samples/blob/master/memcached-operator/pkg/controller/memcached/memcached_controller.go#L55
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &urlwatcherv1alpha1.UrlWatcher{},
	})
	if err != nil {
		return err
	}

	// Watch ingress resources
	// Setting the handler to EnqueueRequestForObject looks like it will watch objects that
	// is not owned by this controller
	// controller is a controller.controller
	err = c.Watch(
		&source.Kind{Type: &v1beta1.Ingress{}},
		&handler.EnqueueRequestsFromMapFunc{
			ToRequests: handler.ToRequestsFunc(func(a handler.MapObject) []reconcile.Request {
				return []reconcile.Request{
					{NamespacedName: types.NamespacedName{
						Name:      "ingress/" + a.Meta.GetName(),
						Namespace: a.Meta.GetNamespace(),
					}},
				}
			}),
		})
	if err != nil {
		// handle it
	}

	return nil
}

// blank assignment to verify that ReconcileUrlWatcher implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileUrlWatcher{}

// ReconcileUrlWatcher reconciles a UrlWatcher object
type ReconcileUrlWatcher struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a UrlWatcher object and makes changes based on the state read
// and what is in the UrlWatcher.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileUrlWatcher) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling UrlWatcher")

	// Get the CRD values
	urlwatcher := &urlwatcherv1alpha1.UrlWatcher{}
	err := r.client.Get(context.TODO(), request.NamespacedName, urlwatcher)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.Info("urlwatcherv1alpha1.UrlWatcher resource not found. Ignoring.")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.Error(err, "Failed to get urlwatcherv1alpha1.UrlWatcher.")
		return reconcile.Result{}, err
	}

	log.Info("urlwatcher", "name", urlwatcher.Name)
	log.Info("urlwatcher", "namepspace", urlwatcher.Namespace)
	log.Info("urlwatcher", "size", urlwatcher.Spec.Size)
	log.Info("urlwatcher", "foo", urlwatcher.Spec.Foo)
	log.Info("urlwatcher", "AllIngresses", urlwatcher.Spec.AllIngresses)


	// Check if the Deployment already exists, if not create a new one
	deployment := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: urlwatcher.Name, Namespace: urlwatcher.Namespace}, deployment)
	if err != nil && errors.IsNotFound(err) {
		// Define a new Deployment
		dep := r.deploymentForUrlWatcher(urlwatcher)
		reqLogger.Info("Creating a new Deployment.", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.client.Create(context.TODO(), dep)
		if err != nil {
			reqLogger.Error(err, "Failed to create new Deployment.", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return reconcile.Result{}, err
		}
		// Deployment created successfully - return and requeue
		// NOTE: that the requeue is made with the purpose to provide the deployment object for the next step to ensure the deployment size is the same as the spec.
		// Also, you could GET the deployment object again instead of requeue if you wish. See more over it here: https://godoc.org/sigs.k8s.io/controller-runtime/pkg/reconcile#Reconciler
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Deployment.")
		return reconcile.Result{}, err
	}


	// Ensure the deployment size is the same as the spec
	size := urlwatcher.Spec.Size
	if *deployment.Spec.Replicas != size {
		deployment.Spec.Replicas = &size
		err = r.client.Update(context.TODO(), deployment)
		if err != nil {
			reqLogger.Error(err, "Failed to update Deployment.", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
			return reconcile.Result{}, err
		}
	}

	enndpointSpec := []urlWatchEndpointSpec{
		{},
	}
	endpoints := urlWatchEndpoints{
		Endpoints: enndpointSpec,
	}
	urlWatchSpecParsed := urlWatchSpec{
		Watch: endpoints,
	}

	// Log out current envars
	for _, envs := range deployment.Spec.Template.Spec.Containers[0].Env {
		log.Info("Deployment", "deployment.Spec.Template.Spec.Containers[].Env", envs.Name+" | "+envs.Value)

		if(envs.Name == "ENDPOINT_TEST_JSON"){


//			fmt.Println("xxxxx")
//			var jsonBlob = []byte(`[
//	{"Name": "Platypus", "Order": "Monotremata"},
//	{"Name": "Quoll",    "Order": "Dasyuromorphia"}
//]`)
//			type Animal struct {
//				Name  string
//				Order string
//			}
//			var animals []Animal
//			err := json.Unmarshal(jsonBlob, &animals)
//			if err != nil {
//				fmt.Println("error:", err)
//			}
//			fmt.Println("%+v", animals)

			fmt.Println("xxxxx")
			//var jsonBlob = []byte(`{"watch":{"endpoints":[{"interval":60,"host":"www.example.com"}]}}`)
			////type Animal struct {
			////	Name  string
			////	Order string
			////}
			////var animals []Animal
			//err := json.Unmarshal(jsonBlob, &urlWatchSpecParsed)
			//if err != nil {
			//	fmt.Println("error:", err)
			//}
			//fmt.Println("%+v", urlWatchSpecParsed)
			fmt.Println("xxxxx")




			//var jsonBlob = []byte(`{"watch":{"endpoints":[]}}`)

			log.Info("Deployment", "deployment.Spec.Template.Spec.Containers[].Env", envs.Name+"|"+string(envs.Value))


			err = json.Unmarshal([]byte(envs.Value), &urlWatchSpecParsed)
			if err != nil {
				reqLogger.Error(err, "Failed to unmarchall json: ENDPOINT_TEST_JSON", "envs.Name", envs.Name, "envs.Value", envs.Value)
			}
		}
	}



	updatedEndpointSpecs := false

	//////////////////////////////////////////////////////////////////////////////
	// Listing the ingresses
	reqLogger.Info("XXXXXXXXXXXXXXXXXXXXXX")
	ingressList := &v1beta1.IngressList{}
	listOps := &client.ListOptions{}

	err = r.client.List(context.TODO(), listOps, ingressList)
	if err != nil {
		reqLogger.Error(err, "Failed to list ingress.", "Memcached.Namespace", request.Namespace, "Memcached.Name", request.Name)
		return reconcile.Result{}, err
	}

	for _, ingressItem := range ingressList.Items {
		reqLogger.Info("XXXXXXXXXXXXXXXXXXXXXX")
		//reqLogger.WithValues("ingress.name", ingressItem.Name)
		log.Info("Ingress list", "ingress.Name", ingressItem.Name)
		log.Info("Ingress list", "ingress.Annotations", ingressItem.Annotations)
		log.Info("Ingress list", "ingress.Status.String", ingressItem.Status.String())
		log.Info("Ingress list", "ingress.GetOwnerReferences", ingressItem.GetOwnerReferences())


		for _, rules := range ingressItem.Spec.Rules {
			log.Info("Ingress list.rules", "ingress.rules.hosts", rules.Host)

			if(!isHostInEndpointsList(urlWatchSpecParsed, rules.Host)){
				// Add to the List and update the deployment

				tempEndpointSpec := urlWatchEndpointSpec{
					Interval: 60,
					Protocol: "http",
					Host: rules.Host,
					Method: "GET",
					Path: "/",
					Payload: "",
					ScrapeTimeout: 30,
				}

				urlWatchSpecParsed.Watch.Endpoints = append(urlWatchSpecParsed.Watch.Endpoints, tempEndpointSpec)

				updatedEndpointSpecs = true
			}

			for _, paths := range rules.IngressRuleValue.HTTP.Paths {
				log.Info("Ingress list.rules.paths", "ingress.rules.hosts.path", paths.Path)
			}
		}

		reqLogger.Info("XXXXXXXXXXXXXXXXXXXXXX")
	}

	// Update the deployment endpoint specs
	if(updatedEndpointSpecs){
		log.Info("updatedEndpointSpecs", "updatedEndpointSpecs", "true")

		b, err := json.Marshal(urlWatchSpecParsed)

		updatedEnv := []corev1.EnvVar{
			{
				Name: "ENDPOINT_TEST_JSON",
				Value: string(b),
			},
		}

		deployment.Spec.Template.Spec.Containers[0].Env = updatedEnv
		err = r.client.Update(context.TODO(), deployment)
		if err != nil {
			reqLogger.Error(err, "Failed to update Deployment.", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
			return reconcile.Result{}, err
		}
	}


	//////////////////////////////////////////////////////////////////////////////


	//// Fetch the UrlWatcher instance
	//instance := &urlwatcherv1alpha1.UrlWatcher{}
	//err = r.client.Get(context.TODO(), request.NamespacedName, instance)
	//if err != nil {
	//	if errors.IsNotFound(err) {
	//		// Request object not found, could have been deleted after reconcile request.
	//		// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
	//		// Return and don't requeue
	//		return reconcile.Result{}, nil
	//	}
	//	// Error reading the object - requeue the request.
	//	return reconcile.Result{}, err
	//}
	//
	//// Define a new Pod object
	//pod := newPodForCR(instance)
	//
	//// Set UrlWatcher instance as the owner and controller
	//if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
	//	return reconcile.Result{}, err
	//}
	//
	//// Check if this Pod already exists
	//found := &corev1.Pod{}
	//err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	//if err != nil && errors.IsNotFound(err) {
	//	reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
	//	err = r.client.Create(context.TODO(), pod)
	//	if err != nil {
	//		return reconcile.Result{}, err
	//	}
	//
	//	// Pod created successfully - don't requeue
	//	return reconcile.Result{}, nil
	//} else if err != nil {
	//	return reconcile.Result{}, err
	//}

	// Pod already exists - don't requeue
	//reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)







	return reconcile.Result{}, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
//func newPodForCR(cr *urlwatcherv1alpha1.UrlWatcher) *corev1.Pod {
//	labels := map[string]string{
//		"app": cr.Name,
//	}
//	return &corev1.Pod{
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      cr.Name + "-pod",
//			Namespace: cr.Namespace,
//			Labels:    labels,
//		},
//		Spec: corev1.PodSpec{
//			Containers: []corev1.Container{
//				{
//					Name:    "busybox",
//					Image:   "busybox",
//					Command: []string{"sleep", "3600"},
//				},
//			},
//		},
//	}
//}


// deploymentForUrlWatcher returns a memcached Deployment object
// Doc: https://github.com/operator-framework/operator-sdk-samples/blob/master/memcached-operator/pkg/controller/memcached/memcached_controller.go#L191
func (r *ReconcileUrlWatcher) deploymentForUrlWatcher(m *urlwatcherv1alpha1.UrlWatcher) *appsv1.Deployment {
	ls := labelsForDeployment(m.Name)
	replicas := m.Spec.Size

	enndpointSpec := []urlWatchEndpointSpec{{}}
	endpoints := urlWatchEndpoints{
		Endpoints: enndpointSpec,
	}
	urlWatchSpecParsed := urlWatchSpec{
		Watch: endpoints,
	}

	b, _ := json.Marshal(urlWatchSpecParsed)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:   "busybox",
						Name:    "url-watcher",
						Command: []string{"sleep", "3600"},
						// https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.15/#envvar-v1-core
						Env: []corev1.EnvVar{
							{
								Name: "ENDPOINT_TEST_JSON",
								Value: string(b),
							},
						},
						Ports: []corev1.ContainerPort{{
							ContainerPort: 9093,
							Name:          "prometheus",
						}},
					}},
				},
			},
		},
	}
	// Set UrlWatcher instance as the owner of the Deployment.
	controllerutil.SetControllerReference(m, dep, r.scheme)
	return dep
}

// labelsForDeployment returns the labels for selecting the resources
// belonging to the given UrlWarcher CR name.
func labelsForDeployment(name string) map[string]string {
	return map[string]string{"app": "memcached", "memcached_cr": name}
}

type urlWatchSpec struct{
	Watch urlWatchEndpoints `json:"watch"`
}

type urlWatchEndpoints struct{
	Endpoints []urlWatchEndpointSpec `json:"endpoints"`
}

type urlWatchEndpointSpec struct{
	Interval int64 `json:"interval,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	Host string `json:"host,omitempty"`
	Method string `json:"method,omitempty"`
	Path string `json:"path,omitempty"`
	Payload string `json:"payload,omitempty"`
	ScrapeTimeout int64 `json:"scrapeTimeout,omitempty"`
}

// Check if the host is in the urlWatchSpecParsed list
func isHostInEndpointsList(urlWatchSpecParsed urlWatchSpec, host string) bool{

	inList := false

	for _, endpointSpec := range urlWatchSpecParsed.Watch.Endpoints{
		if(endpointSpec.Host == host){
			inList = true
		}
	}

	return inList
}