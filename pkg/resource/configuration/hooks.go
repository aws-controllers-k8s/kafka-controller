// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package configuration

import (
	"context"

	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/kafka"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/kafka/types"

	"github.com/aws-controllers-k8s/kafka-controller/pkg/adoption"
)

// shouldLookupARN returns true when the supplied resource needs its ARN
// resolved from its name before a read can be attempted.
func (rm *resourceManager) shouldLookupARN(r *resource) bool {
	return adoption.ShouldLookupARN(r)
}

// getConfigurationARN looks up the ARN of the configuration named by the
// supplied resource's Spec.Name. MSK configuration ARNs embed a
// service-generated UUID, so they cannot be constructed from the name alone and
// must be read back from the API.
//
// Unlike ListClusters, ListConfigurations offers no name filter, so this walks
// every configuration in the region. That is why callers gate the lookup on
// shouldLookupARN: only a resource whose owner opted into adoption pays for it.
//
// Configuration names are unique within an account and region, so at most one
// configuration can match. Returns nil (without error) when none matches, which
// callers translate into ackerr.NotFound.
func (rm *resourceManager) getConfigurationARN(
	ctx context.Context,
	r *resource,
) (arn *string, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.getConfigurationARN")
	defer func() { exit(err) }()

	// Spec.Name is required for create, but an adopted resource may be
	// identified solely by its ARN, in which case there is nothing to look up.
	if r.ko.Spec.Name == nil {
		return nil, nil
	}

	paginator := svcsdk.NewListConfigurationsPaginator(
		rm.sdkapi, &svcsdk.ListConfigurationsInput{},
	)
	for paginator.HasMorePages() {
		resp, err := paginator.NextPage(ctx)
		rm.metrics.RecordAPICall("READ_MANY", "ListConfigurations", err)
		if err != nil {
			return nil, err
		}
		if arn := findConfigurationARNByName(resp.Configurations, *r.ko.Spec.Name); arn != nil {
			return arn, nil
		}
	}
	return nil, nil
}

// findConfigurationARNByName returns the ARN of the configuration named exactly
// name, or nil if the supplied page holds no such configuration.
func findConfigurationARNByName(
	configurations []svcsdktypes.Configuration,
	name string,
) *string {
	for _, c := range configurations {
		if c.Name == nil || *c.Name != name {
			continue
		}
		return c.Arn
	}
	return nil
}
