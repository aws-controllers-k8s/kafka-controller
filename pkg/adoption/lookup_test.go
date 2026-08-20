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

package adoption_test

import (
	"testing"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	acktypes "github.com/aws-controllers-k8s/runtime/pkg/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rtclient "sigs.k8s.io/controller-runtime/pkg/client"

	svcapitypes "github.com/aws-controllers-k8s/kafka-controller/apis/v1alpha1"
	"github.com/aws-controllers-k8s/kafka-controller/pkg/adoption"
)

type fakeResource struct {
	acktypes.AWSResource
	arn         *ackv1alpha1.AWSResourceName
	annotations map[string]string
	nilMeta     bool
}

func (f *fakeResource) Identifiers() acktypes.AWSResourceIdentifiers {
	return &fakeIdentifiers{arn: f.arn}
}

func (f *fakeResource) MetaObject() metav1.Object {
	if f.nilMeta {
		return nil
	}
	return &metav1.ObjectMeta{Annotations: f.annotations}
}

func (f *fakeResource) RuntimeObject() rtclient.Object {
	return &svcapitypes.Cluster{}
}

type fakeIdentifiers struct {
	arn *ackv1alpha1.AWSResourceName
}

func (f *fakeIdentifiers) ARN() *ackv1alpha1.AWSResourceName         { return f.arn }
func (f *fakeIdentifiers) OwnerAccountID() *ackv1alpha1.AWSAccountID { return nil }
func (f *fakeIdentifiers) Region() *ackv1alpha1.AWSRegion            { return nil }
func (f *fakeIdentifiers) Partition() *ackv1alpha1.AWSPartition      { return nil }

func arnPtr(s string) *ackv1alpha1.AWSResourceName {
	arn := ackv1alpha1.AWSResourceName(s)
	return &arn
}

func TestShouldLookupARN(t *testing.T) {
	adoptPolicy := map[string]string{
		ackv1alpha1.AnnotationAdoptionPolicy: "adopt",
	}
	adoptOrCreatePolicy := map[string]string{
		ackv1alpha1.AnnotationAdoptionPolicy: "adopt-or-create",
	}

	tests := []struct {
		name string
		res  acktypes.AWSResource
		want bool
	}{
		{
			name: "nil resource",
			res:  nil,
			want: false,
		},
		{
			name: "no arn, adopt policy",
			res:  &fakeResource{arn: nil, annotations: adoptPolicy},
			want: true,
		},
		{
			name: "no arn, adopt-or-create policy",
			res:  &fakeResource{arn: nil, annotations: adoptOrCreatePolicy},
			want: true,
		},
		{
			name: "no arn, no adoption policy",
			res:  &fakeResource{arn: nil, annotations: nil},
			want: false,
		},
		{
			name: "no arn, unrelated annotations only",
			res: &fakeResource{arn: nil, annotations: map[string]string{
				ackv1alpha1.AnnotationDeletionPolicy: "retain",
			}},
			want: false,
		},
		{
			name: "no arn, empty adoption policy value",
			res: &fakeResource{arn: nil, annotations: map[string]string{
				ackv1alpha1.AnnotationAdoptionPolicy: "",
			}},
			want: false,
		},
		{
			name: "arn already set, adopt policy",
			res: &fakeResource{
				arn:         arnPtr("arn:aws:kafka:us-west-2:123456789012:cluster/abc/uuid-1"),
				annotations: adoptPolicy,
			},
			want: false,
		},
		{
			name: "empty arn string is treated as absent",
			res:  &fakeResource{arn: arnPtr(""), annotations: adoptPolicy},
			want: true,
		},
		{
			name: "nil meta object",
			res:  &fakeResource{arn: nil, nilMeta: true},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := adoption.ShouldLookupARN(tt.res); got != tt.want {
				t.Errorf("ShouldLookupARN() = %v, want %v", got, tt.want)
			}
		})
	}
}
