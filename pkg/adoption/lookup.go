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

package adoption

import (
	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	acktypes "github.com/aws-controllers-k8s/runtime/pkg/types"
)

// ShouldLookupARN returns true when the supplied resource has no ARN but does
// carry an adoption-policy annotation.
//
// MSK identifies clusters and configurations by ARNs that embed a
// service-generated UUID, so an ARN cannot be derived from the resource name
// and has to be read back from the API. Resolving it costs a List call, so the
// lookup is limited to resources whose owner opted in by setting
// services.k8s.aws/adoption-policy. An ordinary resource being created for the
// first time also has no ARN, and must not pay for a lookup that would find
// nothing.
func ShouldLookupARN(res acktypes.AWSResource) bool {
	if res == nil {
		return false
	}

	if arn := res.Identifiers().ARN(); arn != nil && *arn != ackv1alpha1.AWSResourceName("") {
		return false
	}

	mo := res.MetaObject()
	if mo == nil {
		return false
	}
	return mo.GetAnnotations()[ackv1alpha1.AnnotationAdoptionPolicy] != ""
}
