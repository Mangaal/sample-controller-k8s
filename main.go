package main

import (
	"fmt"

	clinetset "sample-controller-k8s/pkg/client/clientset/versioned"
	informer "sample-controller-k8s/pkg/client/informers/externalversions"
	"time"

	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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

	kubeShareInformerFactort.Start(ctx.Done())
	sampleInformerFactory.Start(ctx.Done())

}
