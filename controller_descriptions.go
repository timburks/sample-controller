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
func (c *Controller) runDescriptionWorker(ctx context.Context) {
	for c.processNextDescriptionWorkItem(ctx) {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextDescriptionWorkItem(ctx context.Context) bool {
	obj, shutdown := c.descriptionsWorkqueue.Get()
	logger := klog.FromContext(ctx)

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.descriptionsWorkqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the descriptionsWorkqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the descriptionsWorkqueue and attempted again after a back-off
		// period.
		defer c.descriptionsWorkqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the descriptionsWorkqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// descriptionsWorkqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// descriptionsWorkqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the descriptionsWorkqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.descriptionsWorkqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in descriptionsWorkqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// ApiDescription resource to be synced.
		if err := c.syncDescriptionHandler(ctx, key); err != nil {
			// Put the item back on the descriptionsWorkqueue to handle any transient errors.
			c.descriptionsWorkqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.descriptionsWorkqueue.Forget(obj)
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
// converge the two. It then updates the Status block of the ApiDescription resource
// with the current status of the resource.
func (c *Controller) syncDescriptionHandler(ctx context.Context, key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the ApiDescription resource with this namespace/name
	apiDescription, err := c.apiDescriptionsLister.ApiDescriptions(namespace).Get(name)
	if err != nil {
		// The ApiDescription resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("ApiDescription '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}

	// Finally, we update the status block of the ApiDescription resource to reflect the
	// current state of the world
	err = c.updateApiDescriptionStatus(apiDescription)
	if err != nil {
		return err
	}

	c.recorder.Event(apiDescription, corev1.EventTypeNormal, SuccessSynced, MessageDescriptionSynced)
	return nil
}

func (c *Controller) updateApiDescriptionStatus(apiDescription *apiv1alpha1.ApiDescription) error {
	message := "ok description"
	if apiDescription.Status.Message == message {
		return nil
	}
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	apiDescriptionCopy := apiDescription.DeepCopy()
	apiDescriptionCopy.Status.Message = message
	// If the CustomResourceSubresources feature gate is not enabled,
	// we must use Update instead of UpdateStatus to update the Status block of the ApiDescription resource.
	// UpdateStatus will not allow changes to the Spec of the resource,
	// which is ideal for ensuring nothing other than resource status has been updated.
	_, err := c.sampleclientset.ApicontrollerV1alpha1().ApiDescriptions(apiDescription.Namespace).UpdateStatus(context.TODO(), apiDescriptionCopy, metav1.UpdateOptions{})
	return err
}

// enqueueApiDescription takes an ApiDescription resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than ApiDescription.
func (c *Controller) enqueueApiDescription(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.descriptionsWorkqueue.Add(key)
}
