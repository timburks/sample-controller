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
	"k8s.io/api-controller/pkg/generated/informers/externalversions/apicontroller"
	listers "k8s.io/api-controller/pkg/generated/listers/apicontroller/v1alpha1"
)

const controllerAgentName = "api-controller"

const (
	// SuccessSynced is used as part of the Event 'reason' when an ApiProduct is synced
	SuccessSynced = "Synced"

	// MessageResourceSynced is the message used for an Event fired when a custom resource
	// is synced successfully
	MessageProductSynced     = "ApiProduct synced successfully"
	MessageVersionSynced     = "ApiVersion synced successfully"
	MessageDescriptionSynced = "ApiDescription synced successfully"
	MessageDeploymentSynced  = "ApiDeployment synced successfully"
	MessageArtifactSynced    = "ApiArtifact synced successfully"
)

// Controller is the controller implementation for all of our custom resources
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// clientset is a clientset for our own API group
	clientset clientset.Interface

	apiProductsLister     listers.ApiProductLister
	apiVersionsLister     listers.ApiVersionLister
	apiDescriptionsLister listers.ApiDescriptionLister
	apiDeploymentsLister  listers.ApiDeploymentLister
	apiArtifactsLister    listers.ApiArtifactLister

	apiProductsSynced     cache.InformerSynced
	apiVersionsSynced     cache.InformerSynced
	apiDescriptionsSynced cache.InformerSynced
	apiDeploymentsSynced  cache.InformerSynced
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

// NewController returns a new controller
func NewController(
	ctx context.Context,
	kubeclientset kubernetes.Interface,
	clientset clientset.Interface,
	informers apicontroller.Interface,
) *Controller {
	logger := klog.FromContext(ctx)

	apiProductInformer := informers.V1alpha1().ApiProducts()
	apiVersionInformer := informers.V1alpha1().ApiVersions()
	apiDescriptionInformer := informers.V1alpha1().ApiDescriptions()
	apiDeploymentInformer := informers.V1alpha1().ApiDeployments()
	apiArtifactInformer := informers.V1alpha1().ApiArtifacts()

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
		kubeclientset: kubeclientset,
		clientset:     clientset,

		apiProductsLister:     apiProductInformer.Lister(),
		apiVersionsLister:     apiVersionInformer.Lister(),
		apiDescriptionsLister: apiDescriptionInformer.Lister(),
		apiDeploymentsLister:  apiDeploymentInformer.Lister(),
		apiArtifactsLister:    apiArtifactInformer.Lister(),

		apiProductsSynced:     apiProductInformer.Informer().HasSynced,
		apiVersionsSynced:     apiVersionInformer.Informer().HasSynced,
		apiDescriptionsSynced: apiDescriptionInformer.Informer().HasSynced,
		apiDeploymentsSynced:  apiDeploymentInformer.Informer().HasSynced,
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
	if ok := cache.WaitForCacheSync(ctx.Done(),
		c.apiProductsSynced,
		c.apiVersionsSynced,
		c.apiDescriptionsSynced,
		c.apiDeploymentsSynced,
		c.apiArtifactsSynced,
	); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	logger.Info("Starting workers", "count", workers)
	// Launch workers to process each custom resource type
	for i := 0; i < workers; i++ {
		go wait.UntilWithContext(ctx, c.runProductWorker, time.Second)
		go wait.UntilWithContext(ctx, c.runVersionWorker, time.Second)
		go wait.UntilWithContext(ctx, c.runDescriptionWorker, time.Second)
		go wait.UntilWithContext(ctx, c.runDeploymentWorker, time.Second)
		go wait.UntilWithContext(ctx, c.runArtifactWorker, time.Second)
	}

	logger.Info("Started workers")
	<-ctx.Done()
	logger.Info("Shutting down workers")

	return nil
}
