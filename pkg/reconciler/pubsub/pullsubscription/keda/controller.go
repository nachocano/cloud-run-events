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

package pullsubscription

import (
	"context"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
)

const (
	// reconcilerName is the name of the reconciler
	reconcilerName = "KedaPullSubscriptions"

	// controllerAgentName is the string used by this controller to identify
	// itself when creating events.
	controllerAgentName = "cloud-run-events-pubsub-keda-pullsubscription-controller"
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
	return nil
	//deploymentInformer := deploymentinformer.Get(ctx)
	//pullSubscriptionInformer := pullsubscriptioninformers.Get(ctx)
	//
	//logger := logging.FromContext(ctx).Named(controllerAgentName).Desugar()
	//
	//var env envConfig
	//if err := envconfig.Process("", &env); err != nil {
	//	logger.Fatal("Failed to process env var", zap.Error(err))
	//}
	//
	//pubsubBase := &pubsub.PubSubBase{
	//	Base: reconciler.NewBase(ctx, controllerAgentName, cmw),
	//}
	//
	//r := &Reconciler{
	//	PubSubBase:             pubsubBase,
	//	deploymentLister:       deploymentInformer.Lister(),
	//	pullSubscriptionLister: pullSubscriptionInformer.Lister(),
	//	receiveAdapterImage:    env.ReceiveAdapter,
	//	createClientFn:         gpubsub.NewClient,
	//}
	//
	//impl := controller.NewImpl(r, pubsubBase.Logger, reconcilerName)
	//
	//onlyHpaClass := reconciler.AnnotationFilterFunc(duckv1alpha1.ClassAnnotationKey, autoscaling.HPA, false)
	//
	//pubsubBase.Logger.Info("Setting up event handlers")
	//pullSubscriptionInformer.Informer().AddEventHandlerWithResyncPeriod(controller.HandleAll(impl.Enqueue), reconciler.DefaultResyncPeriod)
	//
	//deploymentInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
	//	FilterFunc: controller.Filter(v1alpha1.SchemeGroupVersion.WithKind("PullSubscription")),
	//	Handler:    controller.HandleAll(impl.EnqueueControllerOf),
	//})
	//
	//r.uriResolver = resolver.NewURIResolver(ctx, impl.EnqueueKey)
	//
	//cmw.Watch(logging.ConfigMapName(), r.UpdateFromLoggingConfigMap)
	//cmw.Watch(metrics.ConfigMapName(), r.UpdateFromMetricsConfigMap)
	//cmw.Watch(tracingconfig.ConfigName, r.UpdateFromTracingConfigMap)
	//
	//// TODO discovery, if keda not install fail.
	//// TODO watch ScaledObjects.
	//// TODO upstream common stuff to pkg.
	//
	//return impl
}
