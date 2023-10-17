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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApiDeployment is a specification for a ApiDeployment resource
type ApiDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApiDeploymentSpec   `json:"spec"`
	Status ApiDeploymentStatus `json:"status"`
}

// ApiDeploymentSpec is the spec for a ApiDeployment resource
type ApiDeploymentSpec struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Parent      string `json:"parent"`
}

// ApiDeploymentStatus is the status for a ApiDeployment resource
type ApiDeploymentStatus struct {
	Message string `json:"message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApiDeploymentList is a list of ApiDeployment resources
type ApiDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ApiDeployment `json:"items"`
}
