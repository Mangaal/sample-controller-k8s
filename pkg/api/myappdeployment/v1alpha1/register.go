package v1alpha1

import (
	myappdeployment "sample-controller-k8s/pkg/api/MyAppDeployment"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var SchemeGroupVesrion = schema.GroupVersion{Group: myappdeployment.GroupName, Version: "v1aplha1"}

func Kind(kind string) schema.GroupKind {
	return SchemeGroupVesrion.WithKind(kind).GroupKind()
}

func Resource(resource string) schema.GroupResource {
	return SchemeGroupVesrion.WithResource(resource).GroupResource()
}

var (
	SchemeBuiler = runtime.NewSchemeBuilder()
)

func addKnownTypes(sheme *runtime.Scheme) error {
	sheme.AddKnownTypes(SchemeGroupVesrion,
		&MyAppDeployment{},
		&MyAppDeploymentList{},
	)
	metav1.AddToGroupVersion(sheme, SchemeGroupVesrion)

	return nil

}
