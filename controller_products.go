package main

import (
	"context"
	"fmt"

	apiv1alpha1 "k8s.io/api-controller/pkg/apis/apicontroller/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
)

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runProductWorker(ctx context.Context) {
	for c.processNextProductWorkItem(ctx) {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextProductWorkItem(ctx context.Context) bool {
	obj, shutdown := c.productsWorkqueue.Get()
	logger := klog.FromContext(ctx)

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.productsWorkqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the productsWorkqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the productsWorkqueue and attempted again after a back-off
		// period.
		defer c.productsWorkqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the productsWorkqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// productsWorkqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// productsWorkqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the productsWorkqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.productsWorkqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in productsWorkqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// ApiProduct resource to be synced.
		if err := c.syncProductHandler(ctx, key); err != nil {
			// Put the item back on the productsWorkqueue to handle any transient errors.
			c.productsWorkqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.productsWorkqueue.Forget(obj)
		logger.Info("Successfully synced", "resourceName", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the ApiProduct resource
// with the current status of the resource.
func (c *Controller) syncProductHandler(ctx context.Context, key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the ApiProduct resource with this namespace/name
	apiProduct, err := c.apiProductsLister.ApiProducts(namespace).Get(name)
	if err != nil {
		// The ApiProduct resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("ApiProduct '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}

	// Finally, we update the status block of the ApiProduct resource to reflect the
	// current state of the world
	err = c.updateApiProductStatus(apiProduct)
	if err != nil {
		return err
	}

	c.recorder.Event(apiProduct, corev1.EventTypeNormal, SuccessSynced, MessageProductSynced)
	return nil
}

func (c *Controller) updateApiProductStatus(apiProduct *apiv1alpha1.ApiProduct) error {
	message := "ok product"
	if apiProduct.Status.Message == message {
		return nil
	}
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	apiProductCopy := apiProduct.DeepCopy()
	apiProductCopy.Status.Message = message
	// If the CustomResourceSubresources feature gate is not enabled,
	// we must use Update instead of UpdateStatus to update the Status block of the ApiProduct resource.
	// UpdateStatus will not allow changes to the Spec of the resource,
	// which is ideal for ensuring nothing other than resource status has been updated.
	_, err := c.sampleclientset.ApicontrollerV1alpha1().ApiProducts(apiProduct.Namespace).UpdateStatus(context.TODO(), apiProductCopy, metav1.UpdateOptions{})
	return err
}
