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

package k8s

import (
	"context"
	"go.uber.org/zap"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/cache"

	"github.com/google/knative-gcp/pkg/apis/pubsub/v1alpha1"
	psreconciler "github.com/google/knative-gcp/pkg/reconciler/pubsub/pullsubscription"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
)

// Reconciler implements controller.Reconciler for PullSubscription resources.
type Reconciler struct {
	*psreconciler.Base
}

// Check that our Reconciler implements controller.Reconciler
var _ controller.Reconciler = (*Reconciler)(nil)

// Reconcile compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the PullSubscription resource
// with the current status of the resource.
func (r *Reconciler) Reconcile(ctx context.Context, key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		logging.FromContext(ctx).Desugar().Error("Invalid resource key")
		return nil
	}
	// Get the PullSubscription resource with this namespace/name
	original, err := r.PullSubscriptionLister.PullSubscriptions(namespace).Get(name)
	if apierrs.IsNotFound(err) {
		// The resource may no longer exist, in which case we stop processing.
		logging.FromContext(ctx).Desugar().Error("PullSubscription in work queue no longer exists")
		return nil
	} else if err != nil {
		return err
	}

	// Don't modify the informers copy
	ps := original.DeepCopy()

	// Reconcile this copy of the PullSubscription and then write back any status
	// updates regardless of whether the reconciliation errored out.
	var reconcileErr = r.reconcile(ctx, ps)

	// If no error is returned, mark the observed generation.
	// This has to be done before updateStatus is called.
	if reconcileErr == nil {
		ps.Status.ObservedGeneration = ps.Generation
	}

	if equality.Semantic.DeepEqual(original.Finalizers, ps.Finalizers) {
		// If we didn't change finalizers then don't call updateFinalizers.

	} else if _, updated, fErr := r.Base.UpdateFinalizers(ctx, ps); fErr != nil {
		logging.FromContext(ctx).Desugar().Warn("Failed to update PullSubscription finalizers", zap.Error(fErr))
		r.Recorder.Eventf(ps, corev1.EventTypeWarning, "UpdateFailed",
			"Failed to update finalizers for PullSubscription %q: %v", ps.Name, fErr)
		return fErr
	} else if updated {
		// There was a difference and updateFinalizers said it updated and did not return an error.
		r.Recorder.Eventf(ps, corev1.EventTypeNormal, "Updated", "Updated PullSubscription %q finalizers", ps.Name)
	}

	if equality.Semantic.DeepEqual(original.Status, ps.Status) {
		// If we didn't change anything then don't call updateStatus.
		// This is important because the copy we loaded from the informer's
		// cache may be stale and we don't want to overwrite a prior update
		// to status with this stale state.

	} else if uErr := r.Base.UpdateStatus(ctx, original, ps); uErr != nil {
		logging.FromContext(ctx).Desugar().Warn("Failed to update ps status", zap.Error(uErr))
		r.Recorder.Eventf(ps, corev1.EventTypeWarning, "UpdateFailed",
			"Failed to update status for PullSubscription %q: %v", ps.Name, uErr)
		return uErr
	} else if reconcileErr == nil {
		// There was a difference and updateStatus did not return an error.
		r.Recorder.Eventf(ps, corev1.EventTypeNormal, "Updated", "Updated PullSubscription %q", ps.Name)
	}
	if reconcileErr != nil {
		r.Recorder.Event(ps, corev1.EventTypeWarning, "InternalError", reconcileErr.Error())
	}

	return reconcileErr
}

func (r *Reconciler) reconcile(ctx context.Context, ps *v1alpha1.PullSubscription) error {
	return r.Base.Reconcile(ctx, ps, r.reconcileDeployment)
}

func (r *Reconciler) reconcileDeployment(ctx context.Context, ra *appsv1.Deployment, src *v1alpha1.PullSubscription) error {
	existing, err := r.Base.GetOrCreateReceiveAdapter(ctx, ra, src)
	if err != nil {
		return err
	}
	if !equality.Semantic.DeepDerivative(ra.Spec, existing.Spec) {
		existing.Spec = ra.Spec
		_, err := r.KubeClientSet.AppsV1().Deployments(src.Namespace).Update(existing)
		if err != nil {
			logging.FromContext(ctx).Desugar().Error("Error updating Receive Adapter", zap.Error(err))
			return err
		}
	}
	return nil
}