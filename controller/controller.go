package controller

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sample-controller-k8s/pkg/api/myappdeployment/v1alpha1"

	//clinetset "sample-controller-k8s/pkg/clinet/clientset/versioned"

	"k8s.io/client-go/kubernetes/scheme"

	//samplescheme "sample-controller-k8s/pkg/clinet/clientset/versioned/scheme"
	//informers "sample-controller-k8s/pkg/clinet/informers/externalversions/myappdeployment/v1alpha1"
	//listers "sample-controller-k8s/pkg/clinet/listers/myappdeployment/v1alpha1"
	"time"

	"golang.org/x/time/rate"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"

	appsinformers "k8s.io/client-go/informers/apps/v1"
	"k8s.io/client-go/kubernetes"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"

	appslisters "k8s.io/client-go/listers/apps/v1"
)

type Controller struct {
	KubeClinetSet kubernetes.Interface

	SampleClinetSet clinetset.Interface

	DeploymentsListers appslisters.DeploymentLister
	DeploymentSynced   cache.InformerSynced

	MyAppDeploymentLister listers.MyAppDeploymentLister
	MyAppDeploymentSynced cache.InformerSynced

	Workqueue workqueue.RateLimitingInterface

	recorder record.EventRecorder
}

const controllerAgentName = "sample-controller"

func NewController(ctx context.Context,
	kubeclientset kubernetes.Interface,
	sampleclientset clinetset.Interface,
	deploymentInformer appsinformers.DeploymentInformer,
	myAppDeploymentInformer informers.MyAppDeploymentInformer) *Controller {

	utilruntime.Must(samplescheme.AddToScheme(scheme.Scheme))
	fmt.Println("Creating event broadcaster is call")

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	ratelimiter := workqueue.NewMaxOfRateLimiter(
		workqueue.NewItemExponentialFailureRateLimiter(5*time.Millisecond, 1000*time.Second),
		&workqueue.BucketRateLimiter{Limiter: rate.NewLimiter(rate.Limit(50), 300)},
	)

	controller := &Controller{
		KubeClinetSet:         kubeclientset,
		SampleClinetSet:       sampleclientset,
		DeploymentsListers:    deploymentInformer.Lister(),
		DeploymentSynced:      deploymentInformer.Informer().HasSynced,
		MyAppDeploymentLister: myAppDeploymentInformer.Lister(),
		MyAppDeploymentSynced: myAppDeploymentInformer.Informer().HasSynced,
		Workqueue:             workqueue.NewRateLimitingQueue(ratelimiter),
		recorder:              recorder,
	}

	myAppDeploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueMyAppDeploymentResource,
		UpdateFunc: func(oldObj, newObj interface{}) {
			controller.enqueueMyAppDeploymentResource(newObj)
		},
	})

	deploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.handleObject,
		UpdateFunc: func(oldObj, newObj interface{}) {
			newDepl := newObj.(*appsv1.Deployment)
			oldDepl := oldObj.(*appsv1.Deployment)

			if !reflect.DeepEqual(newDepl.Spec, oldDepl.Spec) {

				controller.handleObject(newObj)

			}

		},
		DeleteFunc: controller.handleObject,
	})

	return controller
}

