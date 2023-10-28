module sample-controller-k8s

go 1.13

require (
	golang.org/x/time v0.3.0
	k8s.io/api v0.28.2
	k8s.io/apimachinery v0.28.2
	k8s.io/client-go v0.28.2
	sigs.k8s.io/controller-runtime v0.16.2
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3
)

require golang.org/x/net v0.17.0 // indirect
