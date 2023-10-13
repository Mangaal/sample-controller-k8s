/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	json "encoding/json"
	"fmt"
	v1alpha1 "sample-controller-k8s/pkg/api/myappdeployment/v1alpha1"
	myappdeploymentv1alpha1 "sample-controller-k8s/pkg/clinet/applyconfiguration/myappdeployment/v1alpha1"
	scheme "sample-controller-k8s/pkg/clinet/clientset/versioned/scheme"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// MyAppDeploymentsGetter has a method to return a MyAppDeploymentInterface.
// A group's client should implement this interface.
type MyAppDeploymentsGetter interface {
	MyAppDeployments(namespace string) MyAppDeploymentInterface
}

// MyAppDeploymentInterface has methods to work with MyAppDeployment resources.
type MyAppDeploymentInterface interface {
	Create(ctx context.Context, myAppDeployment *v1alpha1.MyAppDeployment, opts v1.CreateOptions) (*v1alpha1.MyAppDeployment, error)
	Update(ctx context.Context, myAppDeployment *v1alpha1.MyAppDeployment, opts v1.UpdateOptions) (*v1alpha1.MyAppDeployment, error)
	UpdateStatus(ctx context.Context, myAppDeployment *v1alpha1.MyAppDeployment, opts v1.UpdateOptions) (*v1alpha1.MyAppDeployment, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.MyAppDeployment, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.MyAppDeploymentList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.MyAppDeployment, err error)
	Apply(ctx context.Context, myAppDeployment *myappdeploymentv1alpha1.MyAppDeploymentApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.MyAppDeployment, err error)
	ApplyStatus(ctx context.Context, myAppDeployment *myappdeploymentv1alpha1.MyAppDeploymentApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.MyAppDeployment, err error)
	MyAppDeploymentExpansion
}

// myAppDeployments implements MyAppDeploymentInterface
type myAppDeployments struct {
	client rest.Interface
	ns     string
}

// newMyAppDeployments returns a MyAppDeployments
func newMyAppDeployments(c *NextgenV1alpha1Client, namespace string) *myAppDeployments {
	return &myAppDeployments{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the myAppDeployment, and returns the corresponding myAppDeployment object, and an error if there is any.
func (c *myAppDeployments) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.MyAppDeployment, err error) {
	result = &v1alpha1.MyAppDeployment{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("myappdeployments").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of MyAppDeployments that match those selectors.
func (c *myAppDeployments) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.MyAppDeploymentList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.MyAppDeploymentList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("myappdeployments").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested myAppDeployments.
func (c *myAppDeployments) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("myappdeployments").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a myAppDeployment and creates it.  Returns the server's representation of the myAppDeployment, and an error, if there is any.
func (c *myAppDeployments) Create(ctx context.Context, myAppDeployment *v1alpha1.MyAppDeployment, opts v1.CreateOptions) (result *v1alpha1.MyAppDeployment, err error) {
	result = &v1alpha1.MyAppDeployment{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("myappdeployments").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(myAppDeployment).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a myAppDeployment and updates it. Returns the server's representation of the myAppDeployment, and an error, if there is any.
func (c *myAppDeployments) Update(ctx context.Context, myAppDeployment *v1alpha1.MyAppDeployment, opts v1.UpdateOptions) (result *v1alpha1.MyAppDeployment, err error) {
	result = &v1alpha1.MyAppDeployment{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("myappdeployments").
		Name(myAppDeployment.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(myAppDeployment).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *myAppDeployments) UpdateStatus(ctx context.Context, myAppDeployment *v1alpha1.MyAppDeployment, opts v1.UpdateOptions) (result *v1alpha1.MyAppDeployment, err error) {
	result = &v1alpha1.MyAppDeployment{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("myappdeployments").
		Name(myAppDeployment.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(myAppDeployment).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the myAppDeployment and deletes it. Returns an error if one occurs.
func (c *myAppDeployments) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("myappdeployments").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *myAppDeployments) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("myappdeployments").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched myAppDeployment.
func (c *myAppDeployments) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.MyAppDeployment, err error) {
	result = &v1alpha1.MyAppDeployment{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("myappdeployments").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied myAppDeployment.
func (c *myAppDeployments) Apply(ctx context.Context, myAppDeployment *myappdeploymentv1alpha1.MyAppDeploymentApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.MyAppDeployment, err error) {
	if myAppDeployment == nil {
		return nil, fmt.Errorf("myAppDeployment provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(myAppDeployment)
	if err != nil {
		return nil, err
	}
	name := myAppDeployment.Name
	if name == nil {
		return nil, fmt.Errorf("myAppDeployment.Name must be provided to Apply")
	}
	result = &v1alpha1.MyAppDeployment{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("myappdeployments").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *myAppDeployments) ApplyStatus(ctx context.Context, myAppDeployment *myappdeploymentv1alpha1.MyAppDeploymentApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.MyAppDeployment, err error) {
	if myAppDeployment == nil {
		return nil, fmt.Errorf("myAppDeployment provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(myAppDeployment)
	if err != nil {
		return nil, err
	}

	name := myAppDeployment.Name
	if name == nil {
		return nil, fmt.Errorf("myAppDeployment.Name must be provided to Apply")
	}

	result = &v1alpha1.MyAppDeployment{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("myappdeployments").
		Name(*name).
		SubResource("status").
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}