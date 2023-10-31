package main

import (
	"fmt"

	myappv1alpha1 "sample-controller-k8s/pkg/apis/appreplica/v1alpha1"
	clinetset "sample-controller-k8s/pkg/client/clientset/versioned"
	informer "sample-controller-k8s/pkg/client/informers/externalversions"
	"time"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

func main() {

	// set up signals so we handle the shutdown signal gracefully
	ctx := signals.SetupSignalHandler()

	cfg, err := rest.InClusterConfig()

	if err != nil {
		fmt.Println(err)
		return
	}

	sampleClinetSet, err := clinetset.NewForConfig(cfg)

	if err != nil {
		fmt.Println(err)
		return
	}

	sampleInformerFactory := informer.NewSharedInformerFactory(sampleClinetSet, 30*time.Second)

	sampleInformer := sampleInformerFactory.Nextgen().V1alpha1().AppReplicas().Informer()

	sampleInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			newobj := obj.(myappv1alpha1.AppReplica)
			fmt.Println(newobj.Name)
		},
	})

	go sampleInformerFactory.Start(ctx.Done())

	if ok := cache.WaitForCacheSync(ctx.Done(), sampleInformer.HasSynced); !ok {
		fmt.Println("failed to wait for caches to sync")

		return
	}
	select {}
}
