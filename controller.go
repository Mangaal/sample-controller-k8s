package main

import (
	"context"
	"fmt"
	appsv1alpha1 "sample-controller-k8s/pkg/apis/appreplica/v1alpha1"
	clinetset "sample-controller-k8s/pkg/client/clientset/versioned"
	appReplicaInformer "sample-controller-k8s/pkg/client/informers/externalversions/appreplica/v1alpha1"
	listers "sample-controller-k8s/pkg/client/listers/appreplica/v1alpha1"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	appsInformers "k8s.io/client-go/informers/apps/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	appslisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
)

type Controller struct {
	kubeclinetset kubernetes.Interface
	appclinetset  clinetset.Interface

	deploymentLister appslisters.DeploymentLister
	deploymentSynced cache.InformerSynced

	appLister listers.AppReplicaLister
	appSynced cache.InformerSynced

	workqueue workqueue.RateLimitingInterface
}

const (
	MessageResourceExists = "Resource %q already exists and is not managed by AppReplica"
)

func NewController(
	kubeclinetset kubernetes.Interface,
	appclinetset clinetset.Interface,
	appReplicaInformer appReplicaInformer.AppReplicaInformer,
	deploymentInformer appsInformers.DeploymentInformer) *Controller {

	controller := &Controller{
		kubeclinetset: kubeclinetset,
		appclinetset:  appclinetset,

		deploymentLister: deploymentInformer.Lister(),
		deploymentSynced: deploymentInformer.Informer().HasSynced,

		appLister: appReplicaInformer.Lister(),
		appSynced: appReplicaInformer.Informer().HasSynced,

		workqueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "AppReplicas"),
	}

	klog.Info("Setting up event handlers")

	appReplicaInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueAppReplica,
		UpdateFunc: func(oldObj, newObj interface{}) {
			controller.enqueueAppReplica(newObj)
		},
	})

	deploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.handleObject,
		UpdateFunc: func(old, new interface{}) {
			newDepl := new.(*appsv1.Deployment)
			oldDepl := old.(*appsv1.Deployment)
			if newDepl.ResourceVersion == oldDepl.ResourceVersion {
				// Periodic resync will send update events for all known Deployments.
				// Two different versions of the same Deployment will always have different RVs.
				return
			}
			controller.handleObject(new)
		},
		DeleteFunc: controller.handleObject,
	})

	return controller

}

func (c *Controller) handleObject(obj interface{}) {
	var object metav1.Object
	var ok bool
	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object, invalid type"))
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
			return
		}
		klog.V(4).Infof("Recovered deleted object '%s' from tombstone", object.GetName())
	}
	klog.V(4).Infof("Processing object: %s", object.GetName())
	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
		// If this object is not owned by a AppReplicas, we should not do anything more
		// with it.
		if ownerRef.Kind != "AppReplica" {
			return
		}

		appReplica, err := c.appLister.AppReplicas(object.GetNamespace()).Get(ownerRef.Name)
		if err != nil {
			klog.V(4).Infof("ignoring orphaned object '%s/%s' of appReplica '%s'", object.GetNamespace(), object.GetName(), ownerRef.Name)
			return
		}

		c.enqueueAppReplica(appReplica)
		return
	}
}

func (c *Controller) enqueueAppReplica(obj interface{}) {

	var key string
	var err error

	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return

	}

	c.workqueue.Add(key)

}

func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	klog.Info("Starting AppReplicas controller")

	// Wait for the caches to be synced before starting workers
	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.deploymentSynced, c.appSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("Starting workers")
	// Launch two workers to process AppReplicas resources
	for i := 0; i < workers; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// AppReplicas resource to be synced.
		if err := c.syncHandler(key); err != nil {
			// Put the item back on the workqueue to handle any transient errors.
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the AppReplicas resource with this namespace/name
	appReplica, err := c.appLister.AppReplicas(namespace).Get(name)
	if err != nil {
		// The AppReplicas resource may no longer exist, in which case we stop
		// processing.

		fmt.Println(namespace, name, "not found")
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("AppReplicas '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}

	deploymentName := appReplica.Spec.DeploymentName
	if deploymentName == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		utilruntime.HandleError(fmt.Errorf("%s: deployment name must be specified", key))
		return nil
	}

	// Get the deployment with the name specified in AppReplicas.spec
	deployment, err := c.deploymentLister.Deployments(appReplica.Namespace).Get(deploymentName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		deployment, err = c.kubeclinetset.AppsV1().Deployments(appReplica.Namespace).Create(context.TODO(), newDeployment(appReplica), metav1.CreateOptions{})
	}

	// If an error occurs during Get/Create, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		fmt.Println(" deployment  get error", err)
		return err
	}

	// If the Deployment is not controlled by this AppReplicas resource, we should log
	// a warning to the event recorder and return error msg.
	if !metav1.IsControlledBy(deployment, appReplica) {
		msg := fmt.Sprintf(MessageResourceExists, deployment.Name)

		return fmt.Errorf("%s", msg)
	}

	// If this number of the replicas on the AppReplicas resource is specified, and the
	// number does not equal the current desired replicas on the Deployment, we
	// should update the Deployment resource.
	if (appReplica.Spec.Replicas != nil && *appReplica.Spec.Replicas != *deployment.Spec.Replicas) || len(deployment.Spec.Template.Spec.Containers) != 1 || deployment.Spec.Template.Spec.Containers[0].Image != appReplica.Spec.DeploymentImage {
		fmt.Printf("AppReplica %s replicas: %d, deployment replicas: %d", name, *appReplica.Spec.Replicas, *deployment.Spec.Replicas)
		fmt.Printf("AppReplica %s image: %s, deployment image: %s", name, appReplica.Spec.DeploymentImage, deployment.Spec.Template.Spec.Containers[0].Image)
		deployment, err = c.kubeclinetset.AppsV1().Deployments(appReplica.Namespace).Update(context.TODO(), newDeployment(appReplica), metav1.UpdateOptions{})
	}

	// If an error occurs during Update, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		fmt.Println(" deployment update error", err)
		return err
	}

	// Finally, we update the status block of the AppReplicas resource to reflect the
	// current state of the world
	err = c.updateAppReplicaStatus(appReplica, deployment)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) updateAppReplicaStatus(appReplica *appsv1alpha1.AppReplica, deployment *appsv1.Deployment) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	appReplicaCopy := appReplica.DeepCopy()
	appReplicaCopy.Status.AvailableReplicas = deployment.Status.AvailableReplicas
	// If the CustomResourceSubresources feature gate is not enabled,
	// we must use Update instead of UpdateStatus to update the Status block of the AppReplica resource.
	// UpdateStatus will not allow changes to the Spec of the resource,
	// which is ideal for ensuring nothing other than resource status has been updated.
	_, err := c.appclinetset.NextgenV1alpha1().AppReplicas(appReplica.Namespace).Update(context.TODO(), appReplicaCopy, metav1.UpdateOptions{})

	return err
}

func newDeployment(appReplica *appsv1alpha1.AppReplica) *appsv1.Deployment {
	labels := map[string]string{
		"app":        appReplica.Spec.DeploymentName,
		"controller": appReplica.Name,
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      appReplica.Spec.DeploymentName,
			Namespace: appReplica.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(appReplica, appsv1alpha1.SchemeGroupVersion.WithKind("AppReplica")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: appReplica.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  appReplica.Spec.DeploymentName,
							Image: appReplica.Spec.DeploymentImage,
						},
					},
				},
			},
		},
	}
}
