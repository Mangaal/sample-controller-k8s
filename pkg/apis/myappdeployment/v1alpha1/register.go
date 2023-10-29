package v1alpha1

import (
	myappdeploymentv1alpha1 "sample-controller-k8s/pkg/apis/myappdeployment"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var SchemeGroupVersion = schema.GroupVersion{Group: myappdeploymentv1alpha1.GroupName, Version: "v1aplha1"}

func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

func addKnownTypes(sheme *runtime.Scheme) error {
	sheme.AddKnownTypes(SchemeGroupVersion,
		&MyAppDeployment{},
		&MyAppDeploymentList{},
	)
	metav1.AddToGroupVersion(sheme, SchemeGroupVersion)

	return nil

}
