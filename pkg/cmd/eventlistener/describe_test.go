// Copyright © 2020 The Tekton Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package eventlistener

// TODO: properly move to v1beta1
import (
	"fmt"
	"testing"

	"github.com/tektoncd/cli/pkg/test"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
	triggerV1beta1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1beta1"
	triggertest "github.com/tektoncd/triggers/test"
	"gotest.tools/v3/golden"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/pkg/apis"
	duckv1alpha1 "knative.dev/pkg/apis/duck/v1beta1"
)

func TestEventListenerDescribe_InvalidNamespace(t *testing.T) {
	cs := test.SeedTestResources(t, triggertest.Resources{})
	p := &test.Params{Triggers: cs.Triggers, Kube: cs.Kube}

	eventListener := Command(p)
	out, err := test.ExecuteCommand(eventListener, "desc", "bar", "-n", "invalid")
	if err == nil {
		t.Errorf("Error expected here")
	}
	golden.Assert(t, out, fmt.Sprintf("%s.golden", t.Name()))
}

func TestEventListenerDescribe_NonExistedName(t *testing.T) {
	cs := test.SeedTestResources(t, triggertest.Resources{Namespaces: []*corev1.Namespace{{
		ObjectMeta: metav1.ObjectMeta{
			Name: "ns",
		},
	}}})
	p := &test.Params{Triggers: cs.Triggers, Kube: cs.Kube}

	eventListener := Command(p)
	out, err := test.ExecuteCommand(eventListener, "desc", "bar", "-n", "ns")
	if err == nil {
		t.Errorf("Error expected here")
	}
	golden.Assert(t, out, fmt.Sprintf("%s.golden", t.Name()))
}

func TestEventListenerDescribe_NoArgProvided(t *testing.T) {
	cs := test.SeedTestResources(t, triggertest.Resources{Namespaces: []*corev1.Namespace{{
		ObjectMeta: metav1.ObjectMeta{
			Name: "ns",
		},
	}}})
	p := &test.Params{Triggers: cs.Triggers, Kube: cs.Kube}

	eventListener := Command(p)
	out, err := test.ExecuteCommand(eventListener, "desc", "-n", "ns")
	if err == nil {
		t.Errorf("Error expected here")
	}
	golden.Assert(t, out, fmt.Sprintf("%s.golden", t.Name()))
}

func TestEventListenerDescribe_WithMinRequiredField(t *testing.T) {
	els := []*triggerV1beta1.EventListener{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "el1",
				Namespace: "ns",
			},
		},
	}

	executeEventListenerCommand(t, els)
}

func TestEventListenerDescribe_OneTriggerWithOneClusterTriggerBinding(t *testing.T) {
	triggerTemplateRef := "tt1"
	els := []*triggerV1beta1.EventListener{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "el1",
				Namespace: "ns",
			},
			Spec: triggerV1beta1.EventListenerSpec{
				Triggers: []triggerV1beta1.EventListenerTrigger{
					{
						Template: &triggerV1beta1.TriggerSpecTemplate{
							Ref:        &triggerTemplateRef,
							APIVersion: "v1beta1",
						},
						Bindings: []*triggerV1beta1.EventListenerBinding{
							{
								Ref:        "tb1",
								Kind:       "ClusterTriggerBinding",
								APIVersion: "v1alpha1",
							},
						},
					},
				},
			},
		},
	}

	executeEventListenerCommand(t, els)
}

