/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	appreplicav1alpha1 "sample-controller-k8s/pkg/apis/appreplica/v1alpha1"
	versioned "sample-controller-k8s/pkg/client/clientset/versioned"
	internalinterfaces "sample-controller-k8s/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "sample-controller-k8s/pkg/client/listers/appreplica/v1alpha1"
	time "time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// AppReplicaInformer provides access to a shared informer and lister for
// AppReplicas.
type AppReplicaInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.AppReplicaLister
}

type appReplicaInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewAppReplicaInformer constructs a new informer for AppReplica type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewAppReplicaInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredAppReplicaInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredAppReplicaInformer constructs a new informer for AppReplica type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredAppReplicaInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.NextgenV1alpha1().AppReplicas(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.NextgenV1alpha1().AppReplicas(namespace).Watch(context.TODO(), options)
			},
		},
		&appreplicav1alpha1.AppReplica{},
		resyncPeriod,
		indexers,
	)
}

func (f *appReplicaInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredAppReplicaInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *appReplicaInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&appreplicav1alpha1.AppReplica{}, f.defaultInformer)
}

func (f *appReplicaInformer) Lister() v1alpha1.AppReplicaLister {
	return v1alpha1.NewAppReplicaLister(f.Informer().GetIndexer())
}
