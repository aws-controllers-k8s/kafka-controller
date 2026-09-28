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

package cluster

import (
	"testing"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/kafka/types"

	svcapitypes "github.com/aws-controllers-k8s/kafka-controller/apis/v1alpha1"
)

func Test_findClusterARNByName(t *testing.T) {
	tests := []struct {
		name     string
		clusters []svcsdktypes.ClusterInfo
		lookup   string
		want     *string
	}{
		{
			name:     "no clusters",
			clusters: nil,
			lookup:   "my-cluster",
			want:     nil,
		},
		{
			name: "exact match",
			clusters: []svcsdktypes.ClusterInfo{
				{
					ClusterName: aws.String("my-cluster"),
					ClusterArn:  aws.String("arn:aws:kafka:us-west-2:123456789012:cluster/my-cluster/abc-1"),
				},
			},
			lookup: "my-cluster",
			want:   aws.String("arn:aws:kafka:us-west-2:123456789012:cluster/my-cluster/abc-1"),
		},
		{
			name: "prefix match is not an exact match",
			clusters: []svcsdktypes.ClusterInfo{
				{
					ClusterName: aws.String("my-cluster-2"),
					ClusterArn:  aws.String("arn:aws:kafka:us-west-2:123456789012:cluster/my-cluster-2/abc-1"),
				},
			},
			lookup: "my-cluster",
			want:   nil,
		},
		{
			name: "exact match found alongside prefix matches",
			clusters: []svcsdktypes.ClusterInfo{
				{
					ClusterName: aws.String("my-cluster-2"),
					ClusterArn:  aws.String("arn:aws:kafka:us-west-2:123456789012:cluster/my-cluster-2/abc-1"),
				},
				{
					ClusterName: aws.String("my-cluster"),
					ClusterArn:  aws.String("arn:aws:kafka:us-west-2:123456789012:cluster/my-cluster/def-2"),
				},
			},
			lookup: "my-cluster",
			want:   aws.String("arn:aws:kafka:us-west-2:123456789012:cluster/my-cluster/def-2"),
		},
		{
			name: "shorter name than the filter does not match",
			clusters: []svcsdktypes.ClusterInfo{
				{
					ClusterName: aws.String("my-clust"),
					ClusterArn:  aws.String("arn:aws:kafka:us-west-2:123456789012:cluster/my-clust/abc-1"),
				},
			},
			lookup: "my-cluster",
			want:   nil,
		},
		{
			name: "name differing only by case does not match",
			clusters: []svcsdktypes.ClusterInfo{
				{
					ClusterName: aws.String("MY-CLUSTER"),
					ClusterArn:  aws.String("arn:aws:kafka:us-west-2:123456789012:cluster/MY-CLUSTER/abc-1"),
				},
			},
			lookup: "my-cluster",
			want:   nil,
		},
		{
			name: "nil cluster name is skipped",
			clusters: []svcsdktypes.ClusterInfo{
				{
					ClusterName: nil,
					ClusterArn:  aws.String("arn:aws:kafka:us-west-2:123456789012:cluster/unknown/abc-1"),
				},
			},
			lookup: "my-cluster",
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

func Test_customPreCompare_storageMode(t *testing.T) {
	tests := []struct {
		name    string
		desired *string
		latest  *string
		want    *string
	}{
		{
			name:    "unset desired adopts observed tiered",
			desired: nil,
			latest:  aws.String(string(svcsdktypes.StorageModeTiered)),
			want:    aws.String(string(svcsdktypes.StorageModeTiered)),
		},
		{
			name:    "unset desired adopts observed local",
			desired: nil,
			latest:  aws.String(string(svcsdktypes.StorageModeLocal)),
			want:    aws.String(string(svcsdktypes.StorageModeLocal)),
		},
		{
			name:    "declared desired is preserved",
			desired: aws.String(string(svcsdktypes.StorageModeTiered)),
			latest:  aws.String(string(svcsdktypes.StorageModeLocal)),
			want:    aws.String(string(svcsdktypes.StorageModeTiered)),
		},
		{
			name:    "both unset stays unset",
			desired: nil,
			latest:  nil,
			want:    nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &resource{ko: &svcapitypes.Cluster{
				Spec: svcapitypes.ClusterSpec{
					BrokerNodeGroupInfo: &svcapitypes.BrokerNodeGroupInfo{},
					StorageMode:         tt.desired,
				},
			}}
			b := &resource{ko: &svcapitypes.Cluster{
				Spec: svcapitypes.ClusterSpec{
					BrokerNodeGroupInfo: &svcapitypes.BrokerNodeGroupInfo{},
					StorageMode:         tt.latest,
				},
			}}
			customPreCompare(ackcompare.NewDelta(), a, b)
			if aws.ToString(a.ko.Spec.StorageMode) != aws.ToString(tt.want) {
				t.Errorf("customPreCompare() StorageMode = %v, want %v",
					aws.ToString(a.ko.Spec.StorageMode), aws.ToString(tt.want))
			}
		})
	}
}
