// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

package enterprisesearch

import (
	entv1beta1 "github.com/elastic/cloud-on-k8s/pkg/apis/enterprisesearch/v1beta1"
	"github.com/elastic/cloud-on-k8s/test/e2e/cmd/run"
	"github.com/elastic/cloud-on-k8s/test/e2e/test"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

func (b Builder) InitTestSteps(k *test.K8sClient) test.StepList {
	return []test.Step{
		{
			Name: "K8S should be accessible",
			Test: test.Eventually(func() error {
				pods := corev1.PodList{}
				return k.Client.List(&pods)
			}),
		},
		{
			Name: "Label test pods",
			Test: test.Eventually(func() error {
				return test.LabelTestPods(
					k.Client,
					test.Ctx(),
					run.TestNameLabel,
					b.EnterpriseSearch.Labels[run.TestNameLabel])
			}),
			Skip: func() bool {
				return test.Ctx().Local
			},
		},
		{
			Name: "EnterpriseSearch CRDs should exist",
			Test: test.Eventually(func() error {
				crds := []runtime.Object{
					&entv1beta1.EnterpriseSearchList{},
				}
				for _, crd := range crds {
					if err := k.Client.List(crd); err != nil {
						return err
					}
				}
				return nil
			}),
		},
		{
			Name: "Remove EnterpriseSearch if it already exists",
			Test: test.Eventually(func() error {
				for _, obj := range b.RuntimeObjects() {
					err := k.Client.Delete(obj)
					if err != nil && !apierrors.IsNotFound(err) {
						return err
					}
				}
				// wait for pods to disappear
				return k.CheckPodCount(0, test.EnterpriseSearchPodListOptions(b.EnterpriseSearch.Namespace, b.EnterpriseSearch.Name)...)
			}),
		},
	}
}
