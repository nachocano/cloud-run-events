/*
Copyright 2019 The Knative Authors

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

package main

import (
	"context"
	"fmt"

	eventsv1alpha1 "github.com/google/knative-gcp/pkg/apis/events/v1alpha1"
	messagingv1alpha1 "github.com/google/knative-gcp/pkg/apis/messaging/v1alpha1"
	pubsubv1alpha1 "github.com/google/knative-gcp/pkg/apis/pubsub/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/injection/sharedmain"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/signals"
	"knative.dev/pkg/system"
	"knative.dev/pkg/webhook"
	"knative.dev/pkg/webhook/certificates"
	"knative.dev/pkg/webhook/resourcesemantics"
)

const (
	component = "webhook"
)

func NewResourceAdmissionController(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
	logger := logging.FromContext(ctx)

	// Decorate contexts with the current state of the config.
	// store := defaultconfig.NewStore(logging.FromContext(ctx).Named("config-store"))
	// store.WatchConfigs(cmw)
	ctxFunc := func(ctx context.Context) context.Context {
		// return v1.WithUpgradeViaDefaulting(store.ToContext(ctx))
		return ctx
	}

	return resourcesemantics.NewAdmissionController(ctx,

		// Name of the resource webhook.
		fmt.Sprintf("webhook.%s.events.cloud.google.com", system.Namespace()),

		// The path on which to serve the webhook.
		"/",

		// The resources to validate and default.
		map[schema.GroupVersionKind]resourcesemantics.GenericCRD{
			// For group messaging.cloud.google.com
			messagingv1alpha1.SchemeGroupVersion.WithKind("Channel"):   &messagingv1alpha1.Channel{},
			messagingv1alpha1.SchemeGroupVersion.WithKind("Decorator"): &messagingv1alpha1.Decorator{},

			// For group events.cloud.google.com
			eventsv1alpha1.SchemeGroupVersion.WithKind("Storage"):   &eventsv1alpha1.Storage{},
			eventsv1alpha1.SchemeGroupVersion.WithKind("Scheduler"): &eventsv1alpha1.Scheduler{},
			eventsv1alpha1.SchemeGroupVersion.WithKind("PubSub"):    &eventsv1alpha1.PubSub{},

			// For group pubsub.cloud.google.com
			pubsubv1alpha1.SchemeGroupVersion.WithKind("PullSubscription"): &pubsubv1alpha1.PullSubscription{},
			pubsubv1alpha1.SchemeGroupVersion.WithKind("Topic"):            &pubsubv1alpha1.Topic{},
		},

		// A function that infuses the context passed to Validate/SetDefaults with custom metadata.
		ctxFunc,

		// Whether to disallow unknown fields.
		true,
	)
}

func main() {
	// Set up a signal context with our webhook options
	ctx := webhook.WithOptions(signals.NewContext(), webhook.Options{
		ServiceName: component,
		Port:        8443,
		SecretName:  "webhook-certs",
	})

	sharedmain.MainWithContext(ctx, component,
		certificates.NewController,
		NewResourceAdmissionController,
	)
}
