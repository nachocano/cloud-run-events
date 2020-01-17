/*
Copyright 2019 Google LLC

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

package fake

import (
	v1alpha1 "github.com/google/knative-gcp/pkg/apis/events/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeCloudStorageSources implements CloudStorageSourceInterface
type FakeCloudStorageSources struct {
	Fake *FakeEventsV1alpha1
	ns   string
}

var storagesResource = schema.GroupVersionResource{Group: "events.cloud.google.com", Version: "v1alpha1", Resource: "storages"}

var storagesKind = schema.GroupVersionKind{Group: "events.cloud.google.com", Version: "v1alpha1", Kind: "CloudStorageSource"}

// Get takes name of the storage, and returns the corresponding storage object, and an error if there is any.
func (c *FakeCloudStorageSources) Get(name string, options v1.GetOptions) (result *v1alpha1.CloudStorageSource, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(storagesResource, c.ns, name), &v1alpha1.CloudStorageSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.CloudStorageSource), err
}

// List takes label and field selectors, and returns the list of CloudStorageSources that match those selectors.
func (c *FakeCloudStorageSources) List(opts v1.ListOptions) (result *v1alpha1.CloudStorageSourceList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(storagesResource, storagesKind, c.ns, opts), &v1alpha1.CloudStorageSourceList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.CloudStorageSourceList{ListMeta: obj.(*v1alpha1.CloudStorageSourceList).ListMeta}
	for _, item := range obj.(*v1alpha1.CloudStorageSourceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested storages.
func (c *FakeCloudStorageSources) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(storagesResource, c.ns, opts))

}

// Create takes the representation of a storage and creates it.  Returns the server's representation of the storage, and an error, if there is any.
func (c *FakeCloudStorageSources) Create(storage *v1alpha1.CloudStorageSource) (result *v1alpha1.CloudStorageSource, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(storagesResource, c.ns, storage), &v1alpha1.CloudStorageSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.CloudStorageSource), err
}

// Update takes the representation of a storage and updates it. Returns the server's representation of the storage, and an error, if there is any.
func (c *FakeCloudStorageSources) Update(storage *v1alpha1.CloudStorageSource) (result *v1alpha1.CloudStorageSource, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(storagesResource, c.ns, storage), &v1alpha1.CloudStorageSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.CloudStorageSource), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeCloudStorageSources) UpdateStatus(storage *v1alpha1.CloudStorageSource) (*v1alpha1.CloudStorageSource, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(storagesResource, "status", c.ns, storage), &v1alpha1.CloudStorageSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.CloudStorageSource), err
}

// Delete takes name of the storage and deletes it. Returns an error if one occurs.
func (c *FakeCloudStorageSources) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(storagesResource, c.ns, name), &v1alpha1.CloudStorageSource{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeCloudStorageSources) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(storagesResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.CloudStorageSourceList{})
	return err
}

// Patch applies the patch and returns the patched storage.
func (c *FakeCloudStorageSources) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.CloudStorageSource, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(storagesResource, c.ns, name, pt, data, subresources...), &v1alpha1.CloudStorageSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.CloudStorageSource), err
}
