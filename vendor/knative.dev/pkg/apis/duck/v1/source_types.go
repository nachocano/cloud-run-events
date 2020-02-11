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

package v1

import (
	"context"
	"fmt"
	"math"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilpointer "k8s.io/utils/pointer"

	"knative.dev/pkg/apis"
	"knative.dev/pkg/apis/duck"
)

// Source is an Implementable "duck type".
var _ duck.Implementable = (*Source)(nil)

// +genduck
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Source is the minimum resource shape to adhere to the Source Specification.
// This duck type is intended to allow implementors of Sources and
// Importers to verify their own resources meet the expectations.
// This is not a real resource.
// NOTE: The Source Specification is in progress and the shape and names could
// be modified until it has been accepted.
type Source struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SourceSpec   `json:"spec"`
	Status SourceStatus `json:"status"`
}

type SourceSpec struct {
	// Sink is a reference to an object that will resolve to a domain name or a
	// URI directly to use as the sink.
	Sink Destination `json:"sink,omitempty"`

	// CloudEventOverrides defines overrides to control the output format and
	// modifications of the event sent to the sink.
	// +optional
	CloudEventOverrides *CloudEventOverrides `json:"ceOverrides,omitempty"`

	// Scaler defines the scaling options for the source, e.g., whether it can
	// scale to zero, the maximum number of pods it can scale to, as well as particular
	// options based on the scaling technology used.
	// If not specified, the source is non-scalable.
	// +optional
	Scaler *ScalerSpec `json:"scaler,omitempty"`
}

// ScalerClass is the class of source scaler that a particular resource has opted into.
type ScalerClass string

const (
	// ScalerClassKeda is the Keda Scaler class.
	ScalerClassKeda ScalerClass = "keda"
	// ScalerClassKsvc is the Knative Service class.
	ScalerClassKsvc ScalerClass = "ksvc"
)

const (
	// defaultScalerClass is the default scaler class.
	defaultScalerClass = ScalerClassKeda
	// defaultMinScale is the default minimum set of Pods the scaler should
	// downscale the source to.
	defaultMinScale int32 = 0
	// defaultMaxScale is the default maximum set of Pods the scaler should
	// upscale the source to.
	defaultMaxScale int32 = 1
)

type ScalerSpec struct {
	// Class defines the class of scaler to use.
	Class ScalerClass `json:"class,omitempty"`

	// MinScale defines the minimum scale for the source.
	// If not specified, defaults to zero.
	// +optional
	MinScale *int32 `json:"minScale,omitempty"`

	// MaxScale defines the maximum scale for the source.
	// If not specified, defaults to one.
	// +optional
	MaxScale *int32 `json:"maxScale,omitempty"`

	// Options defines specific knobs to tune based on the
	// particular scaling backend (e.g., keda or ksvc)
	// +optional
	Options map[string]string `json:"options,omitempty"`
}

// CloudEventOverrides defines arguments for a Source that control the output
// format of the CloudEvents produced by the Source.
type CloudEventOverrides struct {
	// Extensions specify what attribute are added or overridden on the
	// outbound event. Each `Extensions` key-value pair are set on the event as
	// an attribute extension independently.
	// +optional
	Extensions map[string]string `json:"extensions,omitempty"`
}

// SourceStatus shows how we expect folks to embed Addressable in
// their Status field.
type SourceStatus struct {
	// inherits duck/v1beta1 Status, which currently provides:
	// * ObservedGeneration - the 'Generation' of the Service that was last
	//   processed by the controller.
	// * Conditions - the latest available observations of a resource's current
	//   state.
	Status `json:",inline"`

	// SinkURI is the current active sink URI that has been configured for the
	// Source.
	// +optional
	SinkURI *apis.URL `json:"sinkUri,omitempty"`
}

// IsReady returns true if the resource is ready overall.
func (ss *SourceStatus) IsReady() bool {
	for _, c := range ss.Conditions {
		switch c.Type {
		// Look for the "happy" condition, which is the only condition that
		// we can reliably understand to be the overall state of the resource.
		case apis.ConditionReady, apis.ConditionSucceeded:
			return c.IsTrue()
		}
	}
	return false
}

