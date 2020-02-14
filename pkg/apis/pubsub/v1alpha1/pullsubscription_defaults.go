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
	"knative.dev/pkg/ptr"
)

const (
	defaultRetentionDuration = 7 * 24 * time.Hour
	defaultAckDeadline       = 30 * time.Second

	SourceAutoScalerAnnotationKey = "sources.knative.dev/autoscaler"
	SourceMinScaleAnnotationKey   = "sources.knative.dev/minScale"
	SourceMaxScaleAnnotationKey   = "sources.knative.dev/maxScale"

	KedaAutoScaler                    = "keda"
	KedaScalerAnnotationKey           = "keda.knative.dev/scaler"
	KedaPollingIntervalAnnotationKey  = "keda.knative.dev/pollingInterval"
	KedaCooldownPeriodAnnotationKey   = "keda.knative.dev/cooldownPeriod"
	KedaSubscriptionSizeAnnotationKey = "keda.knative.dev/subscriptionSize"

	defaultMinScale = "0"
	defaultMaxScale = "1"

	defaultKedaScaler           = "gcp-pubsub"
	defaultKedaSubscriptionSize = "5"
	defaultKedaPollingInterval  = "30"
	defaultKedaCooldownPeriod   = "120"
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

	// TODO move to some method
	parentMeta := apis.ParentMeta(ctx)
	if _, ok := parentMeta.Annotations[SourceAutoScalerAnnotationKey]; !ok {
		parentMeta.Annotations[SourceAutoScalerAnnotationKey] = KedaAutoScaler
	}

	if parentMeta.Annotations[SourceAutoScalerAnnotationKey] == KedaAutoScaler {
		if _, ok := parentMeta.Annotations[SourceMinScaleAnnotationKey]; !ok {
			parentMeta.Annotations[SourceMinScaleAnnotationKey] = defaultMinScale
		}
		if _, ok := parentMeta.Annotations[SourceMaxScaleAnnotationKey]; !ok {
			parentMeta.Annotations[SourceMaxScaleAnnotationKey] = defaultMaxScale
		}
		if _, ok := parentMeta.Annotations[KedaPollingIntervalAnnotationKey]; !ok {
			parentMeta.Annotations[KedaPollingIntervalAnnotationKey] = defaultKedaPollingInterval
		}
		if _, ok := parentMeta.Annotations[KedaCooldownPeriodAnnotationKey]; !ok {
			parentMeta.Annotations[KedaCooldownPeriodAnnotationKey] = defaultKedaCooldownPeriod
		}
		if _, ok := parentMeta.Annotations[KedaSubscriptionSizeAnnotationKey]; !ok {
			parentMeta.Annotations[KedaSubscriptionSizeAnnotationKey] = defaultKedaSubscriptionSize
		}
		if _, ok := parentMeta.Annotations[KedaScalerAnnotationKey]; !ok {
			parentMeta.Annotations[KedaScalerAnnotationKey] = defaultKedaScaler
		}
	}
}
