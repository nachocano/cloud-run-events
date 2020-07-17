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
package deployment

import (
	"context"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/clock"
	"k8s.io/client-go/tools/cache"
	"knative.dev/pkg/client/injection/kube/informers/apps/v1/deployment"
	"knative.dev/pkg/client/injection/kube/informers/core/v1/secret"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"os"

	"github.com/google/knative-gcp/pkg/apis/duck"
	"github.com/google/knative-gcp/pkg/reconciler"
	pkgreconciler "knative.dev/pkg/reconciler"
)

const (
	// ReconcilerName is the name of the reconciler
	ReconcilerName = "Deployment"

	// controllerAgentName is the string used by this controller to identify
	// itself when creating events.
	controllerAgentName = "cloud-run-events-deployment-controller"

	namespace      = "cloud-run-events"
	secretName     = duck.DefaultSecretName
	deploymentName = "controller"
	envKey         = "GOOGLE_APPLICATION_CREDENTIALS"
)

// NewController initializes the controller and is called by the generated code
// Registers event handlers to enqueue events.
// When the secret `google-cloud-key` of namespace `cloud-run-events` gets updated, we will enqueue the deployment `controller` of namespace `cloud-run-events`.
func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {

	deploymentInformer := deployment.Get(ctx)
	secretInformer := secret.Get(ctx)

	key := types.NamespacedName{Namespace: namespace, Name: deploymentName}

	r := &Reconciler{
		Base: reconciler.NewBase(ctx, controllerAgentName, cmw),
		LeaderAwareFuncs: pkgreconciler.LeaderAwareFuncs{
			// Have this reconciler enqueue our singleton whenever it becomes leader.
			PromoteFunc: func(bkt pkgreconciler.Bucket, enq func(pkgreconciler.Bucket, types.NamespacedName)) error {
				enq(bkt, key)
				return nil
			},
		},
		deploymentLister: deploymentInformer.Lister(),
		clock:            clock.RealClock{},
		key:              key,
	}

	impl := controller.NewImpl(r, r.Logger, ReconcilerName)

	r.Logger.Info("Setting up event handlers")

	secretInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.FilterWithNameAndNamespace(namespace, secretName),
		Handler:    handler(impl.Enqueue),
	})

	return impl
}

func handler(h func(interface{})) cache.ResourceEventHandler {
	return cache.ResourceEventHandlerFuncs{
		// For AddFunc, only enqueue deployment key when envKey is not set.
		// In such case, the controller pod hasn't restarted before.
		// This helps to avoid infinite loop for restarting controller pod.
		AddFunc: func(obj interface{}) {
			if _, ok := os.LookupEnv(envKey); !ok {
				h(obj)
			}
		},
		UpdateFunc: controller.PassNew(h),
		// If secret is deleted, the controller pod will restart, in order to unset the envKey.
		// This is needed when changing authentication configuration from k8s Secret to Workload Identity.
		DeleteFunc: h,
	}
}
