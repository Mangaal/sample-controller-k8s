package main

import (
	"fmt"

	clinetset "sample-controller-k8s/pkg/client/clientset/versioned"
	informer "sample-controller-k8s/pkg/client/informers/externalversions"
	"time"

	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

func main() {

	// set up signals so we handle the shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	cfg, err := rest.InClusterConfig()

	if err != nil {
		fmt.Println(err)
		return
	}

	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		klog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	appreplicaClinetSet, err := clinetset.NewForConfig(cfg)

	if err != nil {
		klog.Fatalf("Error building appreplica clientset: %s", err.Error())

	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, time.Second*30)

	appreplicaInformerFactory := informer.NewSharedInformerFactory(appreplicaClinetSet, 30*time.Second)

	controller := NewController(kubeClient, appreplicaClinetSet,
		appreplicaInformerFactory.Nextgen().V1alpha1().AppReplicas(),
		kubeInformerFactory.Apps().V1().Deployments())

	kubeInformerFactory.Start(stopCh.Done())
	appreplicaInformerFactory.Start(stopCh.Done())

	if err = controller.Run(2, stopCh.Done()); err != nil {
		klog.Fatalf("Error running controller: %s", err.Error())
	}

}
