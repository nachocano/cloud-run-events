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

// Code generated by lister-gen. DO NOT EDIT.

package v1beta1

import (
	v1beta1 "istio.io/client-go/pkg/apis/security/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// PeerAuthenticationLister helps list PeerAuthentications.
type PeerAuthenticationLister interface {
	// List lists all PeerAuthentications in the indexer.
	List(selector labels.Selector) (ret []*v1beta1.PeerAuthentication, err error)
	// PeerAuthentications returns an object that can list and get PeerAuthentications.
	PeerAuthentications(namespace string) PeerAuthenticationNamespaceLister
	PeerAuthenticationListerExpansion
}

// peerAuthenticationLister implements the PeerAuthenticationLister interface.
type peerAuthenticationLister struct {
	indexer cache.Indexer
}

// NewPeerAuthenticationLister returns a new PeerAuthenticationLister.
func NewPeerAuthenticationLister(indexer cache.Indexer) PeerAuthenticationLister {
	return &peerAuthenticationLister{indexer: indexer}
}

// List lists all PeerAuthentications in the indexer.
func (s *peerAuthenticationLister) List(selector labels.Selector) (ret []*v1beta1.PeerAuthentication, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.PeerAuthentication))
	})
	return ret, err
}

// PeerAuthentications returns an object that can list and get PeerAuthentications.
func (s *peerAuthenticationLister) PeerAuthentications(namespace string) PeerAuthenticationNamespaceLister {
	return peerAuthenticationNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// PeerAuthenticationNamespaceLister helps list and get PeerAuthentications.
type PeerAuthenticationNamespaceLister interface {
	// List lists all PeerAuthentications in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1beta1.PeerAuthentication, err error)
	// Get retrieves the PeerAuthentication from the indexer for a given namespace and name.
	Get(name string) (*v1beta1.PeerAuthentication, error)
	PeerAuthenticationNamespaceListerExpansion
}

// peerAuthenticationNamespaceLister implements the PeerAuthenticationNamespaceLister
// interface.
type peerAuthenticationNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all PeerAuthentications in the indexer for a given namespace.
func (s peerAuthenticationNamespaceLister) List(selector labels.Selector) (ret []*v1beta1.PeerAuthentication, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1beta1.PeerAuthentication))
	})
	return ret, err
}

// Get retrieves the PeerAuthentication from the indexer for a given namespace and name.
func (s peerAuthenticationNamespaceLister) Get(name string) (*v1beta1.PeerAuthentication, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1beta1.Resource("peerauthentication"), name)
	}
	return obj.(*v1beta1.PeerAuthentication), nil
}
