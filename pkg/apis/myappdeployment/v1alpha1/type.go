package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type MyAppDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MyAppDeploymentSpec   `json:"spec"`
	Status MyAppDeploymentStatus `json:"status"`
}

type MyAppDeploymentSpec struct {
	DeploymentName  string `json:"deploymentName"`
	DeploymentImage string `json:"deploymentImage"`
	Replicas        *int32 `json:"replicas"`
}

type MyAppDeploymentStatus struct {
	AvailableReplica int32 `json:"availableReplica"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type MyAppDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []MyAppDeployment `json:"items"`
}
