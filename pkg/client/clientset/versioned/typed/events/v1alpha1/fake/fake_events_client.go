/*
Copyright 2020 Google LLC

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
	v1alpha1 "github.com/google/knative-gcp/pkg/client/clientset/versioned/typed/events/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeEventsV1alpha1 struct {
	*testing.Fake
}

func (c *FakeEventsV1alpha1) CloudAuditLogsSources(namespace string) v1alpha1.CloudAuditLogsSourceInterface {
	return &FakeCloudAuditLogsSources{c, namespace}
}

func (c *FakeEventsV1alpha1) CloudPubSubSources(namespace string) v1alpha1.CloudPubSubSourceInterface {
	return &FakeCloudPubSubSources{c, namespace}
}

func (c *FakeEventsV1alpha1) CloudSchedulerSources(namespace string) v1alpha1.CloudSchedulerSourceInterface {
	return &FakeCloudSchedulerSources{c, namespace}
}

func (c *FakeEventsV1alpha1) CloudStorageSources(namespace string) v1alpha1.CloudStorageSourceInterface {
	return &FakeCloudStorageSources{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeEventsV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
