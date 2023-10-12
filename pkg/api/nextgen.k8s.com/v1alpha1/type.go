package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type DeploymentResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeploymentResourceSpec   `json:"spec"`
	Status DeploymentResourceStatus `json:"status"`
}

type DeploymentResourceSpec struct {
	DeploymentName  string `json:"deploymentName"`
	DeploymentImage string `json:"deploymentImage"`
	Replicas        *int32 `json:"replicas"`
}

type DeploymentResourceStatus struct {
	AvailableReplica int32 `json:"availableReplica"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type DeploymentResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []DeploymentResource `json:"items"`
}