func TestEventListenerDescribe_WithOutputStatusURLAndName(t *testing.T) {
	triggerTemplateRef := "tt1"
	els := []*triggerV1beta1.EventListener{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "el1",
				Namespace: "ns",
			},
			Spec: triggerV1beta1.EventListenerSpec{
				Triggers: []triggerV1beta1.EventListenerTrigger{
					{
						Template: &triggerV1beta1.TriggerSpecTemplate{
							Ref:        &triggerTemplateRef,
							APIVersion: "v1beta1",
						},
						Bindings: []*triggerV1beta1.EventListenerBinding{
							{
								Ref:        "tb1",
								Kind:       "ClusterTriggerBinding",
								APIVersion: "v1beta1",
							},
						},
					},
				},
			},
			Status: triggerV1beta1.EventListenerStatus{
				AddressStatus: duckv1alpha1.AddressStatus{
					Address: &duckv1alpha1.Addressable{
						URL: &apis.URL{
							Scheme: "http",
							Host:   "el-listener.default.svc.cluster.local",
						},
					},
				},
				Configuration: triggerV1beta1.EventListenerConfig{
					GeneratedResourceName: "el-listener",
				},
			},
		},
	}

	executeEventListenerCommand(t, els)
}

func TestEventListenerDescribe_OneTriggerWithOneTriggerBinding(t *testing.T) {
	triggerTemplateRef := "tt1"

	els := []*triggerV1beta1.EventListener{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "el1",
				Namespace: "ns",
			},
			Spec: triggerV1beta1.EventListenerSpec{
				ServiceAccountName: "trigger-sa",
				Triggers: []triggerV1beta1.EventListenerTrigger{
					{
						Template: &triggerV1beta1.TriggerSpecTemplate{
							Ref:        &triggerTemplateRef,
							APIVersion: "v1beta1",
						},
						Bindings: []*triggerV1beta1.EventListenerBinding{
							{
								Ref:        "tb1",
								Kind:       "TriggerBinding",
								APIVersion: "",
							},
						},
					},
				},
			},
		},
	}

	executeEventListenerCommand(t, els)
}

func TestEventListenerDescribe_OneTriggerWithMultipleTriggerBinding(t *testing.T) {
	triggerTemplateRef := "tt1"
	els := []*triggerV1beta1.EventListener{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "el1",
				Namespace: "ns",
			},
			Spec: triggerV1beta1.EventListenerSpec{
				Triggers: []triggerV1beta1.EventListenerTrigger{
					{
						Template: &triggerV1beta1.TriggerSpecTemplate{
							Ref:        &triggerTemplateRef,
							APIVersion: "v1beta1",
						},
						Bindings: []*triggerV1beta1.EventListenerBinding{
							{
								Ref:        "tb1",
								Kind:       "TriggerBinding",
								APIVersion: "",
							},
							{
								Ref:        "tb2",
								Kind:       "ClusterTriggerBinding",
								APIVersion: "v1beta1",
							},
							{
								Ref:        "tb3",
								Kind:       "TriggerBinding",
								APIVersion: "v1beta1",
							},
						},
					},
				},
			},
		},
	}

	executeEventListenerCommand(t, els)
}

func TestEventListenerDescribe_OneTriggerWithTriggerBindingName(t *testing.T) {
	bindingval := "somevalue"

	els := []*triggerV1beta1.EventListener{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "el1",
				Namespace: "ns",
			},
			Spec: triggerV1beta1.EventListenerSpec{
				Triggers: []triggerV1beta1.EventListenerTrigger{
					{
						Bindings: []*triggerV1beta1.EventListenerBinding{
							{
								Name:  "binding",
								Value: &bindingval,
							},
						},
						Template: &triggerV1beta1.EventListenerTemplate{
							Ref:        nil,
							APIVersion: "v1alpha1",
							Spec:       nil,
						},
						Name: "tt1",
					},
				},
			},
		},
	}

	executeEventListenerCommand(t, els)
}

func TestEventListenerDescribe_TriggerWithTriggerTemplateRef(t *testing.T) {
	bindingval := "somevalue"
	tempRef := "someref"

	els := []*triggerV1beta1.EventListener{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "el1",
				Namespace: "ns",
			},
			Spec: triggerV1beta1.EventListenerSpec{
				Triggers: []triggerV1beta1.EventListenerTrigger{
					{
						Bindings: []*triggerV1beta1.EventListenerBinding{
							{
								Name:  "binding",
								Value: &bindingval,
							},
						},
						Template: &triggerV1beta1.EventListenerTemplate{
							Ref:        &tempRef,
							APIVersion: "v1beta1",
						},
						Name: "tt1",
					},
				},
			},
		},
	}

	executeEventListenerCommand(t, els)
}

