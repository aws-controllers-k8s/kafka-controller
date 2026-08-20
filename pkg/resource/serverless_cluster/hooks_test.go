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

package serverless_cluster

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/kafka/types"
)

func Test_findClusterARNByName(t *testing.T) {
	tests := []struct {
		name     string
		clusters []svcsdktypes.Cluster
		lookup   string
		want     *string
	}{
		{
			name:     "no clusters",
			clusters: nil,
			lookup:   "abc",
			want:     nil,
		},
		{
			name: "exact match",
			clusters: []svcsdktypes.Cluster{
				{
					ClusterName: aws.String("abc"),
					ClusterArn:  aws.String("arn:aws:kafka:us-west-2:123456789012:cluster/abc/uuid-1"),
				},
			},
			lookup: "abc",
			want:   aws.String("arn:aws:kafka:us-west-2:123456789012:cluster/abc/uuid-1"),
		},
		{
			name: "longer name sharing the filter prefix does not match",
			clusters: []svcsdktypes.Cluster{
				{
					ClusterName: aws.String("abcd"),
					ClusterArn:  aws.String("arn:aws:kafka:us-west-2:123456789012:cluster/abcd/uuid-1"),
				},
			},
			lookup: "abc",
			want:   nil,
		},
		{
			name: "exact match found among prefix matches",
			clusters: []svcsdktypes.Cluster{
				{
					ClusterName: aws.String("abcd"),
					ClusterArn:  aws.String("arn:aws:kafka:us-west-2:123456789012:cluster/abcd/uuid-1"),
				},
				{
					ClusterName: aws.String("abc"),
					ClusterArn:  aws.String("arn:aws:kafka:us-west-2:123456789012:cluster/abc/uuid-2"),
				},
			},
			lookup: "abc",
			want:   aws.String("arn:aws:kafka:us-west-2:123456789012:cluster/abc/uuid-2"),
		},
		{
			name: "shorter name than the filter does not match",
			clusters: []svcsdktypes.Cluster{
				{
					ClusterName: aws.String("ab"),
					ClusterArn:  aws.String("arn:aws:kafka:us-west-2:123456789012:cluster/ab/uuid-1"),
				},
			},
			lookup: "abc",
			want:   nil,
		},
		{
			name: "name differing only by case does not match",
			clusters: []svcsdktypes.Cluster{
				{
					ClusterName: aws.String("ABC"),
					ClusterArn:  aws.String("arn:aws:kafka:us-west-2:123456789012:cluster/ABC/uuid-1"),
				},
			},
			lookup: "abc",
			want:   nil,
		},
		{
			name: "nil cluster name is skipped",
			clusters: []svcsdktypes.Cluster{
				{
					ClusterName: nil,
					ClusterArn:  aws.String("arn:aws:kafka:us-west-2:123456789012:cluster/unknown/uuid-1"),
				},
			},
			lookup: "abc",
			want:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findClusterARNByName(tt.clusters, tt.lookup)
			if aws.ToString(got) != aws.ToString(tt.want) {
				t.Errorf("findClusterARNByName() = %v, want %v",
					aws.ToString(got), aws.ToString(tt.want))
			}
		})
	}
}
