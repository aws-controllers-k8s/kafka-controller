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

package vpc_connection

import (
	"context"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"

	"github.com/aws-controllers-k8s/kafka-controller/pkg/sync"
)

// customUpdate handles updates to the VpcConnection resource. VpcConnection has
// no UpdateVpcConnection API and all spec fields are immutable. Only tags can be
// mutated, via the TagResource/UntagResource APIs.
func (rm *resourceManager) customUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (updated *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.customUpdate")
	defer func() { exit(err) }()

	// Construct the updated resource from the desired spec fields and copy the
	// status from latest.
	updatedRes := rm.concreteResource(desired.DeepCopy())
	updatedRes.ko.Status = latest.ko.Status

	if delta.DifferentAt("Spec.Tags") {
		err = sync.Tags(
			ctx,
			desired.ko.Spec.Tags, latest.ko.Spec.Tags,
			(*string)(latest.ko.Status.ACKResourceMetadata.ARN),
			convertToOrderedACKTags, rm.sdkapi, rm.metrics,
		)
		if err != nil {
			return updatedRes, err
		}
	}

	return updatedRes, nil
}