func TestEventListenerDescribe_TriggerWithTriggerTemplateRefTriggerRef(t *testing.T) {
	bindingval := "somevalue"
	tempRef := "someref"

	els := []*triggerV1beta1.EventListener{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "el1",
				Namespace: "ns",
			},
			Spec: triggerV1beta1.EventListenerSpec{
				Triggers: []triggerV1beta1.EventListenerTrigger{
					{
						Bindings: []*triggerV1beta1.EventListenerBinding{
							{
								Name:  "binding",
								Value: &bindingval,
							},
						},
						Template: &triggerV1beta1.EventListenerTemplate{
							Ref:        &tempRef,
							APIVersion: "v1alpha1",
						},
						Name:               "tt1",
						TriggerRef:         "triggeref",
						ServiceAccountName: "test-sa",
						Interceptors: []*triggerV1beta1.EventInterceptor{
							{
								Webhook: &triggerV1beta1.WebhookInterceptor{
									ObjectRef: &corev1.ObjectReference{
										Kind:       "Service",
										Name:       "testwebhook",
										APIVersion: "v1",
										Namespace:  "ns",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	executeEventListenerCommand(t, els)
}

func TestEventListenerDescribe_OneTriggerWithEmptyTriggerBinding(t *testing.T) {
	triggerTemplateRef := "tt1"

	els := []*triggerV1beta1.EventListener{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "el1",
				Namespace: "ns",
			},
			Spec: triggerV1beta1.EventListenerSpec{
				Triggers: []triggerV1beta1.EventListenerTrigger{
					{
						Template: &triggerV1beta1.TriggerSpecTemplate{
							Ref:        &triggerTemplateRef,
							APIVersion: "v1beta1",
						},
					},
				},
			},
		},
	}

	executeEventListenerCommand(t, els)
}

func TestEventListenerDescribe_MultipleTriggers(t *testing.T) {
	triggerTemplateRef1 := "tt1"
	triggerTemplateRef2 := "tt2"

	els := []*triggerV1beta1.EventListener{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "el1",
				Namespace: "ns",
			},
			Spec: triggerV1beta1.EventListenerSpec{
				Triggers: []triggerV1beta1.EventListenerTrigger{
					{
						Template: &triggerV1beta1.TriggerSpecTemplate{
							Ref:        &triggerTemplateRef1,
							APIVersion: "v1beta1",
						},
						Bindings: []*triggerV1beta1.EventListenerBinding{
							{
								Ref:        "tb1",
								Kind:       "TriggerBinding",
								APIVersion: "",
							},
							{
								Ref:        "tb2",
								Kind:       "ClusterTriggerBinding",
								APIVersion: "v1beta1",
							},
							{
								Ref:        "tb3",
								Kind:       "TriggerBinding",
								APIVersion: "v1beta1",
							},
						},
					},
					{
						Template: &triggerV1beta1.TriggerSpecTemplate{
							Ref:        &triggerTemplateRef2,
							APIVersion: "v1beta1",
						},
						ServiceAccountName: "sa1",
						Bindings: []*triggerV1beta1.EventListenerBinding{
							{
								Ref:        "tb1",
								Kind:       "TriggerBinding",
								APIVersion: "",
							},
							{
								Ref:        "tb2",
								Kind:       "ClusterTriggerBinding",
								APIVersion: "v1beta1",
							},
							{
								Ref:        "tb3",
								Kind:       "TriggerBinding",
								APIVersion: "v1beta1",
							},
						},
					},
				},
			},
		},
	}

	executeEventListenerCommand(t, els)
}

func TestEventListenerDescribe_WithWebhookInterceptors(t *testing.T) {
	triggerTemplateRef := "tt1"
	els := []*triggerV1beta1.EventListener{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "el1",
				Namespace: "ns",
			},
			Status: triggerV1beta1.EventListenerStatus{
				AddressStatus: duckv1alpha1.AddressStatus{
					Address: &duckv1alpha1.Addressable{
						URL: &apis.URL{
							Scheme: "http",
							Host:   "el-listener.default.svc.cluster.local",
						},
					},
				},
				Configuration: triggerV1beta1.EventListenerConfig{
					GeneratedResourceName: "el-listener",
				},
			},
			Spec: triggerV1beta1.EventListenerSpec{
				Triggers: []triggerV1beta1.EventListenerTrigger{

					{
						Template: &triggerV1beta1.TriggerSpecTemplate{
							Ref:        &triggerTemplateRef,
							APIVersion: "v1beta1",
						},
						Name: "foo-trig",
						Bindings: []*triggerV1beta1.EventListenerBinding{
							{
								Ref:        "tb1",
								Kind:       "TriggerBinding",
								APIVersion: "",
							},
							{
								Ref:        "tb2",
								Kind:       "ClusterTriggerBinding",
								APIVersion: "v1beta1",
							},
							{
								Ref:        "tb3",
								Kind:       "TriggerBinding",
								APIVersion: "v1beta1",
							},
						},
						Interceptors: []*triggerV1beta1.TriggerInterceptor{
							{
								Webhook: &triggerV1beta1.WebhookInterceptor{
									ObjectRef: &corev1.ObjectReference{
										Kind:       "Service",
										Name:       "webhookTest",
										Namespace:  "namespace",
										APIVersion: "v1",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	executeEventListenerCommand(t, els)
}

func TestEventListenerDescribe_WithWebhookInterceptorsWithParams(t *testing.T) {
	triggerTemplateRef := "tt1"
	els := []*triggerV1beta1.EventListener{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "el1",
				Namespace: "ns",
			},
			Spec: triggerV1beta1.EventListenerSpec{
				Triggers: []triggerV1beta1.EventListenerTrigger{

					{
						Template: &triggerV1beta1.TriggerSpecTemplate{
							Ref:        &triggerTemplateRef,
							APIVersion: "v1betaa1",
						},
						Name: "foo-trig",
						Bindings: []*triggerV1beta1.EventListenerBinding{
							{
								Ref:        "tb1",
								Kind:       "TriggerBinding",
								APIVersion: "",
							},
							{
								Ref:        "tb2",
								Kind:       "ClusterTriggerBinding",
								APIVersion: "v1alpha1",
							},
							{
								Ref:        "tb3",
								Kind:       "TriggerBinding",
								APIVersion: "v1alpha1",
							},
						},
						Interceptors: []*triggerV1beta1.TriggerInterceptor{
							{
								Webhook: &triggerV1beta1.WebhookInterceptor{
									ObjectRef: &corev1.ObjectReference{
										Kind:       "Service",
										Name:       "foo",
										Namespace:  "namespace",
										APIVersion: "v1",
									},
									Header: []v1beta1.Param{
										{
											Name: "header",
											Value: v1beta1.ArrayOrString{
												Type:     "array",
												ArrayVal: []string{"value"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	executeEventListenerCommand(t, els)
}

func TestEventListenerDescribe_MultipleTriggerWithTriggerRefAndTriggerSpec(t *testing.T) {
	triggerTemplateRef := "tt1"
	els := []*triggerV1beta1.EventListener{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "el1",
				Namespace: "ns",
			},
			Spec: triggerV1beta1.EventListenerSpec{
				Triggers: []triggerV1beta1.EventListenerTrigger{
					{
						TriggerRef: "test-ref",
					},
					{
						Template: &triggerV1beta1.TriggerSpecTemplate{
							Ref:        &triggerTemplateRef,
							APIVersion: "v1alpha1",
						},
						Name: "foo-trig",

						Bindings: []*triggerV1beta1.EventListenerBinding{
							{
								Ref:        "tb1",
								Kind:       "TriggerBinding",
								APIVersion: "",
							},
							{
								Ref:        "tb2",
								Kind:       "ClusterTriggerBinding",
								APIVersion: "v1alpha1",
							},
							{
								Ref:        "tb3",
								Kind:       "TriggerBinding",
								APIVersion: "v1alpha1",
							},
						},
						Interceptors: []*triggerV1beta1.TriggerInterceptor{
							{
								Webhook: &triggerV1beta1.WebhookInterceptor{
									ObjectRef: &corev1.ObjectReference{
										Kind:       "Service",
										Name:       "foo",
										Namespace:  "namespace",
										APIVersion: "v1",
									},
									Header: []v1beta1.Param{
										{
											Name: "header",
											Value: v1beta1.ArrayOrString{
												Type:     "array",
												ArrayVal: []string{"value"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	executeEventListenerCommand(t, els)
}

func TestEventListenerDescribe_WithCELInterceptors(t *testing.T) {
	triggerTemplateRef := "tt1"
	els := []*triggerV1beta1.EventListener{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "el1",
				Namespace: "ns",
			},
			Spec: triggerV1beta1.EventListenerSpec{
				Triggers: []triggerV1beta1.EventListenerTrigger{

					{
						Template: &triggerV1beta1.TriggerSpecTemplate{
							Ref:        &triggerTemplateRef,
							APIVersion: "v1alpha1",
						},
						Name: "foo-trig",
						Bindings: []*triggerV1beta1.EventListenerBinding{
							{
								Ref:        "tb1",
								Kind:       "TriggerBinding",
								APIVersion: "",
							},
							{
								Ref:        "tb2",
								Kind:       "ClusterTriggerBinding",
								APIVersion: "v1alpha1",
							},
							{
								Ref:        "tb3",
								Kind:       "TriggerBinding",
								APIVersion: "v1alpha1",
							},
						},
						Interceptors: []*triggerV1beta1.TriggerInterceptor{
							{
								/*DeprecatedCEL: &v1alpha1.CELInterceptor{
									Filter: "body.value == 'test'",
									Overlays: []v1alpha1.CELOverlay{
										{
											Key:        "value",
											Expression: "'testing'",
										},
									},
								},*/
							},
						},
					},
				},
			},
		},
	}

	executeEventListenerCommand(t, els)
}

func TestEventListenerDescribe_WithMultipleBindingAndInterceptors(t *testing.T) {
	triggerTemplateRef1 := "tt1"
	triggerTemplateRef2 := "tt2"
	interceptorName := "interceptor-one"
	els := []*triggerV1beta1.EventListener{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "el1",
				Namespace: "ns",
			},
			Spec: triggerV1beta1.EventListenerSpec{
				Resources: triggerV1beta1.Resources{
					KubernetesResource: &triggerV1beta1.KubernetesResource{
						ServiceType: "ClusterIP",
					},
				},
				Triggers: []triggerV1beta1.EventListenerTrigger{

					{
						Template: &triggerV1beta1.TriggerSpecTemplate{
							Ref:        &triggerTemplateRef1,
							APIVersion: "v1alpha1",
						},
						Name: "foo-trig",
						Bindings: []*triggerV1beta1.EventListenerBinding{
							{
								Ref:        "tb1",
								Kind:       "TriggerBinding",
								APIVersion: "",
							},
							{
								Ref:        "tb2",
								Kind:       "ClusterTriggerBinding",
								APIVersion: "v1alpha1",
							},
							{
								Ref:        "",
								Kind:       "TriggerBinding",
								APIVersion: "v1alpha1",
							},
						},
						Interceptors: []*triggerV1beta1.TriggerInterceptor{
							{
								Name: &interceptorName,
								Ref: triggerV1beta1.InterceptorRef{
									Name:       "cel",
									Kind:       triggerV1beta1.ClusterInterceptorKind,
									APIVersion: "v1beta1",
								},
								/*DeprecatedCEL: &v1alpha1.CELInterceptor{
									Filter: "body.value == 'test'",
									Overlays: []v1alpha1.CELOverlay{
										{
											Key:        "value",
											Expression: "'testing'",
										},
									},
								},*/
							},
						},
					},
					{
						Template: &triggerV1beta1.TriggerSpecTemplate{
							Ref:        &triggerTemplateRef2,
							APIVersion: "v1alpha1",
						},
						ServiceAccountName: "sa1",
						Name:               "foo-trig",
						Bindings: []*triggerV1beta1.EventListenerBinding{
							{
								Ref:        "tb4",
								Kind:       "TriggerBinding",
								APIVersion: "",
							},
							{
								Ref:        "tb5",
								Kind:       "ClusterTriggerBinding",
								APIVersion: "v1alpha1",
							},
						},
						Interceptors: []*triggerV1beta1.TriggerInterceptor{
							{
								Webhook: &triggerV1beta1.WebhookInterceptor{
									ObjectRef: &corev1.ObjectReference{
										Kind:       "Service",
										Name:       "webhookTest",
										Namespace:  "namespace",
										APIVersion: "v1",
									},
								},
							},
							{
								
							}
							{
								/*DeprecatedCEL: &v1alpha1.CELInterceptor{
									Filter: "body.value == 'test'",
									Overlays: []v1alpha1.CELOverlay{
										{
											Key:        "value",
											Expression: "'testing'",
										},
									},
								},*/
							},
						},
					},
				},
			},
		},
	}
	executeEventListenerCommand(t, els)
}

func TestEventListenerDescribe_OutputYAMLWithMultipleBindingAndInterceptors(t *testing.T) {
	triggerTemplateRef1 := "tt1"
	triggerTemplateRef2 := "tt2"
	els := []*triggerV1beta1.EventListener{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "el1",
				Namespace: "ns",
			},
			Spec: triggerV1beta1.EventListenerSpec{
				Triggers: []triggerV1beta1.EventListenerTrigger{

					{
						Template: &triggerV1beta1.TriggerSpecTemplate{
							Ref:        &triggerTemplateRef1,
							APIVersion: "v1beta1",
						},
						Name: "foo-trig",
						Bindings: []*triggerV1beta1.EventListenerBinding{
							{
								Ref:        "tb1",
								Kind:       "TriggerBinding",
								APIVersion: "",
							},
							{
								Ref:        "tb2",
								Kind:       "ClusterTriggerBinding",
								APIVersion: "v1beta1",
							},
							{
								Ref:        "tb3",
								Kind:       "TriggerBinding",
								APIVersion: "v1beta1",
							},
						},
						Interceptors: []*triggerV1beta1.TriggerInterceptor{
							{
								/*DeprecatedCEL: &v1alpha1.CELInterceptor{
									Filter: "body.value == 'test'",
									Overlays: []v1alpha1.CELOverlay{
										{
											Key:        "value",
											Expression: "'testing'",
										},
									},
								},*/
							},
						},
					},
					{
						Template: &triggerV1beta1.TriggerSpecTemplate{
							Ref:        &triggerTemplateRef2,
							APIVersion: "v1alpha1",
						},
						ServiceAccountName: "sa1",
						Name:               "foo-trig",
						Bindings: []*triggerV1beta1.EventListenerBinding{
							{
								Ref:        "tb4",
								Kind:       "TriggerBinding",
								APIVersion: "",
							},
							{
								Ref:        "tb5",
								Kind:       "ClusterTriggerBinding",
								APIVersion: "v1alpha1",
							},
						},
						Interceptors: []*triggerV1beta1.TriggerInterceptor{
							{
								Webhook: &triggerV1beta1.WebhookInterceptor{
									ObjectRef: &corev1.ObjectReference{
										Kind:       "Service",
										Name:       "webhookTest",
										Namespace:  "namespace",
										APIVersion: "v1",
									},
								},
							},
							{
								/*DeprecatedCEL: &v1alpha1.CELInterceptor{
									Filter: "body.value == 'test'",
									Overlays: []v1alpha1.CELOverlay{
										{
											Key:        "value",
											Expression: "'testing'",
										},
									},
								},*/
							},
						},
					},
				},
			},
		},
	}

	cs := test.SeedTestResources(t, triggertest.Resources{EventListeners: els, Namespaces: []*corev1.Namespace{{
		ObjectMeta: metav1.ObjectMeta{
			Name: "ns",
		},
	}}})
	p := &test.Params{Triggers: cs.Triggers, Kube: cs.Kube}

	eventListener := Command(p)
	out, err := test.ExecuteCommand(eventListener, "desc", "el1", "-n", "ns", "-o", "json")
	if err != nil {
		t.Errorf("Error expected here")
	}
	golden.Assert(t, out, fmt.Sprintf("%s.golden", t.Name()))
}

func TestEventListenerDescribe_WithOutputStatusURL(t *testing.T) {
	els := []*triggerV1beta1.EventListener{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "el1",
				Namespace: "ns",
			},

			Status: triggerV1beta1.EventListenerStatus{
				AddressStatus: duckv1alpha1.AddressStatus{
					Address: &duckv1alpha1.Addressable{
						URL: &apis.URL{
							Scheme: "http",
							Host:   "el-listener.default.svc.cluster.local",
						},
					},
				},
			},
		},
	}

	cs := test.SeedTestResources(t, triggertest.Resources{EventListeners: els})

	p := &test.Params{Triggers: cs.Triggers, Kube: cs.Kube}

	eventListener := Command(p)
	out, err := test.ExecuteCommand(eventListener, "desc", "el1", "-o", "url", "-n", "ns")
	if err != nil {
		t.Errorf("Error")
	}
	test.AssertOutput(t, "http://el-listener.default.svc.cluster.local\n", out)
}

func TestEventListenerDescribe_OutputStatusURL_WithNoURL(t *testing.T) {
	els := []*triggerV1beta1.EventListener{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "el1",
				Namespace: "ns",
			},
		},
	}

	cs := test.SeedTestResources(t, triggertest.Resources{EventListeners: els})

	p := &test.Params{Triggers: cs.Triggers, Kube: cs.Kube}

	eventListener := Command(p)
	out, err := test.ExecuteCommand(eventListener, "desc", "el1", "-o", "url", "-n", "ns")
	if err == nil {
		t.Errorf("Error")
	}

	test.AssertOutput(t, "Error: "+err.Error()+"\n", out)
}

func TestEventListenerDescribe_AutoSelect(t *testing.T) {
	triggerTemplateRef := "tt1"
	els := []*triggerV1beta1.EventListener{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "el1",
				Namespace: "ns",
			},
			Spec: triggerV1beta1.EventListenerSpec{
				Triggers: []triggerV1beta1.EventListenerTrigger{
					{
						Template: &triggerV1beta1.TriggerSpecTemplate{
							Ref:        &triggerTemplateRef,
							APIVersion: "v1beta1",
						},
					},
				},
			},
		},
	}

	cs := test.SeedTestResources(t, triggertest.Resources{EventListeners: els, Namespaces: []*corev1.Namespace{{
		ObjectMeta: metav1.ObjectMeta{
			Name: "ns",
		},
	}}})
	p := &test.Params{Triggers: cs.Triggers, Kube: cs.Kube}

	eventListener := Command(p)
	out, err := test.ExecuteCommand(eventListener, "desc", "-n", "ns")
	if err != nil {
		t.Errorf("Error expected here")
	}
	golden.Assert(t, out, fmt.Sprintf("%s.golden", t.Name()))
}

func executeEventListenerCommand(t *testing.T, els []*triggerV1beta1.EventListener) {
	cs := test.SeedTestResources(t, triggertest.Resources{EventListeners: els, Namespaces: []*corev1.Namespace{{
		ObjectMeta: metav1.ObjectMeta{
			Name: "ns",
		},
	}}})
	p := &test.Params{Triggers: cs.Triggers, Kube: cs.Kube}

	eventListener := Command(p)
	out, err := test.ExecuteCommand(eventListener, "desc", "el1", "-n", "ns")
	if err != nil {
		t.Errorf("Error expected here")
	}
	golden.Assert(t, out, fmt.Sprintf("%s.golden", t.Name()))
}
