package main

import (
	"fmt"

	clinetset "sample-controller-k8s/pkg/clinet/clientset/versioned"
	informer "sample-controller-k8s/pkg/clinet/informers/externalversions"
	"time"

	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

func main() {

	// set up signals so we handle the shutdown signal gracefully
	ctx := signals.SetupSignalHandler()

	cfg, err := clientcmd.BuildConfigFromFlags("", "/root/.kube/config")

	if err != nil {

		fmt.Println(err)

		err = nil

		cfg, err = rest.InClusterConfig()

		if err != nil {
			fmt.Println(err)
			return
		}
	}

	kubernetesClinet, err := kubernetes.NewForConfig(cfg)

	if err != nil {
		fmt.Println(err)
		return
	}

	sampleClinetSet, err := clinetset.NewForConfig(cfg)

	if err != nil {
		fmt.Println(err)
		return
	}

	kubeShareInformerFactort := kubeinformers.NewSharedInformerFactory(kubernetesClinet, 30*time.Second)
	sampleInformerFactory := informer.NewSharedInformerFactory(sampleClinetSet, 30*time.Second)

	controller := NewController(ctx,
		kubernetesClinet,
		sampleClinetSet,
		kubeShareInformerFactort.Apps().V1().Deployments(),
		sampleInformerFactory.Nextgen().V1alpha1().MyAppDeployments(),
	)

	kubeShareInformerFactort.Start(ctx.Done())
	//sampleInformerFactory.Start(ctx.Done())

	if err := controller.Run(ctx, 2); err != nil {

		fmt.Println("Error running Controller: ", err)

		klog.FlushAndExit(klog.ExitFlushTimeout, 1)

	}
}
