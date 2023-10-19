/*
Copyright 2017 The Kubernetes Authors.

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

package main

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	clientset "k8s.io/api-controller/pkg/generated/clientset/versioned"
	apischeme "k8s.io/api-controller/pkg/generated/clientset/versioned/scheme"
	informers "k8s.io/api-controller/pkg/generated/informers/externalversions/apicontroller/v1alpha1"
	listers "k8s.io/api-controller/pkg/generated/listers/apicontroller/v1alpha1"
)

const controllerAgentName = "api-controller"

const (
	// SuccessSynced is used as part of the Event 'reason' when an ApiProduct is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when an ApiProduct fails
	// to sync due to a Deployment of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Deployment already existing
	MessageResourceExists = "Resource %q already exists and is not managed by ApiProduct"
	// MessageResourceSynced is the message used for an Event fired when an ApiProduct
	// is synced successfully
	MessageProductSynced     = "ApiProduct synced successfully"
	MessageVersionSynced     = "ApiVersion synced successfully"
	MessageDescriptionSynced = "ApiDescription synced successfully"
	MessageDeploymentSynced  = "ApiDeployment synced successfully"
	MessageArtifactSynced    = "ApiArtifact synced successfully"
)

// Controller is the controller implementation for ApiProduct resources
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// sampleclientset is a clientset for our own API group
	sampleclientset clientset.Interface

	apiProductsLister     listers.ApiProductLister
	apiProductsSynced     cache.InformerSynced
	apiVersionsLister     listers.ApiVersionLister
	apiVersionsSynced     cache.InformerSynced
	apiDescriptionsLister listers.ApiDescriptionLister
	apiDescriptionsSynced cache.InformerSynced
	apiDeploymentsLister  listers.ApiDeploymentLister
	apiDeploymentsSynced  cache.InformerSynced
	apiArtifactsLister    listers.ApiArtifactLister
	apiArtifactsSynced    cache.InformerSynced

	// workqueue is a rate limited work queue. This is used to queue work to be
	// processed instead of performing it as soon as a change happens. This
	// means we can ensure we only process a fixed amount of resources at a
	// time, and makes it easy to ensure we are never processing the same item
	// simultaneously in two different workers.
	productsWorkqueue     workqueue.RateLimitingInterface
	versionsWorkqueue     workqueue.RateLimitingInterface
	descriptionsWorkqueue workqueue.RateLimitingInterface
	deploymentsWorkqueue  workqueue.RateLimitingInterface
	artifactsWorkqueue    workqueue.RateLimitingInterface

	// recorder is an event recorder for recording Event resources to the
	// Kubernetes API.
	recorder record.EventRecorder
}

// NewController returns a new sample controller
func NewController(
	ctx context.Context,
	kubeclientset kubernetes.Interface,
	sampleclientset clientset.Interface,
	apiProductInformer informers.ApiProductInformer,
	apiVersionInformer informers.ApiVersionInformer,
	apiDescriptionInformer informers.ApiDescriptionInformer,
	apiDeploymentInformer informers.ApiDeploymentInformer,
	apiArtifactInformer informers.ApiArtifactInformer,
) *Controller {
	logger := klog.FromContext(ctx)

	// Create event broadcaster
	// Add api-controller types to the default Kubernetes Scheme so Events can be
	// logged for api-controller types.
	utilruntime.Must(apischeme.AddToScheme(scheme.Scheme))
	logger.V(4).Info("Creating event broadcaster")

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeclientset:   kubeclientset,
		sampleclientset: sampleclientset,

		apiProductsLister:     apiProductInformer.Lister(),
		apiProductsSynced:     apiProductInformer.Informer().HasSynced,
		apiVersionsLister:     apiVersionInformer.Lister(),
		apiVersionsSynced:     apiVersionInformer.Informer().HasSynced,
		apiDescriptionsLister: apiDescriptionInformer.Lister(),
		apiDescriptionsSynced: apiDescriptionInformer.Informer().HasSynced,
		apiDeploymentsLister:  apiDeploymentInformer.Lister(),
		apiDeploymentsSynced:  apiDeploymentInformer.Informer().HasSynced,
		apiArtifactsLister:    apiArtifactInformer.Lister(),
		apiArtifactsSynced:    apiArtifactInformer.Informer().HasSynced,

		productsWorkqueue:     workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ApiProducts"),
		versionsWorkqueue:     workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ApiVersions"),
		descriptionsWorkqueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ApiDescriptionss"),
		deploymentsWorkqueue:  workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ApiDeployments"),
		artifactsWorkqueue:    workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "ApiArtifacts"),

		recorder: recorder,
	}

	logger.Info("Setting up event handlers")
	// Set up an event handler for when ApiProduct resources change
	apiProductInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueApiProduct,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueApiProduct(new)
		},
	})
	// Set up an event handler for when ApiVersion resources change
	apiVersionInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueApiVersion,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueApiVersion(new)
		},
	})
	// Set up an event handler for when ApiDescription resources change
	apiDescriptionInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueApiDescription,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueApiDescription(new)
		},
	})
	// Set up an event handler for when ApiDeployment resources change
	apiDeploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueApiDeployment,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueApiDeployment(new)
		},
	})
	// Set up an event handler for when ApiArtifact resources change
	apiArtifactInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueApiArtifact,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueApiArtifact(new)
		},
	})

	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(ctx context.Context, workers int) error {
	defer utilruntime.HandleCrash()
	defer c.productsWorkqueue.ShutDown()
	defer c.versionsWorkqueue.ShutDown()
	defer c.descriptionsWorkqueue.ShutDown()
	defer c.deploymentsWorkqueue.ShutDown()
	defer c.artifactsWorkqueue.ShutDown()
	logger := klog.FromContext(ctx)

	// Start the informer factories to begin populating the informer caches
	logger.Info("Starting controller")

	// Wait for the caches to be synced before starting workers
	logger.Info("Waiting for informer caches to sync")

	if ok := cache.WaitForCacheSync(ctx.Done(), c.apiProductsSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}
	if ok := cache.WaitForCacheSync(ctx.Done(), c.apiVersionsSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}
	if ok := cache.WaitForCacheSync(ctx.Done(), c.apiDescriptionsSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}
	if ok := cache.WaitForCacheSync(ctx.Done(), c.apiDeploymentsSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}
	if ok := cache.WaitForCacheSync(ctx.Done(), c.apiArtifactsSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	logger.Info("Starting workers", "count", workers)
	// Launch two workers to process ApiProduct resources
	for i := 0; i < workers; i++ {
		go wait.UntilWithContext(ctx, c.runProductWorker, time.Second)
	}
	// Launch two workers to process ApiVersion resources
	for i := 0; i < workers; i++ {
		go wait.UntilWithContext(ctx, c.runVersionWorker, time.Second)
	}
	// Launch two workers to process ApiDescription resources
	for i := 0; i < workers; i++ {
		go wait.UntilWithContext(ctx, c.runDescriptionWorker, time.Second)
	}
	// Launch two workers to process ApiDeployment resources
	for i := 0; i < workers; i++ {
		go wait.UntilWithContext(ctx, c.runDeploymentWorker, time.Second)
	}
	// Launch two workers to process ApiArtifact resources
	for i := 0; i < workers; i++ {
		go wait.UntilWithContext(ctx, c.runArtifactWorker, time.Second)
	}

	logger.Info("Started workers")
	<-ctx.Done()
	logger.Info("Shutting down workers")

	return nil
}

// enqueueApiProduct takes an ApiProduct resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than ApiProduct.
func (c *Controller) enqueueApiProduct(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.productsWorkqueue.Add(key)
}

// enqueueApiVersion takes an ApiVersion resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than ApiVersion.
func (c *Controller) enqueueApiVersion(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.versionsWorkqueue.Add(key)
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

// enqueueApiDeployment takes an ApiDeployment resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than ApiDeployment.
func (c *Controller) enqueueApiDeployment(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.deploymentsWorkqueue.Add(key)
}

// enqueueApiArtifact takes an ApiArtifact resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than ApiArtifact.
func (c *Controller) enqueueApiArtifact(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.artifactsWorkqueue.Add(key)
}