// Validate the ScalerSpec has all the necessary fields.
func (ss *ScalerSpec) Validate(ctx context.Context) *apis.FieldError {
	if ss == nil {
		return nil
	}
	var errs *apis.FieldError
	if ss.Class == "" {
		errs = errs.Also(apis.ErrMissingField("class"))
	}
	if ss.MinScale == nil {
		errs = errs.Also(apis.ErrMissingField("minScale"))
	} else if *ss.MinScale < 0 {
		errs = errs.Also(apis.ErrOutOfBoundsValue(*ss.MinScale, 0, math.MaxInt32, "minScale"))
	}

	if ss.MaxScale == nil {
		errs = errs.Also(apis.ErrMissingField("maxScale"))
	} else if *ss.MaxScale < 1 {
		errs = errs.Also(apis.ErrOutOfBoundsValue(*ss.MaxScale, 1, math.MaxInt32, "maxScale"))
	}

	if ss.MinScale != nil && ss.MaxScale != nil && *ss.MaxScale < *ss.MinScale {
		errs = errs.Also(&apis.FieldError{
			Message: fmt.Sprintf("maxScale=%d is less than minScale=%d", *ss.MaxScale, *ss.MinScale),
			Paths:   []string{"maxScale", "minScale"},
		})
	}

	return errs
}

func (ss *ScalerSpec) SetDefault(ctx context.Context) {
	if ss == nil {
		return
	}
	if ss.Class == "" {
		ss.Class = defaultScalerClass
	}
	if ss.MinScale == nil {
		ss.MinScale = utilpointer.Int32Ptr(defaultMinScale)
	}
	if ss.MaxScale == nil {
		ss.MaxScale = utilpointer.Int32Ptr(defaultMaxScale)
	}
}

var (
	// Verify Source resources meet duck contracts.
	_ duck.Populatable = (*Source)(nil)
	_ apis.Listable    = (*Source)(nil)
)

const (
	// SourceConditionSinkProvided has status True when the Source
	// has been configured with a sink target that is resolvable.
	SourceConditionSinkProvided apis.ConditionType = "SinkProvided"

	// SourceScalerProvided has status True when the Source
	// has been configured with an scaler.
	SourceScalerProvided apis.ConditionType = "ScalerProvided"
)

// GetFullType implements duck.Implementable
func (*Source) GetFullType() duck.Populatable {
	return &Source{}
}

// Populate implements duck.Populatable
func (s *Source) Populate() {
	s.Spec.Sink = Destination{
		URI: &apis.URL{
			Scheme:   "https",
			Host:     "tableflip.dev",
			RawQuery: "flip=mattmoor",
		},
	}
	s.Spec.CloudEventOverrides = &CloudEventOverrides{
		Extensions: map[string]string{"boosh": "kakow"},
	}
	s.Spec.Scaler = &ScalerSpec{
		Class:    ScalerClassKsvc,
		MinScale: utilpointer.Int32Ptr(0),
		MaxScale: utilpointer.Int32Ptr(1),
		Options:  map[string]string{"myoption": "myoptionvalue"},
	}
	s.Status.ObservedGeneration = 42
	s.Status.Conditions = Conditions{{
		// Populate ALL fields
		Type:               SourceConditionSinkProvided,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: apis.VolatileTime{Inner: metav1.NewTime(time.Date(1984, 02, 28, 18, 52, 00, 00, time.UTC))},
	}, {
		Type:               SourceScalerProvided,
		Status:             corev1.ConditionTrue,
		LastTransitionTime: apis.VolatileTime{Inner: metav1.NewTime(time.Date(1984, 02, 28, 18, 52, 00, 00, time.UTC))},
	}}
	s.Status.SinkURI = &apis.URL{
		Scheme:   "https",
		Host:     "tableflip.dev",
		RawQuery: "flip=mattmoor",
	}
}

// IsScalable returns true if the SourceSpec has been configured with scaling options.
func (ss *SourceSpec) IsScalable() bool {
	return ss.Scaler != nil
}

// GetListType implements apis.Listable
func (*Source) GetListType() runtime.Object {
	return &SourceList{}
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SourceList is a list of Source resources
type SourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Source `json:"items"`
}
