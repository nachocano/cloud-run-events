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

package v1alpha1

import (
	"context"
	"time"

	duckv1alpha1 "github.com/google/knative-gcp/pkg/apis/duck/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"knative.dev/pkg/apis"
	pkgduckv1alpha1 "knative.dev/pkg/apis/duck/v1alpha1"
	"knative.dev/pkg/ptr"
)

const (
	defaultRetentionDuration = 7 * 24 * time.Hour
	defaultAckDeadline       = 30 * time.Second

	defaultScalerType       = "gcp-pubsub"
	defaultSubscriptionSize = "5"
	subscriptionSize        = "subscriptionSize"
)

func (s *PullSubscription) SetDefaults(ctx context.Context) {
	withParent := apis.WithinParent(ctx, s.ObjectMeta)
	s.Spec.SetDefaults(withParent)
}

func (ss *PullSubscriptionSpec) SetDefaults(ctx context.Context) {
	if ss.AckDeadline == nil {
		ackDeadline := defaultAckDeadline
		ss.AckDeadline = ptr.String(ackDeadline.String())
	}

	if ss.RetentionDuration == nil {
		retentionDuration := defaultRetentionDuration
		ss.RetentionDuration = ptr.String(retentionDuration.String())
	}

	if ss.Secret == nil || equality.Semantic.DeepEqual(ss.Secret, &corev1.SecretKeySelector{}) {
		ss.Secret = duckv1alpha1.DefaultGoogleCloudSecretSelector()
	}

	switch ss.Mode {
	case ModeCloudEventsBinary, ModeCloudEventsStructured, ModePushCompatible:
		// Valid Mode.
	default:
		// Default is CloudEvents Binary Mode.
		ss.Mode = ModeCloudEventsBinary
	}

	if ss.Scaler != nil && !equality.Semantic.DeepEqual(ss.Scaler, &pkgduckv1alpha1.KedaScalerSpec{}) {
		// Set up common defaults for Keda scalers.
		ss.Scaler.SetDefault(ctx)

		// Set up our own defaults.
		if ss.Scaler.Type == "" {
			ss.Scaler.Type = defaultScalerType
		}

		if ss.Scaler.Metadata == nil {
			ss.Scaler.Metadata = make(map[string]string)
		}
		if _, ok := ss.Scaler.Metadata[subscriptionSize]; !ok {
			ss.Scaler.Metadata[subscriptionSize] = defaultSubscriptionSize
		}

		parentMeta := apis.ParentMeta(ctx)
		if _, ok := parentMeta.Annotations[pkgduckv1alpha1.SourceScalerAnnotationKey]; !ok {
			parentMeta.Annotations[pkgduckv1alpha1.SourceScalerAnnotationKey] = pkgduckv1alpha1.KEDA
		}
	}
}
