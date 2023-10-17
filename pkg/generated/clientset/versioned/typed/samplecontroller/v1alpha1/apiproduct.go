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

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
	v1alpha1 "k8s.io/api-controller/pkg/apis/samplecontroller/v1alpha1"
	scheme "k8s.io/api-controller/pkg/generated/clientset/versioned/scheme"
)

// ApiProductsGetter has a method to return a ApiProductInterface.
// A group's client should implement this interface.
type ApiProductsGetter interface {
	ApiProducts(namespace string) ApiProductInterface
}

// ApiProductInterface has methods to work with ApiProduct resources.
type ApiProductInterface interface {
	Create(ctx context.Context, apiProduct *v1alpha1.ApiProduct, opts v1.CreateOptions) (*v1alpha1.ApiProduct, error)
	Update(ctx context.Context, apiProduct *v1alpha1.ApiProduct, opts v1.UpdateOptions) (*v1alpha1.ApiProduct, error)
	UpdateStatus(ctx context.Context, apiProduct *v1alpha1.ApiProduct, opts v1.UpdateOptions) (*v1alpha1.ApiProduct, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.ApiProduct, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.ApiProductList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ApiProduct, err error)
	ApiProductExpansion
}

// apiProducts implements ApiProductInterface
type apiProducts struct {
	client rest.Interface
	ns     string
}

// newApiProducts returns a ApiProducts
func newApiProducts(c *SamplecontrollerV1alpha1Client, namespace string) *apiProducts {
	return &apiProducts{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the apiProduct, and returns the corresponding apiProduct object, and an error if there is any.
func (c *apiProducts) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ApiProduct, err error) {
	result = &v1alpha1.ApiProduct{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("apiproducts").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ApiProducts that match those selectors.
func (c *apiProducts) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ApiProductList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.ApiProductList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("apiproducts").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested apiProducts.
func (c *apiProducts) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("apiproducts").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a apiProduct and creates it.  Returns the server's representation of the apiProduct, and an error, if there is any.
func (c *apiProducts) Create(ctx context.Context, apiProduct *v1alpha1.ApiProduct, opts v1.CreateOptions) (result *v1alpha1.ApiProduct, err error) {
	result = &v1alpha1.ApiProduct{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("apiproducts").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(apiProduct).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a apiProduct and updates it. Returns the server's representation of the apiProduct, and an error, if there is any.
func (c *apiProducts) Update(ctx context.Context, apiProduct *v1alpha1.ApiProduct, opts v1.UpdateOptions) (result *v1alpha1.ApiProduct, err error) {
	result = &v1alpha1.ApiProduct{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("apiproducts").
		Name(apiProduct.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(apiProduct).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *apiProducts) UpdateStatus(ctx context.Context, apiProduct *v1alpha1.ApiProduct, opts v1.UpdateOptions) (result *v1alpha1.ApiProduct, err error) {
	result = &v1alpha1.ApiProduct{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("apiproducts").
		Name(apiProduct.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(apiProduct).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the apiProduct and deletes it. Returns an error if one occurs.
func (c *apiProducts) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("apiproducts").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *apiProducts) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("apiproducts").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched apiProduct.
func (c *apiProducts) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ApiProduct, err error) {
	result = &v1alpha1.ApiProduct{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("apiproducts").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
