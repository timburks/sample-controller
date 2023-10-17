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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	v1alpha1 "k8s.io/sample-controller/pkg/apis/samplecontroller/v1alpha1"
)

// ApiProductLister helps list ApiProducts.
// All objects returned here must be treated as read-only.
type ApiProductLister interface {
	// List lists all ApiProducts in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.ApiProduct, err error)
	// ApiProducts returns an object that can list and get ApiProducts.
	ApiProducts(namespace string) ApiProductNamespaceLister
	ApiProductListerExpansion
}

// apiProductLister implements the ApiProductLister interface.
type apiProductLister struct {
	indexer cache.Indexer
}

// NewApiProductLister returns a new ApiProductLister.
func NewApiProductLister(indexer cache.Indexer) ApiProductLister {
	return &apiProductLister{indexer: indexer}
}

// List lists all ApiProducts in the indexer.
func (s *apiProductLister) List(selector labels.Selector) (ret []*v1alpha1.ApiProduct, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ApiProduct))
	})
	return ret, err
}

// ApiProducts returns an object that can list and get ApiProducts.
func (s *apiProductLister) ApiProducts(namespace string) ApiProductNamespaceLister {
	return apiProductNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ApiProductNamespaceLister helps list and get ApiProducts.
// All objects returned here must be treated as read-only.
type ApiProductNamespaceLister interface {
	// List lists all ApiProducts in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.ApiProduct, err error)
	// Get retrieves the ApiProduct from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.ApiProduct, error)
	ApiProductNamespaceListerExpansion
}

// apiProductNamespaceLister implements the ApiProductNamespaceLister
// interface.
type apiProductNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ApiProducts in the indexer for a given namespace.
func (s apiProductNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.ApiProduct, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ApiProduct))
	})
	return ret, err
}

// Get retrieves the ApiProduct from the indexer for a given namespace and name.
func (s apiProductNamespaceLister) Get(name string) (*v1alpha1.ApiProduct, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("apiproduct"), name)
	}
	return obj.(*v1alpha1.ApiProduct), nil
}
