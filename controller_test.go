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
	"reflect"
	"testing"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/diff"
	kubeinformers "k8s.io/client-go/informers"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	core "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2/ktesting"

	apicontroller "k8s.io/api-controller/pkg/apis/apicontroller/v1alpha1"
	"k8s.io/api-controller/pkg/generated/clientset/versioned/fake"
	informers "k8s.io/api-controller/pkg/generated/informers/externalversions"
)

var (
	alwaysReady        = func() bool { return true }
	noResyncPeriodFunc = func() time.Duration { return 0 }
)

type fixture struct {
	t *testing.T

	client     *fake.Clientset
	kubeclient *k8sfake.Clientset
	// Objects to put in the store.
	apiProductLister []*apicontroller.ApiProduct
	// Actions expected to happen on the client.
	kubeactions []core.Action
	actions     []core.Action
	// Objects from here preloaded into NewSimpleFake.
	kubeobjects []runtime.Object
	objects     []runtime.Object
}

func newFixture(t *testing.T) *fixture {
	f := &fixture{}
	f.t = t
	f.objects = []runtime.Object{}
	f.kubeobjects = []runtime.Object{}
	return f
}

func newApiProduct(name string, replicas *int32) *apicontroller.ApiProduct {
	return &apicontroller.ApiProduct{
		TypeMeta: metav1.TypeMeta{APIVersion: apicontroller.SchemeGroupVersion.String()},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: metav1.NamespaceDefault,
		},
		Spec: apicontroller.ApiProductSpec{
			Title:       "Sample API Product",
			Description: "A sample API Product.",
		},
	}
}

func (f *fixture) newController(ctx context.Context) (*Controller, informers.SharedInformerFactory, kubeinformers.SharedInformerFactory) {
	f.client = fake.NewSimpleClientset(f.objects...)
	f.kubeclient = k8sfake.NewSimpleClientset(f.kubeobjects...)

	i := informers.NewSharedInformerFactory(f.client, noResyncPeriodFunc())
	k8sI := kubeinformers.NewSharedInformerFactory(f.kubeclient, noResyncPeriodFunc())

	c := NewController(ctx, f.kubeclient, f.client,
		i.Apicontroller().V1alpha1().ApiProducts(),
		i.Apicontroller().V1alpha1().ApiVersions(),
		i.Apicontroller().V1alpha1().ApiDescriptions(),
		i.Apicontroller().V1alpha1().ApiDeployments(),
	)

	c.apiProductsSynced = alwaysReady
	c.recorder = &record.FakeRecorder{}

	for _, f := range f.apiProductLister {
		i.Apicontroller().V1alpha1().ApiProducts().Informer().GetIndexer().Add(f)
	}

	return c, i, k8sI
}

func (f *fixture) run(ctx context.Context, apiProductName string) {
	f.runController(ctx, apiProductName, true, false)
}

func (f *fixture) runController(ctx context.Context, apiProductName string, startInformers bool, expectError bool) {
	c, i, k8sI := f.newController(ctx)
	if startInformers {
		i.Start(ctx.Done())
		k8sI.Start(ctx.Done())
	}

	err := c.syncProductHandler(ctx, apiProductName)
	if !expectError && err != nil {
		f.t.Errorf("error syncing apiProduct: %v", err)
	} else if expectError && err == nil {
		f.t.Error("expected error syncing apiProduct, got nil")
	}

	actions := filterInformerActions(f.client.Actions())
	for i, action := range actions {
		if len(f.actions) < i+1 {
			f.t.Errorf("%d unexpected actions: %+v", len(actions)-len(f.actions), actions[i:])
			break
		}

		expectedAction := f.actions[i]
		checkAction(expectedAction, action, f.t)
	}

	if len(f.actions) > len(actions) {
		f.t.Errorf("%d additional expected actions:%+v", len(f.actions)-len(actions), f.actions[len(actions):])
	}

}

// checkAction verifies that expected and actual actions are equal and both have
// same attached resources
func checkAction(expected, actual core.Action, t *testing.T) {
	if !(expected.Matches(actual.GetVerb(), actual.GetResource().Resource) && actual.GetSubresource() == expected.GetSubresource()) {
		t.Errorf("Expected\n\t%#v\ngot\n\t%#v", expected, actual)
		return
	}

	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		t.Errorf("Action has wrong type. Expected: %t. Got: %t", expected, actual)
		return
	}

	switch a := actual.(type) {
	case core.CreateActionImpl:
		e, _ := expected.(core.CreateActionImpl)
		expObject := e.GetObject()
		object := a.GetObject()

		if !reflect.DeepEqual(expObject, object) {
			t.Errorf("Action %s %s has wrong object\nDiff:\n %s",
				a.GetVerb(), a.GetResource().Resource, diff.ObjectGoPrintSideBySide(expObject, object))
		}
	case core.UpdateActionImpl:
		e, _ := expected.(core.UpdateActionImpl)
		expObject := e.GetObject()
		object := a.GetObject()

		if !reflect.DeepEqual(expObject, object) {
			t.Errorf("Action %s %s has wrong object\nDiff:\n %s",
				a.GetVerb(), a.GetResource().Resource, diff.ObjectGoPrintSideBySide(expObject, object))
		}
	case core.PatchActionImpl:
		e, _ := expected.(core.PatchActionImpl)
		expPatch := e.GetPatch()
		patch := a.GetPatch()

		if !reflect.DeepEqual(expPatch, patch) {
			t.Errorf("Action %s %s has wrong patch\nDiff:\n %s",
				a.GetVerb(), a.GetResource().Resource, diff.ObjectGoPrintSideBySide(expPatch, patch))
		}
	default:
		t.Errorf("Uncaptured Action %s %s, you should explicitly add a case to capture it",
			actual.GetVerb(), actual.GetResource().Resource)
	}
}

// filterInformerActions filters list and watch actions for testing resources.
// Since list and watch don't change resource state we can filter it to lower
// nose level in our tests.
func filterInformerActions(actions []core.Action) []core.Action {
	ret := []core.Action{}
	for _, action := range actions {
		if len(action.GetNamespace()) == 0 &&
			(action.Matches("list", "apiproducts") ||
				action.Matches("watch", "apiproducts")) {
			continue
		}
		ret = append(ret, action)
	}

	return ret
}

func (f *fixture) expectUpdateApiProductStatusAction(apiProduct *apicontroller.ApiProduct) {
	action := core.NewUpdateSubresourceAction(schema.GroupVersionResource{Resource: "apiproducts"}, "status", apiProduct.Namespace, apiProduct)
	f.actions = append(f.actions, action)
}

func getKey(apiProduct *apicontroller.ApiProduct, t *testing.T) string {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(apiProduct)
	if err != nil {
		t.Errorf("Unexpected error getting key for apiProduct %v: %v", apiProduct.Name, err)
		return ""
	}
	return key
}

func TestDoNothing(t *testing.T) {
	f := newFixture(t)
	apiProduct := newApiProduct("test", int32Ptr(1))
	_, ctx := ktesting.NewTestContext(t)

	f.apiProductLister = append(f.apiProductLister, apiProduct)
	f.objects = append(f.objects, apiProduct)

	apiProduct.Status.Message = "ok"
	f.expectUpdateApiProductStatusAction(apiProduct)
	f.run(ctx, getKey(apiProduct, t))
}

func int32Ptr(i int32) *int32 { return &i }
