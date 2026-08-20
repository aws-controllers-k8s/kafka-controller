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
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/kafka/types"
)

func Test_findConfigurationARNByName(t *testing.T) {
	tests := []struct {
		name           string
		configurations []svcsdktypes.Configuration
		lookup         string
		want           *string
	}{
		{
			name:           "no configurations",
			configurations: nil,
			lookup:         "abc",
			want:           nil,
		},
		{
			name: "exact match",
			configurations: []svcsdktypes.Configuration{
				{
					Name: aws.String("abc"),
					Arn:  aws.String("arn:aws:kafka:us-west-2:123456789012:configuration/abc/uuid-1"),
				},
			},
			lookup: "abc",
			want:   aws.String("arn:aws:kafka:us-west-2:123456789012:configuration/abc/uuid-1"),
		},
		{
			name: "longer name sharing a prefix does not match",
			configurations: []svcsdktypes.Configuration{
				{
					Name: aws.String("abcd"),
					Arn:  aws.String("arn:aws:kafka:us-west-2:123456789012:configuration/abcd/uuid-1"),
				},
			},
			lookup: "abc",
			want:   nil,
		},
		{
			name: "exact match found among other configurations",
			configurations: []svcsdktypes.Configuration{
				{
					Name: aws.String("abcd"),
					Arn:  aws.String("arn:aws:kafka:us-west-2:123456789012:configuration/abcd/uuid-1"),
				},
				{
					Name: aws.String("abc"),
					Arn:  aws.String("arn:aws:kafka:us-west-2:123456789012:configuration/abc/uuid-2"),
				},
			},
			lookup: "abc",
			want:   aws.String("arn:aws:kafka:us-west-2:123456789012:configuration/abc/uuid-2"),
		},
		{
			name: "name differing only by case does not match",
			configurations: []svcsdktypes.Configuration{
				{
					Name: aws.String("ABC"),
					Arn:  aws.String("arn:aws:kafka:us-west-2:123456789012:configuration/ABC/uuid-1"),
				},
			},
			lookup: "abc",
			want:   nil,
		},
		{
			name: "nil configuration name is skipped",
			configurations: []svcsdktypes.Configuration{
				{
					Name: nil,
					Arn:  aws.String("arn:aws:kafka:us-west-2:123456789012:configuration/unknown/uuid-1"),
				},
			},
			lookup: "abc",
			want:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findConfigurationARNByName(tt.configurations, tt.lookup)
			if aws.ToString(got) != aws.ToString(tt.want) {
				t.Errorf("findConfigurationARNByName() = %v, want %v",
					aws.ToString(got), aws.ToString(tt.want))
			}
		})
	}
}