func (c *Controller) Run(ctx context.Context, workers int) error {
	defer c.Workqueue.ShutDown()

	fmt.Println("Starting worker foo")

	// Wait for the caches to be synced before starting workers
	fmt.Println("Waiting for informer caches to sync")

	if ok := cache.WaitForCacheSync(ctx.Done(), c.DeploymentSynced, c.MyAppDeploymentSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	fmt.Println("Starting workers", "count", workers)
	// Launch two workers to process Foo resources
	for i := 0; i < workers; i++ {
		go wait.UntilWithContext(ctx, c.runWorker, time.Second)
	}

	fmt.Println("Started workers")
	<-ctx.Done()
	fmt.Println("Shutting down workers")

	return nil

}

func (c *Controller) runWorker(ctx context.Context) {

	for c.processNextItem(ctx) {

	}

}

func (c *Controller) processNextItem(ctx context.Context) bool {

	obj, sutdown := c.Workqueue.Get()

	if sutdown {
		return false

	}

	err := func(obj interface{}) error {

		defer c.Workqueue.Done(obj)

		var key string

		var ok bool

		if key, ok = obj.(string); !ok {
			c.Workqueue.Forget(obj)

			fmt.Printf("Error: expected string in workqueue but got %#v", obj)
			return nil

		}

		if err := c.syncHandeler(ctx, key); err != nil {

			c.Workqueue.AddRateLimited(key)

			return fmt.Errorf("error syncing %s : %s . renquing ", key, err.Error())

		}

		c.Workqueue.Forget(key)

		fmt.Printf("Successfully sync , resource %s", key)

		return nil

	}(obj)

	if err != nil {

		fmt.Println(err.Error())

		return true
	}

	return true

}

func (c *Controller) syncHandeler(ctx context.Context, key string) error {

	name, namespace, err := cache.SplitMetaNamespaceKey(key)

	if err != nil {
		fmt.Println("Error geting name and namespace from key")
		return nil
	}
	myresource, err := c.MyAppDeploymentLister.MyAppDeployments(namespace).Get(name)

	if err != nil {

		fmt.Println("Error getting myappdeployment resource ", myresource, err)

		return err
	}

	deploymentName := myresource.Spec.DeploymentName

	if deploymentName == "" {
		fmt.Println("Deployment must be specfiled ")
		return nil
	}

	deployment, err := c.DeploymentsListers.Deployments(namespace).Get(deploymentName)

	if err != nil {

		var iferror error

		fmt.Printf("Error getting deployment %s for  resource %s. %v .\n Creating a new deployment \n", deploymentName, myresource.Name, err)

		deployment, iferror = c.KubeClinetSet.AppsV1().Deployments(namespace).Create(ctx, newDeployment(myresource), metav1.CreateOptions{})

		if iferror != nil {

			fmt.Printf("Error Creating Deployment %s for resync resource %s", deploymentName, myresource.Name)

			return iferror
		}

	}

	if !metav1.IsControlledBy(deployment, myresource) {
		errormsg := fmt.Sprintf("The deployment %s is not controller by resource %s \n", deploymentName, myresource.Name)
		return errors.New(errormsg)
	}

	if myresource.Spec.Replicas != nil && *myresource.Spec.Replicas != *deployment.Spec.Replicas {
		var iferror error
		deployment, iferror = c.KubeClinetSet.AppsV1().Deployments(namespace).Update(context.TODO(), newDeployment(myresource), metav1.UpdateOptions{})

		if iferror != nil {
			return iferror

		}

	} else {

		if len(deployment.Spec.Template.Spec.Containers) != 1 && deployment.Spec.Template.Spec.Containers[0].Name != myresource.Spec.DeploymentName && deployment.Spec.Template.Spec.Containers[0].Image != myresource.Spec.DeploymentImage {
			var iferror error
			deployment, iferror = c.KubeClinetSet.AppsV1().Deployments(namespace).Update(context.TODO(), newDeployment(myresource), metav1.UpdateOptions{})

			if iferror != nil {
				return iferror

			}
		}

	}

	err = c.UpdateMyAppDeployment(myresource, deployment)

	if err != nil {
		return err
	}

	return nil

}

func (c *Controller) UpdateMyAppDeployment(myresource *v1alpha1.MyAppDeployment, deployment *appsv1.Deployment) error {

	myresourceCopy := myresource.DeepCopy()

	myresourceCopy.Status.AvailableReplica = *deployment.Spec.Replicas

	_, err := c.SampleClinetSet.NextgenV1alpha1().MyAppDeployments(myresource.Namespace).UpdateStatus(context.TODO(), myresourceCopy, metav1.UpdateOptions{})

	if err != nil {
		return err
	}

	return nil

}

func (c *Controller) handleObject(obj interface{}) {

	var object metav1.Object
	var ok bool

	if object, ok = obj.(metav1.Object); !ok {

		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)

		if !ok {
			fmt.Println("Error decoding object, inavlid type")

			return
		}

		object, ok = tombstone.Obj.(metav1.Object)

		if !ok {
			fmt.Println("Error decoding tombstone object, inavlid type")

			return
		}
	}

	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {

		if ownerRef.Kind != "MyAppDeployment" {
			return
		}

		mydepoyment, err := c.MyAppDeploymentLister.MyAppDeployments(object.GetNamespace()).Get(ownerRef.Name)

		if err != nil {
			fmt.Println("Ignoring orphan object ", object.GetName())
		}

		c.enqueueMyAppDeploymentResource(mydepoyment)

		return

	}
}

func (c *Controller) enqueueMyAppDeploymentResource(obj interface{}) {
	var key string
	var err error

	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		fmt.Println(err)
		return
	}

	c.Workqueue.Add(key)
}

func newDeployment(myDeploymentApp *v1alpha1.MyAppDeployment) *appsv1.Deployment {

	label := map[string]string{
		"app":        myDeploymentApp.Spec.DeploymentName,
		"controller": myDeploymentApp.Name,
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      myDeploymentApp.Spec.DeploymentName,
			Namespace: myDeploymentApp.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(myDeploymentApp, v1alpha1.SchemeGroupVersion.WithKind("MyAppDeployment")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: myDeploymentApp.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: label,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: label,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  myDeploymentApp.Spec.DeploymentName,
							Image: myDeploymentApp.Spec.DeploymentImage,
						},
					},
				},
			},
		},
	}

}
