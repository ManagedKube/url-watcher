package watcher

import (
	"context"
	"k8s.io/api/extensions/v1beta1"

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

var log = logf.Log.WithName("controller_watcher")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Watcher Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileWatcher{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("watcher-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Watcher
	err = c.Watch(&source.Kind{Type: &urlwatcherv1alpha1.Watcher{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Watcher
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &urlwatcherv1alpha1.Watcher{},
	})
	if err != nil {
		return err
	}

	// Watch ingress resources
	// Setting the handler to EnqueueRequestForObject looks like it will watch objects that
	// is not owned by this controller
	// Doc: https://godoc.org/sigs.k8s.io/controller-runtime/pkg/handler
	//err = c.Watch(&source.Kind{Type: &v1beta1.Ingress{}}, &handler.EnqueueRequestForObject{})
	//if err != nil {
	//	return err
	//}
	//
	//
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

// blank assignment to verify that ReconcileWatcher implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileWatcher{}

// ReconcileWatcher reconciles a Watcher object
type ReconcileWatcher struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Watcher object and makes changes based on the state read
// and what is in the Watcher.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileWatcher) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name, "other", request.NamespacedName.Name)
	reqLogger.Info("Reconciling Watcher")
	reqLogger.Info("doing something....")

	// Get the CRD values
	urlwatcher := &urlwatcherv1alpha1.Watcher{}
	err := r.client.Get(context.TODO(), request.NamespacedName, urlwatcher)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.Info("Memcached resource not found. Ignoring since object must be deleted.")
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.Error(err, "Failed to get Memcached.")
		return reconcile.Result{}, err
	}

	log.Info("urlwatcher", "name", urlwatcher.Name)
	log.Info("urlwatcher", "namepspace", urlwatcher.Namespace)
	log.Info("urlwatcher", "size", urlwatcher.Spec.Size)
	log.Info("urlwatcher", "foo", urlwatcher.Spec.Foo)



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

	reqLogger.Info("XXXXXXXXXXXXXXXXXXXXXX")
	for _, ingressItem := range ingressList.Items {
		//reqLogger.WithValues("ingress.name", ingressItem.Name)
		log.Info("Ingress list", "ingress.Name", ingressItem.Name)
		log.Info("Ingress list", "ingress.Annotations", ingressItem.Annotations)
		log.Info("Ingress list", "ingress.Status.String", ingressItem.Status.String())
		log.Info("Ingress list", "ingress.GetOwnerReferences", ingressItem.GetOwnerReferences())
	}

	reqLogger.Info("XXXXXXXXXXXXXXXXXXXXXX")
	//////////////////////////////////////////////////////////////////////////////


	// Fetch the Watcher instance
	instance := &urlwatcherv1alpha1.Watcher{}
	err = r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Define a new Pod object
	pod := newPodForCR(instance)

	// Set Watcher instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	found := &corev1.Pod{}
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
	}







	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	return reconcile.Result{}, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *urlwatcherv1alpha1.Watcher) *corev1.Pod {
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
