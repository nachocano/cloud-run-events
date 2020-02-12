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

package keda

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/google/knative-gcp/pkg/apis/pubsub/v1alpha1"
	psreconciler "github.com/google/knative-gcp/pkg/reconciler/pubsub/pullsubscription"
	"github.com/google/knative-gcp/pkg/reconciler/pubsub/pullsubscription/keda/resources"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
)

// Reconciler implements controller.Reconciler for PullSubscription resources.
type Reconciler struct {
	*psreconciler.Base
}

// Check that our Reconciler implements controller.Reconciler
var _ controller.Reconciler = (*Reconciler)(nil)

func (r *Reconciler) Reconcile(ctx context.Context, key string) error {
	return r.Base.Reconcile(ctx, key)
}

func (r *Reconciler) ReconcileScaledObject(ctx context.Context, ra *appsv1.Deployment, src *v1alpha1.PullSubscription) error {
	// TODO discovery
	// TODO reconcileDeployment but do not consider replicas.
	// TODO reconcileScaledObject... if we created a new deployment
	// TODO upstream to pkg
	// TODO tracker

	gvr, _ := meta.UnsafeGuessKindToResource(resources.ScaledObjectGVK)
	scaledObjectResourceInterface := r.DynamicClientSet.Resource(gvr).Namespace(src.Namespace)
	if scaledObjectResourceInterface == nil {
		return fmt.Errorf("unable to create dynamic client for ScaledObject")
	}

	so := resources.MakeScaledObject(ctx, ra, src)
	_, err := scaledObjectResourceInterface.Get(so.GetName(), metav1.GetOptions{})
	if err != nil {
		if apierrs.IsNotFound(err) {
			// If there is no scaledObject, we need to create the deployment as well,
			// so that the ScaledObject can set it as its scaleTargetRef.
			_, err := r.Base.GetOrCreateReceiveAdapter(ctx, ra, src)
			if err != nil {
				logging.FromContext(ctx).Desugar().Error("Failed to get or create Receive Adapter in order to create the ScaledObject", zap.Error(err))
				return err
			}
			// Note that the ScaledObject is the one in charge of reconciling the Deployment, e.g., by changing the number
			// of replicas.
			_, err = scaledObjectResourceInterface.Create(so, metav1.CreateOptions{})
			if err != nil {
				logging.FromContext(ctx).Desugar().Error("Failed to create ScaledObject", zap.Any("so", so), zap.Error(err))
				return err
			}
		} else {
			logging.FromContext(ctx).Desugar().Error("Failed to get ScaledObject", zap.Any("so", so), zap.Error(err))
		}
	}
	return nil
}
