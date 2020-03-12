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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	time "time"

	securityv1alpha1 "github.com/google/knative-gcp/pkg/apis/security/v1alpha1"
	versioned "github.com/google/knative-gcp/pkg/client/clientset/versioned"
	internalinterfaces "github.com/google/knative-gcp/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/google/knative-gcp/pkg/client/listers/security/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// EventPolicyBindingInformer provides access to a shared informer and lister for
// EventPolicyBindings.
type EventPolicyBindingInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.EventPolicyBindingLister
}

type eventPolicyBindingInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewEventPolicyBindingInformer constructs a new informer for EventPolicyBinding type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewEventPolicyBindingInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredEventPolicyBindingInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredEventPolicyBindingInformer constructs a new informer for EventPolicyBinding type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredEventPolicyBindingInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.SecurityV1alpha1().EventPolicyBindings(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.SecurityV1alpha1().EventPolicyBindings(namespace).Watch(options)
			},
		},
		&securityv1alpha1.EventPolicyBinding{},
		resyncPeriod,
		indexers,
	)
}

func (f *eventPolicyBindingInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredEventPolicyBindingInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *eventPolicyBindingInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&securityv1alpha1.EventPolicyBinding{}, f.defaultInformer)
}

func (f *eventPolicyBindingInformer) Lister() v1alpha1.EventPolicyBindingLister {
	return v1alpha1.NewEventPolicyBindingLister(f.Informer().GetIndexer())
}
