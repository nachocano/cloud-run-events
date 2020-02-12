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

package resources

import (
	"context"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/google/knative-gcp/pkg/apis/pubsub/v1alpha1"
	"k8s.io/api/apps/v1"
)

var (
	ScaledObjectGVK = schema.GroupVersionKind{
		Group:   "keda.k8s.io",
		Version: "v1alpha1",
		Kind:    "ScaledObject",
	}
)

func MakeScaledObject(ctx context.Context, ra *v1.Deployment, ps *v1alpha1.PullSubscription) *unstructured.Unstructured {
	metadata := ps.Spec.Scaler.Metadata
	metadata["subscriptionName"] = ps.Status.SubscriptionID
	metadata["credentials"] = "GOOGLE_APPLICATION_CREDENTIALS_JSON"

	so := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "keda.k8s.io/v1alpha1",
			"kind":       "ScaledObject",
			"metadata": map[string]interface{}{
				"namespace": ps.Namespace,
				// TODO autogenerate the name using the UID
				"name": ps.Name,
				"labels": map[string]interface{}{
					"deploymentName": ra.Name,
				},
				"ownerReferences": []map[string]interface{}{{
					"apiVersion":         ps.GetGroupVersion().String(),
					"kind":               ps.GetGroupVersionKind().Kind,
					"blockOwnerDeletion": true,
					"controller":         true,
					"name":               ps.Name,
					"uid":                ps.UID,
				}},
			},
			"spec": map[string]interface{}{
				"scaleTargetRef": map[string]interface{}{
					"deploymentName": ra.Name,
				},
				"minReplicaCount": *ps.Spec.Scaler.MinScale,
				"maxReplicaCount": *ps.Spec.Scaler.MaxScale,
				"triggers": []map[string]interface{}{{
					"type":     ps.Spec.Scaler.Type,
					"metadata": metadata,
				}},
			},
		},
	}
	return so
}
