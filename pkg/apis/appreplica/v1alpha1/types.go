package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)


// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AppReplica struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AppReplicaSpec   `json:"spec"`
	Status AppReplicaStatus `json:"status"`
}

// FooSpec is the spec for a Foo resource
type AppReplicaSpec struct {
	DeploymentName  string `json:"deploymentName"`
	DeploymentImage string `json:"deploymentImage"`
	Replicas        *int32 `json:"replicas"`
}

// FooStatus is the status for a Foo resource
type AppReplicaStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FooList is a list of Foo resources
type AppReplicaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []AppReplica `json:"items"`
}
