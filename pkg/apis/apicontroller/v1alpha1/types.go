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

// ApiProduct is a specification for a ApiProduct resource
type ApiProduct struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApiProductSpec   `json:"spec"`
	Status ApiProductStatus `json:"status"`
}

// ApiProductSpec is the spec for a ApiProduct resource
type ApiProductSpec struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// ApiProductStatus is the status for a ApiProduct resource
type ApiProductStatus struct {
	Message string `json:"message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApiProductList is a list of ApiProduct resources
type ApiProductList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ApiProduct `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApiVersion is a specification for a ApiVersion resource
type ApiVersion struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApiVersionSpec   `json:"spec"`
	Status ApiVersionStatus `json:"status"`
}

// ApiVersionSpec is the spec for a ApiVersion resource
type ApiVersionSpec struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Parent      string `json:"parent"`
}

// ApiVersionStatus is the status for a ApiVersion resource
type ApiVersionStatus struct {
	Message string `json:"message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApiVersionList is a list of ApiVersion resources
type ApiVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ApiVersion `json:"items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApiDescription is a specification for a ApiDescription resource
type ApiDescription struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApiDescriptionSpec   `json:"spec"`
	Status ApiDescriptionStatus `json:"status"`
}

// ApiDescriptionSpec is the spec for a ApiDescription resource
type ApiDescriptionSpec struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Parent      string `json:"parent"`
}

// ApiDescriptionStatus is the status for a ApiDescription resource
type ApiDescriptionStatus struct {
	Message string `json:"message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApiDescriptionList is a list of ApiDescription resources
type ApiDescriptionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ApiDescription `json:"items"`
}

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

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApiArtifact is a specification for a ApiArtifact resource
type ApiArtifact struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApiArtifactSpec   `json:"spec"`
	Status ApiArtifactStatus `json:"status"`
}

// ApiArtifactSpec is the spec for a ApiArtifact resource
type ApiArtifactSpec struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Parent      string `json:"parent"`
}

// ApiArtifactStatus is the status for a ApiArtifact resource
type ApiArtifactStatus struct {
	Message string `json:"message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApiDeploymentList is a list of ApiArtifact resources
type ApiArtifactList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []ApiArtifact `json:"items"`
}
