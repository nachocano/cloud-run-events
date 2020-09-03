/*
Copyright 2019 Google LLC.

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

// Package auditlogs implements the CloudAuditLogsSource controller.
package auditlogs

import (
	"context"

	"knative.dev/pkg/injection"

	"k8s.io/client-go/tools/cache"
	serviceaccountinformers "knative.dev/pkg/client/injection/kube/informers/core/v1/serviceaccount"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"

	"github.com/google/knative-gcp/pkg/apis/configs/gcpauth"
	v1 "github.com/google/knative-gcp/pkg/apis/events/v1"
	"github.com/google/knative-gcp/pkg/pubsub/adapter/converters"
	"github.com/google/knative-gcp/pkg/reconciler"
	"github.com/google/knative-gcp/pkg/reconciler/identity"
	"github.com/google/knative-gcp/pkg/reconciler/identity/iam"
	"github.com/google/knative-gcp/pkg/reconciler/intevents"

	cloudauditlogssourceinformers "github.com/google/knative-gcp/pkg/client/injection/informers/events/v1/cloudauditlogssource"
	pullsubscriptioninformers "github.com/google/knative-gcp/pkg/client/injection/informers/intevents/v1/pullsubscription"
	topicinformers "github.com/google/knative-gcp/pkg/client/injection/informers/intevents/v1/topic"
	cloudauditlogssourcereconciler "github.com/google/knative-gcp/pkg/client/injection/reconciler/events/v1/cloudauditlogssource"
	glogadmin "github.com/google/knative-gcp/pkg/gclient/logging/logadmin"
	gpubsub "github.com/google/knative-gcp/pkg/gclient/pubsub"
)

const (
	// reconcilerName is the name of the reconciler
	reconcilerName = "CloudAuditLogsSource"

	// controllerAgentName is the string used by this controller to identify
	// itself when creating events.
	controllerAgentName = "cloud-run-events-cloudauditlogssource-controller"

	// receiveAdapterName is the string used as name for the receive adapter pod.
	receiveAdapterName = "cloudauditlogssource.events.cloud.google.com"
)

type Constructor injection.ControllerConstructor

// NewConstructor creates a constructor to make a CloudAuditLogsSource controller.
func NewConstructor(ipm iam.IAMPolicyManager, gcpas *gcpauth.StoreSingleton) Constructor {
	return func(ctx context.Context, cmw configmap.Watcher) *controller.Impl {
		return newController(ctx, cmw, ipm, gcpas.Store(ctx, cmw))
	}
}

func newController(
	ctx context.Context,
	cmw configmap.Watcher,
	ipm iam.IAMPolicyManager,
	gcpas *gcpauth.Store,
) *controller.Impl {
	pullsubscriptionInformer := pullsubscriptioninformers.Get(ctx)
	topicInformer := topicinformers.Get(ctx)
	cloudauditlogssourceInformer := cloudauditlogssourceinformers.Get(ctx)
	serviceAccountInformer := serviceaccountinformers.Get(ctx)

	r := &Reconciler{
		PubSubBase: intevents.NewPubSubBase(ctx,
			&intevents.PubSubBaseArgs{
				ControllerAgentName: controllerAgentName,
				ReceiveAdapterName:  receiveAdapterName,
				ReceiveAdapterType:  string(converters.CloudAuditLogs),
				ConfigWatcher:       cmw,
			}),
		Identity:               identity.NewIdentity(ctx, ipm, gcpas),
		auditLogsSourceLister:  cloudauditlogssourceInformer.Lister(),
		logadminClientProvider: glogadmin.NewClient,
		pubsubClientProvider:   gpubsub.NewClient,
		serviceAccountLister:   serviceAccountInformer.Lister(),
	}
	impl := cloudauditlogssourcereconciler.NewImpl(ctx, r)

	r.Logger.Info("Setting up event handlers")
	cloudauditlogssourceInformer.Informer().AddEventHandlerWithResyncPeriod(
		controller.HandleAll(impl.Enqueue), reconciler.DefaultResyncPeriod)

	auditLogsGK := v1.Kind("CloudAuditLogsSource")

	topicInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.FilterControllerGK(auditLogsGK),
		Handler:    controller.HandleAll(impl.EnqueueControllerOf),
	})

	pullsubscriptionInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.FilterControllerGK(auditLogsGK),
		Handler:    controller.HandleAll(impl.EnqueueControllerOf),
	})

	serviceAccountInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.FilterControllerGK(auditLogsGK),
		Handler:    controller.HandleAll(impl.EnqueueControllerOf),
	})

	return impl
}
