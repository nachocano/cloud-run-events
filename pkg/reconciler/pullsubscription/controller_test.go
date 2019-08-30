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
	"os"
	"testing"

	"knative.dev/pkg/logging"
	"knative.dev/pkg/metrics"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"knative.dev/pkg/configmap"
	logtesting "knative.dev/pkg/logging/testing"
	. "knative.dev/pkg/reconciler/testing"
	"knative.dev/pkg/system"

	_ "knative.dev/pkg/metrics/testing"

	// Fake injection informers
	_ "knative.dev/pkg/injection/informers/kubeinformers/appsv1/deployment/fake"
	_ "knative.dev/pkg/injection/informers/kubeinformers/batchv1/job/fake"

	_ "github.com/google/knative-gcp/pkg/client/injection/informers/pubsub/v1alpha1/pullsubscription/fake"
)

func TestNew(t *testing.T) {
	defer logtesting.ClearAll()
	ctx, _ := SetupFakeContext(t)

	_ = os.Setenv("PUBSUB_RA_IMAGE", "PUBSUB_RA_IMAGE")
	_ = os.Setenv("PUBSUB_SUB_IMAGE", "PUBSUB_SUB_IMAGE")

	c := NewController(ctx, configmap.NewStaticWatcher(
		&corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      logging.ConfigMapName(),
				Namespace: system.Namespace(),
			},
			Data: map[string]string{},
		},
		&corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      metrics.ConfigMapName(),
				Namespace: system.Namespace(),
			},
			Data: map[string]string{},
		},
	))

	if c == nil {
		t.Fatal("Expected NewController to return a non-nil value")
	}
}
