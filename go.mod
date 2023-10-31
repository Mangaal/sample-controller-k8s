module sample-controller-k8s

go 1.16

require (
	k8s.io/apimachinery v0.28.3
	k8s.io/client-go v0.28.3
	k8s.io/klog/v2 v2.100.1
	sigs.k8s.io/controller-runtime v0.16.3
)

replace (
	k8s.io/api => k8s.io/api v0.24.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.24.0
	k8s.io/client-go => k8s.io/client-go v0.24.0
	k8s.io/code-generator => k8s.io/code-generator v0.24.0
)
