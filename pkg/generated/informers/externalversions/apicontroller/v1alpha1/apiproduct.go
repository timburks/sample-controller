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
	time "time"

	apicontrollerv1alpha1 "k8s.io/api-controller/pkg/apis/apicontroller/v1alpha1"
	versioned "k8s.io/api-controller/pkg/generated/clientset/versioned"
	internalinterfaces "k8s.io/api-controller/pkg/generated/informers/externalversions/internalinterfaces"
	v1alpha1 "k8s.io/api-controller/pkg/generated/listers/apicontroller/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ApiProductInformer provides access to a shared informer and lister for
// ApiProducts.
type ApiProductInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.ApiProductLister
}

type apiProductInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewApiProductInformer constructs a new informer for ApiProduct type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewApiProductInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredApiProductInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredApiProductInformer constructs a new informer for ApiProduct type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredApiProductInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ApicontrollerV1alpha1().ApiProducts(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.ApicontrollerV1alpha1().ApiProducts(namespace).Watch(context.TODO(), options)
			},
		},
		&apicontrollerv1alpha1.ApiProduct{},
		resyncPeriod,
		indexers,
	)
}

func (f *apiProductInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredApiProductInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *apiProductInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&apicontrollerv1alpha1.ApiProduct{}, f.defaultInformer)
}

func (f *apiProductInformer) Lister() v1alpha1.ApiProductLister {
	return v1alpha1.NewApiProductLister(f.Informer().GetIndexer())
}
