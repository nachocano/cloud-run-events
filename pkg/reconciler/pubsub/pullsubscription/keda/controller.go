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

	gpubsub "github.com/google/knative-gcp/pkg/gclient/pubsub"
	"github.com/google/knative-gcp/pkg/reconciler"
	"github.com/google/knative-gcp/pkg/reconciler/pubsub"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
	"k8s.io/client-go/tools/cache"
	duckv1alpha1 "knative.dev/pkg/apis/duck/v1alpha1"
	pkgreconciler "knative.dev/pkg/reconciler"

	psreconciler "github.com/google/knative-gcp/pkg/reconciler/pubsub/pullsubscription"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/metrics"
	"knative.dev/pkg/resolver"
	tracingconfig "knative.dev/pkg/tracing/config"
	"knative.dev/pkg/tracker"

	pullsubscriptioninformers "github.com/google/knative-gcp/pkg/client/injection/informers/pubsub/v1alpha1/pullsubscription"
	deploymentinformer "knative.dev/pkg/client/injection/kube/informers/apps/v1/deployment"
)

const (
	// reconcilerName is the name of the reconciler
	reconcilerName = "KedaPullSubscriptions"

	// controllerAgentName is the string used by this controller to identify
	// itself when creating events.
	controllerAgentName = "cloud-run-events-pubsub-keda-pullsubscription-controller"

	finalizerName = controllerAgentName
)

type envConfig struct {
	// ReceiveAdapter is the receive adapters image. Required.
	ReceiveAdapter string `envconfig:"PUBSUB_RA_IMAGE" required:"true"`
}

// NewController initializes the controller and is called by the generated code
// Registers event handlers to enqueue events
func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {

	deploymentInformer := deploymentinformer.Get(ctx)
	pullSubscriptionInformer := pullsubscriptioninformers.Get(ctx)

	logger := logging.FromContext(ctx).Named(controllerAgentName).Desugar()

	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		logger.Fatal("Failed to process env var", zap.Error(err))
	}

	pubsubBase := &pubsub.PubSubBase{
		Base: reconciler.NewBase(ctx, controllerAgentName, cmw),
	}

	r := &Reconciler{
		Base: &psreconciler.Base{
			PubSubBase:             pubsubBase,
			DeploymentLister:       deploymentInformer.Lister(),
			PullSubscriptionLister: pullSubscriptionInformer.Lister(),
			ReceiveAdapterImage:    env.ReceiveAdapter,
			CreateClientFn:         gpubsub.NewClient,
			ControllerAgentName:    controllerAgentName,
			FinalizerName:          finalizerName,
		},
	}

	impl := controller.NewImpl(r, pubsubBase.Logger, reconcilerName)

	pubsubBase.Logger.Info("Setting up event handlers")
	onlyKedaScaler := pkgreconciler.AnnotationFilterFunc(duckv1alpha1.SourceScalerAnnotationKey, duckv1alpha1.KEDA, false)

	pullSubscriptionHandler := cache.FilteringResourceEventHandler{
		FilterFunc: onlyKedaScaler,
		Handler:    controller.HandleAll(impl.Enqueue),
	}
	pullSubscriptionInformer.Informer().AddEventHandlerWithResyncPeriod(pullSubscriptionHandler, reconciler.DefaultResyncPeriod)

	deploymentInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: onlyKedaScaler,
		Handler:    controller.HandleAll(impl.EnqueueControllerOf),
	})

	r.UriResolver = resolver.NewURIResolver(ctx, impl.EnqueueKey)
	r.ReconcileDataPlaneFn = r.ReconcileScaledObject
	// Tracker is used to notify us that a Keda ScaledObject has changed so that we can reconcile.
	// TODO not working. Need to call the OnChanged callback
	r.tracker = tracker.New(impl.EnqueueKey, controller.GetTrackerLease(ctx))

	cmw.Watch(logging.ConfigMapName(), r.UpdateFromLoggingConfigMap)
	cmw.Watch(metrics.ConfigMapName(), r.UpdateFromMetricsConfigMap)
	cmw.Watch(tracingconfig.ConfigName, r.UpdateFromTracingConfigMap)

	return impl
}
